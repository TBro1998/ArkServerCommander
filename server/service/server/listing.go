package server

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"ark-server-commander/models"
)

// 收录资料（简介、标签、QQ 群、Discord）：在管理端填写，随发布请求提交给收录站，显示在 arkfgf.com 的详情页上。
// 这里先校验，让服主在保存时就看到错误

const (
	listingDescriptionMax = 500
	listingTagsMax        = 8
	listingTagMax         = 16
	// 运行中修改收录资料会重新发布；距上次发布不足这个时间时等到满了再发，避免连续修改时触发收录站的频率限制
	listingRepublishGap = time.Minute
	publishStatsTTL     = 5 * time.Minute
)

var (
	qqGroupRe = regexp.MustCompile(`^[1-9][0-9]{4,11}$`)
	discordRe = regexp.MustCompile(`^https://(discord\.gg|discord\.com/invite)/[A-Za-z0-9-]{2,32}$`)
)

// normalizeListing 校验并整理收录资料：去掉首尾空白，标签去重
func normalizeListing(p models.ListingProfile) (models.ListingProfile, error) {
	out := models.ListingProfile{
		Description: strings.TrimSpace(strings.ReplaceAll(p.Description, "\r\n", "\n")),
		QQGroup:     strings.TrimSpace(p.QQGroup),
		Discord:     strings.TrimSpace(p.Discord),
		Tags:        []string{},
	}
	if utf8.RuneCountInString(out.Description) > listingDescriptionMax {
		return out, fmt.Errorf("%w：简介最多 %d 字", ErrInvalidListing, listingDescriptionMax)
	}
	seen := map[string]bool{}
	for _, tag := range p.Tags {
		tag = strings.TrimSpace(strings.NewReplacer(",", " ", "，", " ").Replace(tag))
		if tag == "" || seen[strings.ToLower(tag)] {
			continue
		}
		if utf8.RuneCountInString(tag) > listingTagMax {
			return out, fmt.Errorf("%w：标签「%s」超过 %d 字", ErrInvalidListing, tag, listingTagMax)
		}
		seen[strings.ToLower(tag)] = true
		out.Tags = append(out.Tags, tag)
	}
	if len(out.Tags) > listingTagsMax {
		return out, fmt.Errorf("%w：标签最多 %d 个", ErrInvalidListing, listingTagsMax)
	}
	if out.QQGroup != "" && !qqGroupRe.MatchString(out.QQGroup) {
		return out, fmt.Errorf("%w：QQ 群号应为 5 到 12 位数字", ErrInvalidListing)
	}
	if out.Discord != "" && !discordRe.MatchString(out.Discord) {
		return out, fmt.Errorf("%w：Discord 邀请链接应为 https://discord.gg/xxxx 的形式", ErrInvalidListing)
	}
	return out, nil
}

// listingOf 从数据库记录读出收录资料
func listingOf(server models.Server) models.ListingProfile {
	p := models.ListingProfile{
		Description: server.ListingDescription,
		QQGroup:     server.ListingQQ,
		Discord:     server.ListingDiscord,
		Tags:        []string{},
	}
	for _, tag := range strings.Split(server.ListingTags, ",") {
		if tag = strings.TrimSpace(tag); tag != "" {
			p.Tags = append(p.Tags, tag)
		}
	}
	return p
}

// setListing 把收录资料写进记录，返回是否有变化
func setListing(server *models.Server, p models.ListingProfile) bool {
	tags := strings.Join(p.Tags, ",")
	changed := server.ListingDescription != p.Description || server.ListingTags != tags ||
		server.ListingQQ != p.QQGroup || server.ListingDiscord != p.Discord
	server.ListingDescription, server.ListingTags, server.ListingQQ, server.ListingDiscord = p.Description, tags, p.QQGroup, p.Discord
	return changed
}

// listingChanged 收录资料改了：已经发布的服务器重新发布一次（距上次发布至少 listingRepublishGap）
func (p *Publisher) listingChanged(serverID uint) {
	p.mu.Lock()
	defer p.mu.Unlock()
	st, ok := p.states[serverID]
	if !ok || !st.published {
		return // 还没发布，下次可以连接时发布的就是新资料
	}
	st.published = false
	if earliest := st.lastOK.Add(listingRepublishGap); st.nextAt.Before(earliest) {
		st.nextAt = earliest
	}
}

type cachedStats struct {
	stats   *models.PublishStats
	fetched time.Time
}

var (
	statsMu    sync.Mutex
	statsCache = map[uint]cachedStats{}
)

// fetchStats 向收录站查询这台服务器的访问统计，结果缓存 publishStatsTTL
func (p *Publisher) fetchStats(server models.Server) (*models.PublishStats, error) {
	if p.api() == "" {
		return nil, fmt.Errorf("%w：收录站发布功能已关闭（ARKFGF_API 为空）", ErrPublishUnavailable)
	}
	token := accountManager.token()
	if token == "" {
		return nil, ErrAccountRequired
	}
	if !server.Published {
		return nil, ErrNotPublished
	}
	statsMu.Lock()
	if c, ok := statsCache[server.ID]; ok && p.now().Sub(c.fetched) < publishStatsTTL {
		statsMu.Unlock()
		return c.stats, nil
	}
	statsMu.Unlock()

	var stats models.PublishStats
	path := "/publish/stats?key=" + url.QueryEscape(publishKey(server.ID))
	status, _, code, err := arkfgfRequest(p.client, http.MethodGet, p.api()+path, token, nil, &stats, "收录站")
	if err != nil {
		switch {
		case isInvalidAccount(status, code):
			accountManager.invalidate()
			return nil, ErrAccountRequired
		case code == "not_published" || status == http.StatusNotFound:
			return nil, ErrNotPublished
		}
		return nil, fmt.Errorf("%w：%v", ErrPublishUnavailable, err)
	}
	statsMu.Lock()
	statsCache[server.ID] = cachedStats{stats: &stats, fetched: p.now()}
	statsMu.Unlock()
	return &stats, nil
}

// CheckReachability 用公网地址重新检测服务器的查询端口，返回检测后的状态
func (s *ServerService) CheckReachability(userID uint, serverID string) (*models.Reachability, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	return reachability.Recheck(server)
}

// GetPublishStats 查询服务器在收录站上的访问统计、详情页和徽章地址
func (s *ServerService) GetPublishStats(userID uint, serverID string) (*models.PublishStats, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	return publisher.fetchStats(server)
}

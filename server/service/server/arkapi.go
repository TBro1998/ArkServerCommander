package server

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"ark-server-commander/assets"
	"ark-server-commander/config"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// ArkApi 由管理端安装到服务端文件卷的 ShooterGame/Binaries/Win64（version.dll、msdia140.dll、config.json 等），
// 不再随服务器镜像发布。可用版本组成一个版本库：
//   - 内置版本（assets，随管理端发布，离线可用）
//   - 数据目录 arkapi/<版本>.zip：从 GitHub 下载或用户上传的版本
//   - GitHub（Michidu/ARK-Server-API）上发布、还没下载的版本，安装时再下载进版本库
// 已安装的版本号写在文件卷的 .arkapi_version 里

const (
	arkApiRepo         = "Michidu/ARK-Server-API"
	maxArkApiZip       = 100 << 20 // 压缩包大小上限
	arkApiVersionFile  = ".arkapi_version"
	arkApiDownloadTime = 5 * time.Minute
)

// arkApiVersionPattern 版本号同时用作版本库里的文件名
var arkApiVersionPattern = regexp.MustCompile(`^[0-9][0-9A-Za-z._-]{0,31}$`)

// arkApiNameVersion 从文件名或目录名里取版本号，例如 ArkApi_3.56.zip、ArkApi v3.54
var arkApiNameVersion = regexp.MustCompile(`(?i)arkapi[ _-]*v?([0-9]+(?:\.[0-9]+)+)`)

var arkApiDownloadClient = &http.Client{Timeout: arkApiDownloadTime}

// arkApiLibraryLock 保护版本库目录的写入（下载、上传、删除）
var arkApiLibraryLock sync.Mutex

// ArkApiVersion 版本库里的一个版本
type ArkApiVersion struct {
	Version     string `json:"version"`
	Builtin     bool   `json:"builtin,omitempty"`      // 随管理端发布
	Local       bool   `json:"local"`                  // 已在本地（内置、已下载或上传），安装时不需要联网
	Deletable   bool   `json:"deletable,omitempty"`    // 本地文件，可以从版本库删除
	OnGitHub    bool   `json:"on_github,omitempty"`    // GitHub 上有这个版本
	Size        int64  `json:"size,omitempty"`         // 本地压缩包大小
	PublishedAt string `json:"published_at,omitempty"` // GitHub 发布时间
	ReleaseURL  string `json:"release_url,omitempty"`  // GitHub 发布页
}

// ArkApiVersions 版本库列表
type ArkApiVersions struct {
	Versions    []ArkApiVersion `json:"versions"` // 从新到旧
	Latest      string          `json:"latest,omitempty"`
	GitHubError string          `json:"github_error,omitempty"` // 查询 GitHub 失败的原因，此时只列出本地版本
	CheckedAt   string          `json:"checked_at,omitempty"`
}

// arkApiRelease GitHub 上的一个发布
type arkApiRelease struct {
	version, url, published, download string
}

var arkApiReleases struct {
	sync.Mutex
	list    []arkApiRelease
	err     string
	checked time.Time
}

// arkApiLibraryDir 版本库目录
func arkApiLibraryDir() string {
	return filepath.Join(config.DataDir, "arkapi")
}

func arkApiLibraryPath(version string) string {
	return filepath.Join(arkApiLibraryDir(), version+".zip")
}

// arkApiWin64Dir 服务端文件卷里 ArkApi 所在的目录
func arkApiWin64Dir(serverID uint) string {
	return filepath.Join(docker_manager.ServerFilesDir(serverID), "ShooterGame", "Binaries", "Win64")
}

// compareVersions 按点分隔的数字比较版本号，数字部分相同时按字符串比较，返回 -1、0、1
func compareVersions(a, b string) int {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(pa) || i < len(pb); i++ {
		var x, y string
		if i < len(pa) {
			x = pa[i]
		}
		if i < len(pb) {
			y = pb[i]
		}
		nx, ex := strconv.Atoi(x)
		ny, ey := strconv.Atoi(y)
		switch {
		case ex == nil && ey == nil:
			if nx != ny {
				if nx < ny {
					return -1
				}
				return 1
			}
		case x != y:
			if x < y {
				return -1
			}
			return 1
		}
	}
	return 0
}

// localArkApiVersions 本地可用的版本：内置版本和版本库目录里的压缩包。目录里的同名文件优先于内置版本
func localArkApiVersions() map[string]ArkApiVersion {
	result := map[string]ArkApiVersion{
		assets.BuiltinArkApiVersion: {
			Version: assets.BuiltinArkApiVersion,
			Builtin: true,
			Local:   true,
			Size:    int64(len(assets.BuiltinArkApiZip)),
		},
	}
	entries, err := os.ReadDir(arkApiLibraryDir())
	if err != nil {
		return result
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".zip") {
			continue
		}
		version := strings.TrimSuffix(name, ".zip")
		if !arkApiVersionPattern.MatchString(version) {
			continue
		}
		v := ArkApiVersion{Version: version, Local: true, Deletable: true}
		if fi, err := e.Info(); err == nil {
			v.Size = fi.Size()
		}
		v.Builtin = version == assets.BuiltinArkApiVersion
		result[version] = v
	}
	return result
}

// latestLocalArkApi 本地最新的版本
func latestLocalArkApi() string {
	latest := ""
	for v := range localArkApiVersions() {
		if latest == "" || compareVersions(v, latest) > 0 {
			latest = v
		}
	}
	return latest
}

// cachedArkApiReleases 返回 GitHub 上的发布列表。成功结果缓存 6 小时，失败 10 分钟；refresh 时重新查询
func cachedArkApiReleases(refresh bool) ([]arkApiRelease, string, time.Time) {
	arkApiReleases.Lock()
	defer arkApiReleases.Unlock()
	ttl := 6 * time.Hour
	if arkApiReleases.err != "" {
		ttl = 10 * time.Minute
	}
	if refresh || arkApiReleases.checked.IsZero() || time.Since(arkApiReleases.checked) > ttl {
		list, err := fetchArkApiReleases()
		arkApiReleases.checked = time.Now()
		if err != nil {
			utils.Warn("查询 ArkApi 版本失败", zap.Error(err))
			arkApiReleases.err = err.Error()
		} else {
			arkApiReleases.list, arkApiReleases.err = list, ""
		}
	}
	return arkApiReleases.list, arkApiReleases.err, arkApiReleases.checked
}

func fetchArkApiReleases() ([]arkApiRelease, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/"+arkApiRepo+"/releases?per_page=30", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "ArkServerCommander/"+utils.Version)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub 接口返回 %d", resp.StatusCode)
	}
	var body []githubRelease
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&body); err != nil {
		return nil, err
	}
	return parseArkApiReleases(body), nil
}

// githubRelease GitHub 发布接口返回的字段
type githubRelease struct {
	TagName     string `json:"tag_name"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
	Assets      []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

// parseArkApiReleases 取每个正式发布的版本号和 zip 附件（优先名称里带 ArkApi 的）
func parseArkApiReleases(body []githubRelease) []arkApiRelease {
	var list []arkApiRelease
	for _, r := range body {
		if r.Draft || r.Prerelease {
			continue
		}
		version := strings.TrimPrefix(strings.TrimSpace(r.TagName), "v")
		if !arkApiVersionPattern.MatchString(version) {
			continue
		}
		download := ""
		for _, a := range r.Assets {
			if !strings.HasSuffix(strings.ToLower(a.Name), ".zip") {
				continue
			}
			if download == "" || strings.Contains(strings.ToLower(a.Name), "arkapi") {
				download = a.URL
			}
		}
		if download == "" {
			continue
		}
		list = append(list, arkApiRelease{version: version, url: r.HTMLURL, published: r.PublishedAt, download: download})
	}
	return list
}

// ListArkApiVersions 列出本地版本和 GitHub 上的版本
func (s *ServerService) ListArkApiVersions(refresh bool) ArkApiVersions {
	versions := localArkApiVersions()
	releases, ghErr, checked := cachedArkApiReleases(refresh)
	for _, r := range releases {
		v := versions[r.version]
		v.Version = r.version
		v.OnGitHub = true
		v.PublishedAt = r.published
		v.ReleaseURL = r.url
		versions[r.version] = v
	}

	result := ArkApiVersions{Versions: []ArkApiVersion{}, GitHubError: ghErr}
	if !checked.IsZero() {
		result.CheckedAt = checked.Format(time.RFC3339)
	}
	for _, v := range versions {
		result.Versions = append(result.Versions, v)
	}
	sort.Slice(result.Versions, func(i, j int) bool {
		return compareVersions(result.Versions[i].Version, result.Versions[j].Version) > 0
	})
	if len(result.Versions) > 0 {
		result.Latest = result.Versions[0].Version
	}
	return result
}

// openArkApiZip 打开版本库里的压缩包；只在 GitHub 上的版本先下载进版本库
func openArkApiZip(version string) ([]byte, error) {
	if !arkApiVersionPattern.MatchString(version) {
		return nil, ErrArkApiNotFound
	}
	if data, err := os.ReadFile(arkApiLibraryPath(version)); err == nil {
		return data, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("读取 ArkApi 压缩包失败: %w", err)
	}
	if version == assets.BuiltinArkApiVersion {
		return assets.BuiltinArkApiZip, nil
	}

	releases, _, _ := cachedArkApiReleases(false)
	for _, r := range releases {
		if r.version == version {
			return downloadArkApi(r)
		}
	}
	return nil, ErrArkApiNotFound
}

// downloadArkApi 从 GitHub 下载一个版本，校验后存入版本库
func downloadArkApi(r arkApiRelease) ([]byte, error) {
	utils.Info("下载 ArkApi", zap.String("version", r.version), zap.String("url", r.download))
	ctx, cancel := context.WithTimeout(context.Background(), arkApiDownloadTime)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.download, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrArkApiDownload, err)
	}
	req.Header.Set("User-Agent", "ArkServerCommander/"+utils.Version)
	resp, err := arkApiDownloadClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrArkApiDownload, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: GitHub 返回 %d", ErrArkApiDownload, resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxArkApiZip+1))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrArkApiDownload, err)
	}
	if len(data) > maxArkApiZip {
		return nil, fmt.Errorf("%w: 压缩包超过 100 MB", ErrArkApiDownload)
	}
	if _, _, err := readArkApiZip(bytes.NewReader(data), int64(len(data))); err != nil {
		return nil, err
	}
	if err := saveArkApiZip(r.version, data); err != nil {
		return nil, err
	}
	return data, nil
}

// saveArkApiZip 把压缩包写入版本库（先写临时文件再改名）
func saveArkApiZip(version string, data []byte) error {
	arkApiLibraryLock.Lock()
	defer arkApiLibraryLock.Unlock()
	if err := os.MkdirAll(arkApiLibraryDir(), 0o755); err != nil {
		return fmt.Errorf("创建 ArkApi 版本库目录失败: %w", err)
	}
	dest := arkApiLibraryPath(version)
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("保存 ArkApi 压缩包失败: %w", err)
	}
	if err := os.Rename(tmp, dest); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("保存 ArkApi 压缩包失败: %w", err)
	}
	return nil
}

// readArkApiZip 打开压缩包并找到 ArkApi 的根目录（version.dll 所在目录，可以有多层前缀），同时检查条目数量、大小和路径
func readArkApiZip(r io.ReaderAt, size int64) (*zip.Reader, string, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, "", fmt.Errorf("%w: 不是有效的 zip 文件", ErrInvalidArkApi)
	}
	if len(zr.File) > maxPluginZipEntries {
		return nil, "", fmt.Errorf("%w: 压缩包文件太多", ErrInvalidArkApi)
	}
	var total uint64
	root, found := "", false
	for _, f := range zr.File {
		p := f.Name
		if strings.HasPrefix(p, "__MACOSX/") {
			continue
		}
		if !safeZipPath(p) {
			return nil, "", fmt.Errorf("%w: 压缩包包含不安全的路径 %s", ErrInvalidArkApi, p)
		}
		if f.Mode()&os.ModeSymlink != 0 {
			return nil, "", fmt.Errorf("%w: 压缩包不能包含符号链接", ErrInvalidArkApi)
		}
		total += f.UncompressedSize64
		if total > maxPluginUnpacked {
			return nil, "", fmt.Errorf("%w: 解压后超过 500 MB", ErrInvalidArkApi)
		}
		if f.FileInfo().IsDir() || !strings.EqualFold(path.Base(p), "version.dll") {
			continue
		}
		dir := ""
		if d := path.Dir(p); d != "." {
			dir = d + "/"
		}
		// 插件目录里不会有 version.dll；有多个时取层级最浅的
		if !found || len(dir) < len(root) {
			root, found = dir, true
		}
	}
	if !found {
		return nil, "", fmt.Errorf("%w: 压缩包里没有 version.dll", ErrInvalidArkApi)
	}
	return zr, root, nil
}

// arkApiVersionFromZip 从文件名、压缩包里的目录名推断版本号，推断不出返回空字符串
func arkApiVersionFromZip(filename, root string) string {
	for _, s := range []string{filename, root} {
		if m := arkApiNameVersion.FindStringSubmatch(s); m != nil {
			return m[1]
		}
	}
	return ""
}

// UploadArkApi 上传 ArkApi 压缩包到版本库，返回版本号。version 为空时从文件名或压缩包里的目录名推断
func (s *ServerService) UploadArkApi(r io.ReaderAt, size int64, filename, version string) (string, error) {
	if size > maxArkApiZip {
		return "", fmt.Errorf("%w: 压缩包超过 100 MB", ErrInvalidArkApi)
	}
	_, root, err := readArkApiZip(r, size)
	if err != nil {
		return "", err
	}
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	if version == "" {
		version = arkApiVersionFromZip(filename, root)
	}
	if version == "" {
		return "", fmt.Errorf("%w: 无法从文件名识别版本号，请填写版本号", ErrInvalidArkApi)
	}
	if !arkApiVersionPattern.MatchString(version) {
		return "", fmt.Errorf("%w: 版本号只能包含数字、字母、点、下划线和连字符，并以数字开头", ErrInvalidArkApi)
	}
	data := make([]byte, size)
	if _, err := r.ReadAt(data, 0); err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("读取上传文件失败: %w", err)
	}
	if err := saveArkApiZip(version, data); err != nil {
		return "", err
	}
	utils.Info("已上传 ArkApi", zap.String("version", version))
	return version, nil
}

// DeleteArkApiVersion 从版本库删除一个本地压缩包（内置版本不能删除）
func (s *ServerService) DeleteArkApiVersion(version string) error {
	if !arkApiVersionPattern.MatchString(version) {
		return ErrArkApiNotFound
	}
	arkApiLibraryLock.Lock()
	defer arkApiLibraryLock.Unlock()
	err := os.Remove(arkApiLibraryPath(version))
	switch {
	case err == nil:
		return nil
	case errors.Is(err, os.ErrNotExist):
		if version == assets.BuiltinArkApiVersion {
			return ErrArkApiBuiltin
		}
		return ErrArkApiNotFound
	default:
		return fmt.Errorf("删除 ArkApi 压缩包失败: %w", err)
	}
}

// installedArkApiVersion 服务器已安装的 ArkApi 版本；installed 表示 version.dll 存在（版本号可能未知）
func installedArkApiVersion(serverID uint) (version string, installed bool) {
	if _, err := os.Stat(filepath.Join(arkApiWin64Dir(serverID), "version.dll")); err == nil {
		installed = true
	}
	if data, err := os.ReadFile(filepath.Join(docker_manager.ServerFilesDir(serverID), arkApiVersionFile)); err == nil {
		version = strings.TrimSpace(string(data))
	}
	return version, installed
}

// InstallArkApi 把指定版本（为空时安装最新版本）安装到服务器，返回安装的版本。只能在服务器停止时进行：运行中的游戏加载着 version.dll
func (s *ServerService) InstallArkApi(userID uint, serverID, version string) (string, error) {
	server, err := stoppedServer(userID, serverID)
	if err != nil {
		return "", err
	}
	if version == "" {
		version = s.ListArkApiVersions(false).Latest
	}
	data, err := openArkApiZip(version)
	if err != nil {
		return "", err
	}
	unlock := lockPlugins(server.ID)
	defer unlock()
	if err := installArkApi(server.ID, data, version); err != nil {
		return "", err
	}
	utils.Info("已安装 ArkApi", zap.Uint("server_id", server.ID), zap.String("version", version))
	return version, nil
}

// ensureArkApiInstalled 启动启用了 ArkApi 的服务器前调用：文件卷里还没有 version.dll 时安装本地最新的版本（不联网）
func ensureArkApiInstalled(serverID uint) error {
	if _, installed := installedArkApiVersion(serverID); installed {
		return nil
	}
	version := latestLocalArkApi()
	data, err := openArkApiZip(version)
	if err != nil {
		return fmt.Errorf("安装 ArkApi 失败: %w", err)
	}
	unlock := lockPlugins(serverID)
	defer unlock()
	if err := installArkApi(serverID, data, version); err != nil {
		return fmt.Errorf("安装 ArkApi 失败: %w", err)
	}
	utils.Info("启动前自动安装了 ArkApi", zap.Uint("server_id", serverID), zap.String("version", version))
	return nil
}

// installArkApi 把压缩包里 ArkApi 根目录下的文件解压到 Win64，写入版本号。
// 已有的 config.json 保留（用户可能改过）；ArkApi/Plugins 下的内容跳过，插件在单独的插件卷里管理
func installArkApi(serverID uint, data []byte, version string) error {
	zr, root, err := readArkApiZip(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	dest := arkApiWin64Dir(serverID)
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	_, hasConfig := os.Stat(filepath.Join(dest, "config.json"))

	for _, f := range zr.File {
		if f.FileInfo().IsDir() || strings.HasPrefix(f.Name, "__MACOSX/") || !strings.HasPrefix(f.Name, root) {
			continue
		}
		rel := strings.TrimPrefix(f.Name, root)
		lower := strings.ToLower(rel)
		if strings.HasPrefix(lower, "arkapi/plugins/") {
			continue
		}
		if lower == "config.json" && hasConfig == nil {
			continue
		}
		target := filepath.Join(dest, filepath.FromSlash(rel))
		// 先解压到临时文件再替换，中断时不会留下写了一半的 dll
		tmp := target + ".arkapi-tmp"
		if err := extractZipFile(f, tmp); err != nil {
			os.Remove(tmp)
			return fmt.Errorf("%w: %v", ErrInvalidArkApi, err)
		}
		if err := os.Rename(tmp, target); err != nil {
			os.Remove(tmp)
			return fmt.Errorf("安装 ArkApi 失败: %w", err)
		}
	}
	if err := os.WriteFile(filepath.Join(docker_manager.ServerFilesDir(serverID), arkApiVersionFile), []byte(version+"\n"), 0o644); err != nil {
		return fmt.Errorf("写入 ArkApi 版本号失败: %w", err)
	}
	return nil
}

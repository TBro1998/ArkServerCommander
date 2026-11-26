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
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// 插件框架由管理端安装到服务端文件卷的 ShooterGame/Binaries/Win64，不随服务器镜像发布。两个游戏各一套：
//   - ASE 的 ArkApi（version.dll、msdia140.dll、config.json 等），GitHub Michidu/ARK-Server-API
//   - ASA 的 AsaApi（用 AsaApiLoader.exe 启动服务端并注入插件），GitHub ArkServerApi/AsaApi
// 每套的可用版本组成一个版本库：
//   - 内置版本（assets，随管理端发布，离线可用；只有 ArkApi 有）
//   - 数据目录 arkapi/<版本>.zip（AsaApi 是 asaapi/）：从 GitHub 下载或用户上传的版本
//   - GitHub 上发布、还没下载的版本，安装时再下载进版本库
// 已安装的版本号写在文件卷的 .arkapi_version 里（两个游戏的服务器不共用文件卷，共用文件名）

const (
	maxArkApiZip       = 100 << 20 // 压缩包大小上限
	arkApiVersionFile  = ".arkapi_version"
	arkApiDownloadTime = 5 * time.Minute
)

// arkApiVersionPattern 版本号同时用作版本库里的文件名
var arkApiVersionPattern = regexp.MustCompile(`^[0-9][0-9A-Za-z._-]{0,31}$`)

var arkApiDownloadClient = &http.Client{Timeout: arkApiDownloadTime}

// arkApiLibraryLock 保护版本库目录的写入（下载、上传、删除），两个框架共用
var arkApiLibraryLock sync.Mutex

// apiFlavor 一个游戏的插件框架
type apiFlavor struct {
	game   string
	name   string // 显示名称：ArkApi / AsaApi
	repo   string // GitHub 仓库
	libDir string // 数据目录下的版本库目录
	// marker 压缩包里标志根目录的文件（安装后在 Win64 下），也用来判断服务器是否已安装
	marker string
	// nameVersion 从文件名或目录名里取版本号，例如 ArkApi_3.56.zip、ArkApi v3.54、AsaApi_2.11.zip
	nameVersion *regexp.Regexp
	builtin     bool // 有随管理端发布的内置版本（assets）
	catalog     bool // arkfgf.com 插件目录（这个游戏的那一份）里有这个框架的本体

	releases struct {
		sync.Mutex
		list    []arkApiRelease
		err     string
		checked time.Time
	}
}

var (
	aseApi = &apiFlavor{
		game: utils.GameASE, name: "ArkApi", repo: "Michidu/ARK-Server-API", libDir: "arkapi", marker: "version.dll",
		nameVersion: regexp.MustCompile(`(?i)arkapi[ _-]*v?([0-9]+(?:\.[0-9]+)+)`), builtin: true, catalog: true,
	}
	asaApi = &apiFlavor{
		game: utils.GameASA, name: "AsaApi", repo: "ArkServerApi/AsaApi", libDir: "asaapi", marker: "AsaApiLoader.exe",
		nameVersion: regexp.MustCompile(`(?i)asaapi[ _-]*v?([0-9]+(?:\.[0-9]+)+)`), catalog: true,
	}
)

// apiFor 游戏对应的插件框架
func apiFor(game string) *apiFlavor {
	if utils.NormalizeGame(game) == utils.GameASA {
		return asaApi
	}
	return aseApi
}

// ArkApiVersion 版本库里的一个版本
type ArkApiVersion struct {
	Version     string `json:"version"`
	Builtin     bool   `json:"builtin,omitempty"`      // 随管理端发布
	Local       bool   `json:"local"`                  // 已在本地（内置、已下载或上传），安装时不需要联网
	Deletable   bool   `json:"deletable,omitempty"`    // 本地文件，可以从版本库删除
	OnGitHub    bool   `json:"on_github,omitempty"`    // GitHub 上有这个版本
	OnCatalog   bool   `json:"on_catalog,omitempty"`   // arkfgf.com 插件目录（ark-server-api.com）里有这个版本
	Size        int64  `json:"size,omitempty"`         // 本地压缩包大小
	PublishedAt string `json:"published_at,omitempty"` // GitHub 发布时间
	ReleaseURL  string `json:"release_url,omitempty"`  // GitHub 发布页
}

// ArkApiVersions 版本库列表
type ArkApiVersions struct {
	Game        string          `json:"game"`     // ase（ArkApi）或 asa（AsaApi）
	Name        string          `json:"name"`     // ArkApi / AsaApi
	Versions    []ArkApiVersion `json:"versions"` // 从新到旧
	Latest      string          `json:"latest,omitempty"`
	GitHubError string          `json:"github_error,omitempty"` // 查询 GitHub 失败的原因，此时只列出本地版本
	CheckedAt   string          `json:"checked_at,omitempty"`
}

// arkApiRelease GitHub 上的一个发布
type arkApiRelease struct {
	version, url, published, download string
}

// libraryDir 版本库目录
func (a *apiFlavor) libraryDir() string {
	return filepath.Join(config.DataDir, a.libDir)
}

func (a *apiFlavor) libraryPath(version string) string {
	return filepath.Join(a.libraryDir(), version+".zip")
}

// arkApiWin64Dir 服务端文件卷里插件框架所在的目录
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

// localVersions 本地可用的版本：内置版本和版本库目录里的压缩包。目录里的同名文件优先于内置版本
func (a *apiFlavor) localVersions() map[string]ArkApiVersion {
	result := map[string]ArkApiVersion{}
	if a.builtin {
		result[assets.BuiltinArkApiVersion] = ArkApiVersion{
			Version: assets.BuiltinArkApiVersion,
			Builtin: true,
			Local:   true,
			Size:    int64(len(assets.BuiltinArkApiZip)),
		}
	}
	entries, err := os.ReadDir(a.libraryDir())
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
		v.Builtin = a.builtin && version == assets.BuiltinArkApiVersion
		result[version] = v
	}
	return result
}

// latestLocal 本地最新的版本，没有时为空
func (a *apiFlavor) latestLocal() string {
	latest := ""
	for v := range a.localVersions() {
		if latest == "" || compareVersions(v, latest) > 0 {
			latest = v
		}
	}
	return latest
}

// cachedReleases 返回 GitHub 上的发布列表。成功结果缓存 6 小时，失败 10 分钟；refresh 时重新查询
func (a *apiFlavor) cachedReleases(refresh bool) ([]arkApiRelease, string, time.Time) {
	a.releases.Lock()
	defer a.releases.Unlock()
	ttl := 6 * time.Hour
	if a.releases.err != "" {
		ttl = 10 * time.Minute
	}
	if refresh || a.releases.checked.IsZero() || time.Since(a.releases.checked) > ttl {
		list, err := a.fetchReleases()
		a.releases.checked = time.Now()
		if err != nil {
			utils.Warn("查询 "+a.name+" 版本失败", zap.Error(err))
			a.releases.err = err.Error()
		} else {
			a.releases.list, a.releases.err = list, ""
		}
	}
	return a.releases.list, a.releases.err, a.releases.checked
}

func (a *apiFlavor) fetchReleases() ([]arkApiRelease, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/"+a.repo+"/releases?per_page=30", nil)
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
	return a.parseReleases(body), nil
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

// parseReleases 取每个正式发布的版本号和 zip 附件（优先名称里带框架名的）
func (a *apiFlavor) parseReleases(body []githubRelease) []arkApiRelease {
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
		for _, asset := range r.Assets {
			if !strings.HasSuffix(strings.ToLower(asset.Name), ".zip") {
				continue
			}
			if download == "" || strings.Contains(strings.ToLower(asset.Name), strings.ToLower(a.name)) {
				download = asset.URL
			}
		}
		if download == "" {
			continue
		}
		list = append(list, arkApiRelease{version: version, url: r.HTMLURL, published: r.PublishedAt, download: download})
	}
	return list
}

// ListArkApiVersions 列出游戏（ase / asa，空为 ase）的插件框架的本地版本和 GitHub 上的版本
func (s *ServerService) ListArkApiVersions(game string, refresh bool) (ArkApiVersions, error) {
	if !utils.ValidGame(game) {
		return ArkApiVersions{}, ErrInvalidGame
	}
	return apiFor(game).list(refresh), nil
}

func (a *apiFlavor) list(refresh bool) ArkApiVersions {
	versions := a.localVersions()
	releases, ghErr, checked := a.cachedReleases(refresh)
	for _, r := range releases {
		v := versions[r.version]
		v.Version = r.version
		v.OnGitHub = true
		v.PublishedAt = r.published
		v.ReleaseURL = r.url
		versions[r.version] = v
	}
	// arkfgf.com 插件目录里的框架本体（ArkApi 在 GitHub 上只到 3.54，更新的版本发布在 ark-server-api.com；
	// AsaApi 两边都有）
	if a.catalog {
		catalog := pluginCatalogFor(a.game)
		if refresh {
			_, _ = catalog.get(true)
		} else if catalog.peek() == nil && catalog.api() != "" {
			_, _ = catalog.get(false)
		}
		if c := catalogArkApi(a.game); c != nil {
			v := versions[c.Version]
			v.Version = c.Version
			v.OnCatalog = true
			if v.PublishedAt == "" && c.UpdatedAt != nil {
				v.PublishedAt = c.UpdatedAt.UTC().Format(time.RFC3339)
			}
			versions[c.Version] = v
		}
	}

	result := ArkApiVersions{Game: a.game, Name: a.name, Versions: []ArkApiVersion{}, GitHubError: ghErr}
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

// openZip 打开版本库里的压缩包；只在 GitHub 上的版本先下载进版本库
func (a *apiFlavor) openZip(version string) ([]byte, error) {
	if !arkApiVersionPattern.MatchString(version) {
		return nil, ErrArkApiNotFound
	}
	if data, err := os.ReadFile(a.libraryPath(version)); err == nil {
		return data, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("读取 %s 压缩包失败: %w", a.name, err)
	}
	if a.builtin && version == assets.BuiltinArkApiVersion {
		return assets.BuiltinArkApiZip, nil
	}

	releases, _, _ := a.cachedReleases(false)
	for _, r := range releases {
		if r.version == version {
			return a.download(r)
		}
	}
	if a.catalog {
		if c := catalogArkApi(a.game); c != nil && c.Version == version {
			return a.download(arkApiRelease{version: c.Version, url: c.PageURL, download: c.DownloadURL})
		}
	}
	return nil, ErrArkApiNotFound
}

// download 从 GitHub 或插件目录给出的地址下载一个版本，校验后存入版本库
func (a *apiFlavor) download(r arkApiRelease) ([]byte, error) {
	utils.Info("下载 "+a.name, zap.String("version", r.version), zap.String("url", r.download))
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
		return nil, fmt.Errorf("%w: 下载地址返回 %d", ErrArkApiDownload, resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxArkApiZip+1))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrArkApiDownload, err)
	}
	if len(data) > maxArkApiZip {
		return nil, fmt.Errorf("%w: 压缩包超过 100 MB", ErrArkApiDownload)
	}
	if _, _, err := a.readZip(bytes.NewReader(data), int64(len(data))); err != nil {
		return nil, err
	}
	if err := a.saveZip(r.version, data); err != nil {
		return nil, err
	}
	return data, nil
}

// saveZip 把压缩包写入版本库（先写临时文件再改名）
func (a *apiFlavor) saveZip(version string, data []byte) error {
	arkApiLibraryLock.Lock()
	defer arkApiLibraryLock.Unlock()
	if err := os.MkdirAll(a.libraryDir(), 0o755); err != nil {
		return fmt.Errorf("创建 %s 版本库目录失败: %w", a.name, err)
	}
	dest := a.libraryPath(version)
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("保存 %s 压缩包失败: %w", a.name, err)
	}
	if err := os.Rename(tmp, dest); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("保存 %s 压缩包失败: %w", a.name, err)
	}
	return nil
}

// readZip 打开压缩包并找到框架的根目录（标志文件 version.dll / AsaApiLoader.exe 所在目录，可以有多层前缀），
// 同时检查条目数量、大小和路径
func (a *apiFlavor) readZip(r io.ReaderAt, size int64) (*zip.Reader, string, error) {
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
		p := zipEntryName(f)
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
		// AsaApi 的压缩包里也可能有 version.dll，靠 AsaApiLoader.exe 区分，免得把 ASA 的框架装到 ASE 服务器上
		if a != asaApi && !zipEntryIsDir(f) && strings.EqualFold(path.Base(p), asaApi.marker) {
			return nil, "", fmt.Errorf("%w: 这是 ASA 用的 AsaApi，不是 %s", ErrInvalidArkApi, a.name)
		}
		if zipEntryIsDir(f) || !strings.EqualFold(path.Base(p), a.marker) {
			continue
		}
		dir := ""
		if d := path.Dir(p); d != "." {
			dir = d + "/"
		}
		// 插件目录里不会有标志文件；有多个时取层级最浅的
		if !found || len(dir) < len(root) {
			root, found = dir, true
		}
	}
	if !found {
		return nil, "", fmt.Errorf("%w: 压缩包里没有 %s，不是 %s", ErrInvalidArkApi, a.marker, a.name)
	}
	return zr, root, nil
}

// versionFromZip 从文件名、压缩包里的目录名推断版本号，推断不出返回空字符串
func (a *apiFlavor) versionFromZip(filename, root string) string {
	for _, s := range []string{filename, root} {
		if m := a.nameVersion.FindStringSubmatch(s); m != nil {
			return m[1]
		}
	}
	return ""
}

// UploadArkApi 上传游戏（ase / asa）的插件框架压缩包到版本库，返回版本号。version 为空时从文件名或压缩包里的目录名推断
func (s *ServerService) UploadArkApi(game string, r io.ReaderAt, size int64, filename, version string) (string, error) {
	if !utils.ValidGame(game) {
		return "", ErrInvalidGame
	}
	a := apiFor(game)
	if size > maxArkApiZip {
		return "", fmt.Errorf("%w: 压缩包超过 100 MB", ErrInvalidArkApi)
	}
	_, root, err := a.readZip(r, size)
	if err != nil {
		return "", err
	}
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	if version == "" {
		version = a.versionFromZip(filename, root)
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
	if err := a.saveZip(version, data); err != nil {
		return "", err
	}
	utils.Info("已上传 "+a.name, zap.String("version", version))
	return version, nil
}

// DeleteArkApiVersion 从游戏（ase / asa）的版本库删除一个本地压缩包（内置版本不能删除）
func (s *ServerService) DeleteArkApiVersion(game, version string) error {
	if !utils.ValidGame(game) {
		return ErrInvalidGame
	}
	a := apiFor(game)
	if !arkApiVersionPattern.MatchString(version) {
		return ErrArkApiNotFound
	}
	arkApiLibraryLock.Lock()
	defer arkApiLibraryLock.Unlock()
	err := os.Remove(a.libraryPath(version))
	switch {
	case err == nil:
		return nil
	case errors.Is(err, os.ErrNotExist):
		if a.builtin && version == assets.BuiltinArkApiVersion {
			return ErrArkApiBuiltin
		}
		return ErrArkApiNotFound
	default:
		return fmt.Errorf("删除 %s 压缩包失败: %w", a.name, err)
	}
}

// installedArkApiVersion 服务器已安装的插件框架版本；installed 表示标志文件（ASE 的 version.dll、
// ASA 的 AsaApiLoader.exe）存在（版本号可能未知）
func installedArkApiVersion(server models.Server) (version string, installed bool) {
	if _, err := os.Stat(filepath.Join(arkApiWin64Dir(server.ID), apiFor(server.Game).marker)); err == nil {
		installed = true
	}
	if data, err := os.ReadFile(filepath.Join(docker_manager.ServerFilesDir(server.ID), arkApiVersionFile)); err == nil {
		version = strings.TrimSpace(string(data))
	}
	return version, installed
}

// InstallArkApi 把指定版本（为空时安装最新版本）安装到服务器（ASE 装 ArkApi，ASA 装 AsaApi），返回安装的版本。
// 只能在服务器停止时进行：运行中的游戏加载着框架的文件
func (s *ServerService) InstallArkApi(userID uint, serverID, version string) (string, error) {
	server, err := stoppedServer(userID, serverID)
	if err != nil {
		return "", err
	}
	a := apiFor(server.Game)
	if version == "" {
		version = a.list(false).Latest
	}
	data, err := a.openZip(version)
	if err != nil {
		return "", err
	}
	unlock := lockPlugins(server.ID)
	defer unlock()
	if err := a.install(server.ID, data, version); err != nil {
		return "", err
	}
	utils.Info("已安装 "+a.name, zap.Uint("server_id", server.ID), zap.String("version", version))
	return version, nil
}

// ensureArkApiInstalled 启动启用了 ArkApi / AsaApi 的服务器前调用：文件卷里还没有时安装本地最新的版本（不联网）；
// 本地没有任何版本时（AsaApi 没有内置版本）安装 GitHub 上最新的版本
func ensureArkApiInstalled(server models.Server) error {
	if _, installed := installedArkApiVersion(server); installed {
		return nil
	}
	a := apiFor(server.Game)
	version := a.latestLocal()
	if version == "" {
		version = a.list(false).Latest
	}
	if version == "" {
		return fmt.Errorf("安装 %s 失败: %w", a.name, ErrArkApiNotFound)
	}
	data, err := a.openZip(version)
	if err != nil {
		return fmt.Errorf("安装 %s 失败: %w", a.name, err)
	}
	unlock := lockPlugins(server.ID)
	defer unlock()
	if err := a.install(server.ID, data, version); err != nil {
		return fmt.Errorf("安装 %s 失败: %w", a.name, err)
	}
	utils.Info("启动前自动安装了 "+a.name, zap.Uint("server_id", server.ID), zap.String("version", version))
	return nil
}

// install 把压缩包里框架根目录下的文件解压到 Win64，写入版本号。
// 已有的 config.json 保留（用户可能改过）；ArkApi/Plugins 下的内容跳过，插件在单独的插件卷里管理
func (a *apiFlavor) install(serverID uint, data []byte, version string) error {
	zr, root, err := a.readZip(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	dest := arkApiWin64Dir(serverID)
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	_, hasConfig := os.Stat(filepath.Join(dest, "config.json"))

	for _, f := range zr.File {
		name := zipEntryName(f)
		if zipEntryIsDir(f) || strings.HasPrefix(name, "__MACOSX/") || !strings.HasPrefix(name, root) {
			continue
		}
		rel := strings.TrimPrefix(name, root)
		lower := strings.ToLower(rel)
		if strings.HasPrefix(lower, "arkapi/plugins/") {
			continue
		}
		if lower == "config.json" && hasConfig == nil {
			continue
		}
		target := filepath.Join(dest, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("安装 %s 失败: %w", a.name, err)
		}
		// 先解压到临时文件再替换，中断时不会留下写了一半的 dll
		tmp := target + ".arkapi-tmp"
		if err := extractZipFile(f, tmp); err != nil {
			os.Remove(tmp)
			return fmt.Errorf("%w: %v", ErrInvalidArkApi, err)
		}
		if err := os.Rename(tmp, target); err != nil {
			os.Remove(tmp)
			return fmt.Errorf("安装 %s 失败: %w", a.name, err)
		}
	}
	if err := os.WriteFile(filepath.Join(docker_manager.ServerFilesDir(serverID), arkApiVersionFile), []byte(version+"\n"), 0o644); err != nil {
		return fmt.Errorf("写入 %s 版本号失败: %w", a.name, err)
	}
	return nil
}

package server

import (
	"crypto/ed25519"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	_ "github.com/glebarez/go-sqlite" // database/sql 驱动 "sqlite"
)

// 玩家站登录后看自己的点数：Commander 定期只读 ArkShop 的 SQLite 库（复制一份再读，不锁住插件），
// 把每个 SteamID 的点数、已花点数、礼包剩余次数写进快照目录的 players/<服务器 ID>.json。
// 这些文件不公开，玩家站只给登录的本人返回他自己的那一条。只支持 ASE + SQLite：ASA 的 ArkShop 按 EOS ID 记账，
// 玩家站用 Steam 登录对不上；用 MySQL 的（UseMysql）或改了库路径的（DbPathOverride）不读

const (
	arkShopDBName       = "ArkShop.db"
	maxArkShopDBSize    = 256 << 20
	maxShopPlayers      = 200000
	identityKeyFile     = "identity.key"
	portalContainerUser = 65532 // 玩家站镜像（distroless nonroot）的用户
)

type portalPlayer struct {
	Points     int64          `json:"points"`
	TotalSpent int64          `json:"total_spent"`
	Kits       map[string]int `json:"kits,omitempty"`
}

type portalPlayers struct {
	UpdatedAt string                  `json:"updated_at"`
	Players   map[string]portalPlayer `json:"players"`
}

// arkShopUsesOwnSQLite ArkShop 用的是插件目录里自己的 SQLite 库（没有改成 MySQL、没有改库路径）
func arkShopUsesOwnSQLite(config []byte) bool {
	var cfg struct {
		Mysql struct {
			UseMysql bool `json:"UseMysql"`
		} `json:"Mysql"`
		General struct {
			DbPathOverride string `json:"DbPathOverride"`
		} `json:"General"`
	}
	if err := json.Unmarshal(trimBOM(config), &cfg); err != nil {
		return false
	}
	return !cfg.Mysql.UseMysql && strings.TrimSpace(cfg.General.DbPathOverride) == ""
}

func trimBOM(b []byte) []byte { return []byte(strings.TrimPrefix(string(b), "\xEF\xBB\xBF")) }

// arkShopDBStat ArkShop 库文件的状态（判断有没有变化）；没有时返回 false
func arkShopDBStat(serverID uint) (fileStat, bool) {
	root, err := os.OpenRoot(docker_manager.ServerPluginsDir(serverID))
	if err != nil {
		return fileStat{}, false
	}
	defer root.Close()
	dir, ok := arkShopDirName(root)
	if !ok {
		return fileStat{}, false
	}
	info, err := root.Stat(dir + "/" + arkShopDBName)
	if err != nil || !info.Mode().IsRegular() {
		return fileStat{}, false
	}
	return fileStat{modTime: info.ModTime(), size: info.Size()}, true
}

type fileStat struct {
	modTime time.Time
	size    int64
}

func arkShopDirName(root *os.Root) (string, bool) {
	entries, err := fs.ReadDir(root.FS(), ".")
	if err != nil {
		return "", false
	}
	for _, e := range entries {
		if e.IsDir() && strings.EqualFold(e.Name(), arkShopPluginName) {
			return e.Name(), true
		}
	}
	return "", false
}

// readArkShopPlayers 复制 ArkShop 的库到临时文件，只读打开，读出所有玩家
func readArkShopPlayers(serverID uint) (map[string]portalPlayer, error) {
	root, err := os.OpenRoot(docker_manager.ServerPluginsDir(serverID))
	if err != nil {
		return nil, err
	}
	defer root.Close()
	dir, ok := arkShopDirName(root)
	if !ok {
		return nil, errNoArkShop
	}
	src, err := root.Open(dir + "/" + arkShopDBName)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	tmp, err := os.CreateTemp("", "arkshop-*.db")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	n, err := io.Copy(tmp, io.LimitReader(src, maxArkShopDBSize+1))
	tmp.Close()
	if err != nil {
		return nil, err
	}
	if n > maxArkShopDBSize {
		return nil, errors.New("ArkShop 的库太大")
	}

	db, err := sql.Open("sqlite", "file:"+filepath.ToSlash(tmp.Name())+"?mode=ro")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT SteamId, Points, TotalSpent, Kits FROM Players LIMIT ?", maxShopPlayers)
	if err != nil {
		return nil, fmt.Errorf("读取 ArkShop 的玩家表失败: %w", err)
	}
	defer rows.Close()
	out := map[string]portalPlayer{}
	for rows.Next() {
		var steamID any
		var points, spent sql.NullInt64
		var kits sql.NullString
		if err := rows.Scan(&steamID, &points, &spent, &kits); err != nil {
			continue
		}
		id := fmt.Sprint(steamID)
		if b, ok := steamID.([]byte); ok {
			id = string(b)
		}
		if !steamIDRe.MatchString(id) {
			continue
		}
		p := portalPlayer{Points: points.Int64, TotalSpent: spent.Int64}
		var k map[string]struct {
			Amount int `json:"Amount"`
		}
		if kits.Valid && json.Unmarshal([]byte(kits.String), &k) == nil {
			for name, v := range k {
				if len(p.Kits) < 200 {
					if p.Kits == nil {
						p.Kits = map[string]int{}
					}
					p.Kits[shopText(name, maxShopShortRunes)] = v.Amount
				}
			}
		}
		out[id] = p
	}
	return out, rows.Err()
}

// playersCache 每台服务器上次读到的库文件状态和内容，库文件没变化时不重新读
type playersCache struct {
	stat    fileStat
	players map[string]portalPlayer
}

// shopPlayers 一台服务器的玩家点数；不支持（ASA、MySQL、改了路径）或读不出来时返回 nil
func (m *PortalManager) shopPlayers(s models.Server, config []byte) map[string]portalPlayer {
	if utils.NormalizeGame(s.Game) != utils.GameASE || !arkShopUsesOwnSQLite(config) {
		return nil
	}
	st, ok := arkShopDBStat(s.ID)
	if !ok {
		return nil
	}
	m.mu.Lock()
	c, cached := m.players[s.ID]
	m.mu.Unlock()
	if cached && c.stat == st {
		return c.players
	}
	players, err := readArkShopPlayers(s.ID)
	if err != nil {
		utils.Debug("读取 ArkShop 玩家点数失败：" + err.Error())
		players = nil
	}
	m.mu.Lock()
	m.players[s.ID] = playersCache{stat: st, players: players}
	m.mu.Unlock()
	return players
}

// ensurePortalIdentity 玩家站的身份密钥（Ed25519 种子，base64），没有时生成。玩家站用它向联机助手证明自己，
// arkfgf.com 签发的登录令牌只对这个身份有效；也用来签玩家站自己的会话。文件只给玩家站容器的用户读
func ensurePortalIdentity(dir string) error {
	path := filepath.Join(dir, identityKeyFile)
	if data, err := os.ReadFile(path); err == nil {
		if seed, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(data))); err == nil && len(seed) == ed25519.SeedSize {
			return nil
		}
	}
	seed := make([]byte, ed25519.SeedSize)
	if _, err := rand.Read(seed); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(base64.StdEncoding.EncodeToString(seed)), 0o600); err != nil {
		return err
	}
	// 只给玩家站容器的用户读；不支持 chown 的系统（Windows 上开发）放宽成所有人可读
	if err := os.Chown(tmp, portalContainerUser, portalContainerUser); err != nil {
		_ = os.Chmod(tmp, 0o644)
	}
	return os.Rename(tmp, path)
}

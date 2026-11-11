package config

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// jwtSecretFileName 自动生成的 JWT 密钥文件名，与数据库放在同一目录
	jwtSecretFileName  = "jwt_secret"
	minJWTSecretLength = 32
)

var (
	JWTSecret  []byte
	DBPath     = "ark_server.db"
	ServerPort = "8080"

	// ContainerTZ ARK 服务器容器的时区，通过 CONTAINER_TZ 修改，仅对新建或重建的容器生效
	ContainerTZ = "Asia/Shanghai"

	// ListenHost 监听地址，默认仅本机回环；容器内由 docker-entrypoint.sh 指定
	ListenHost = "127.0.0.1"

	// CORSAllowedOrigins 允许跨域的来源白名单，为空时仅同源访问
	CORSAllowedOrigins []string

	// BackupDir 存档备份目录，默认为数据库所在目录下的 backups（镜像内为 /data/backups）
	BackupDir string

	// PublicHost 玩家连接用的公网地址（IP 或域名），为空时自动探测公网 IP
	PublicHost string
	// DataDir 数据目录（数据库所在目录，镜像内为 /data）。ARK 服务器的服务端文件、存档、插件和集群数据
	// 放在其下的 servers/<id>/ 和 clusters/<cluster_id>/，宿主机上可以直接管理
	DataDir string
	// HostDataDir DataDir 在宿主机上的路径，通过 HOST_DATA_DIR 指定；为空时自动探测（见 docker_manager.InitStorage）
	HostDataDir string
)

// 弱密钥黑名单，仅用于校验手动指定的 JWT_SECRET
var weakSecrets = []string{
	"ark-server-commander-secret-key",
	"secret",
	"password",
	"123456",
	"default",
	"changeme",
	"test",
}

func InitConfig() error {
	// 读取其他配置（密钥文件放在数据库所在目录，需要先确定 DBPath）
	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		DBPath = dbPath
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		ServerPort = port
	}

	if tz := os.Getenv("CONTAINER_TZ"); tz != "" {
		ContainerTZ = tz
	}

	if host := os.Getenv("LISTEN_HOST"); host != "" {
		ListenHost = host
	}

	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		CORSAllowedOrigins = strings.Split(origins, ",")
	}

	BackupDir = os.Getenv("BACKUP_DIR")
	if BackupDir == "" {
		BackupDir = filepath.Join(filepath.Dir(DBPath), "backups")
	}
	PublicHost = strings.TrimSpace(os.Getenv("PUBLIC_HOST"))

	dataDir, err := filepath.Abs(filepath.Dir(DBPath))
	if err != nil {
		return fmt.Errorf("解析数据目录失败: %w", err)
	}
	DataDir = dataDir
	HostDataDir = strings.TrimSpace(os.Getenv("HOST_DATA_DIR"))

	secret, err := loadJWTSecret()
	if err != nil {
		return err
	}
	JWTSecret = secret

	return nil
}

// loadJWTSecret 优先使用环境变量 JWT_SECRET；未设置时从数据库目录下的密钥文件读取，文件不存在则生成
func loadJWTSecret() ([]byte, error) {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		if err := validateJWTSecret(secret); err != nil {
			return nil, err
		}
		return []byte(secret), nil
	}

	path := filepath.Join(filepath.Dir(DBPath), jwtSecretFileName)
	data, err := os.ReadFile(path)
	if err == nil {
		secret := strings.TrimSpace(string(data))
		if len(secret) < minJWTSecretLength {
			return nil, fmt.Errorf("JWT 密钥文件 %s 内容无效，删除后重启会自动重新生成", path)
		}
		return []byte(secret), nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("读取 JWT 密钥文件失败: %w", err)
	}

	buf := make([]byte, 48)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("生成 JWT 密钥失败: %w", err)
	}
	secret := base64.StdEncoding.EncodeToString(buf)

	// O_EXCL 避免覆盖已有密钥（覆盖会让所有已登录会话失效）
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return nil, fmt.Errorf("保存 JWT 密钥文件失败: %w", err)
	}
	defer f.Close()
	if _, err := f.WriteString(secret + "\n"); err != nil {
		return nil, fmt.Errorf("保存 JWT 密钥文件失败: %w", err)
	}

	return []byte(secret), nil
}

// validateJWTSecret 校验手动指定的 JWT 密钥强度
func validateJWTSecret(secret string) error {
	if len(secret) < minJWTSecretLength {
		return fmt.Errorf("JWT_SECRET must be at least %d characters long (current: %d)", minJWTSecretLength, len(secret))
	}

	secretLower := strings.ToLower(secret)
	for _, weak := range weakSecrets {
		if strings.Contains(secretLower, weak) {
			return fmt.Errorf("JWT_SECRET contains weak/common password pattern: '%s'", weak)
		}
	}

	return nil
}

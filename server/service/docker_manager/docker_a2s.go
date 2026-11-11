package docker_manager

import (
	"encoding/json"
	"fmt"
	"time"
)

// a2sScript 在 ARK 容器内向本机查询端口发送 A2S_INFO，输出 JSON。
// 能拿到响应说明游戏已经加载完地图、可以连接；同时得到在线人数和带版本号的服务器名。
// 新版协议可能先返回 challenge（0x41），需要带上 challenge 重发一次
const a2sScript = `import os, socket, json, sys
port = int(os.environ['QUERY_PORT'])
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.settimeout(3)
req = b'\xff\xff\xff\xffTSource Engine Query\x00'
s.sendto(req, ('127.0.0.1', port))
data, _ = s.recvfrom(4096)
if len(data) >= 9 and data[4] == 0x41:
    s.sendto(req + data[5:9], ('127.0.0.1', port))
    data, _ = s.recvfrom(4096)
if len(data) < 6 or data[4] != 0x49:
    sys.exit('unexpected A2S response')
i = 6
def cstr():
    global i
    j = data.index(b'\x00', i)
    v = data[i:j].decode('utf-8', 'replace')
    i = j + 1
    return v
name = cstr(); mp = cstr(); folder = cstr(); game = cstr()
i += 2
print(json.dumps({'name': name, 'map': mp, 'players': data[i], 'max_players': data[i + 1]}))
`

// A2SInfo 查询端口返回的服务器信息
type A2SInfo struct {
	Name       string `json:"name"`
	Map        string `json:"map"`
	Players    int    `json:"players"`
	MaxPlayers int    `json:"max_players"`
}

// QueryServerInfo 通过查询端口获取服务器信息；游戏还没加载完时返回错误
func (dm *DockerManager) QueryServerInfo(containerName string, queryPort int) (*A2SInfo, error) {
	out, err := dm.execPython(containerName, a2sScript, []string{fmt.Sprintf("QUERY_PORT=%d", queryPort)}, 10*time.Second)
	if err != nil {
		return nil, err
	}
	var info A2SInfo
	if err := json.Unmarshal([]byte(out), &info); err != nil {
		return nil, fmt.Errorf("解析查询结果失败: %v", err)
	}
	return &info, nil
}

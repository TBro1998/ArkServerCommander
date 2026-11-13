package servers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ark-server-commander/models"

	"github.com/gin-gonic/gin"
)

// RestartServer 重启服务器
// @Summary 重启服务器
// @Description 运行中的服务器先保存世界并停止，再按当前配置启动；已停止的直接启动。后台执行
// @Tags 服务器管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]string "已开始重启"
// @Failure 404 {object} map[string]string "服务器不存在"
// @Failure 409 {object} map[string]string "服务器正在启停"
// @Router /servers/{id}/restart [post]
func RestartServer(c *gin.Context) {
	if err := serverService.RestartServer(c.GetUint("user_id"), c.Param("id")); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "服务器重启已开始"})
}

// UpdateGameRequest 更新游戏本体 / Mod 请求
type UpdateGameRequest struct {
	Game   bool     `json:"game"`              // 校验并更新游戏本体
	Mods   bool     `json:"mods"`              // 重新下载全部 Mod
	ModIDs []string `json:"mod_ids,omitempty"` // 只重新下载这些 Mod（mods 为 false 时生效）
}

// UpdateGame 更新游戏本体和 Mod
// @Summary 更新游戏本体和 Mod
// @Description 清除安装标记后启动服务器（运行中的先停止），入口脚本在启动时用 steamcmd 校验更新游戏本体、重新下载全部 Mod 或 mod_ids 指定的 Mod
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param body body UpdateGameRequest true "更新内容"
// @Success 200 {object} map[string]string "已开始更新"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 409 {object} map[string]string "服务器正在启停"
// @Router /servers/{id}/update-game [post]
func UpdateGame(c *gin.Context) {
	var req UpdateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	if err := serverService.UpdateGame(c.GetUint("user_id"), c.Param("id"), req.Game, req.Mods, req.ModIDs); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新已开始，服务器启动时会下载更新"})
}

// GetLogs 获取容器日志
// @Summary 获取服务器日志
// @Description 返回容器最近 tail 行日志（入口脚本、steamcmd 和游戏的标准输出）
// @Tags 服务器管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param tail query int false "行数（默认 500，最多 5000）"
// @Success 200 {object} map[string][]string "日志行"
// @Router /servers/{id}/logs [get]
func GetLogs(c *gin.Context) {
	tail, _ := strconv.Atoi(c.DefaultQuery("tail", "500"))
	lines, err := serverService.GetLogs(c.GetUint("user_id"), c.Param("id"), tail)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": lines})
}

// StreamLogs 实时日志（SSE）
// @Summary 实时跟随服务器日志
// @Description Server-Sent Events：先发送最近 tail 行，之后持续推送新日志。每条消息的 data 是一行日志；容器停止时发送 end 事件
// @Tags 服务器管理
// @Produce text/event-stream
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param tail query int false "先发送的历史行数（默认 200）"
// @Router /servers/{id}/logs/stream [get]
func StreamLogs(c *gin.Context) {
	userID := c.GetUint("user_id")
	serverID := c.Param("id")
	tail, _ := strconv.Atoi(c.DefaultQuery("tail", "200"))

	// 先检查服务器归属，出错时还能返回普通的 JSON 错误
	if _, err := serverService.GetLogs(userID, serverID, 1); err != nil {
		respondServiceError(c, err)
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	c.Writer.Flush()

	write := func(line string) error {
		if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", strings.ReplaceAll(line, "\n", " ")); err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}

	err := serverService.FollowLogs(c.Request.Context(), userID, serverID, tail, write)
	if c.Request.Context().Err() != nil {
		return
	}
	if err != nil {
		fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", err.Error())
	}
	fmt.Fprint(c.Writer, "event: end\ndata: \n\n")
	c.Writer.Flush()
}

// RconCommandRequest RCON 命令请求
type RconCommandRequest struct {
	Command string `json:"command" binding:"required"`
}

// RconCommand 执行 RCON 命令
// @Summary 执行 RCON 命令
// @Description 在运行中的服务器上执行一条 RCON 命令（如 Broadcast、SaveWorld、KickPlayer），返回输出
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param body body RconCommandRequest true "命令"
// @Success 200 {object} map[string]string "命令输出"
// @Failure 400 {object} map[string]string "服务器未运行或命令无效"
// @Router /servers/{id}/rcon/command [post]
func RconCommand(c *gin.Context) {
	var req RconCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	out, err := serverService.RconExec(c.GetUint("user_id"), c.Param("id"), req.Command)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"output": out}})
}

// GetGameLog 游戏日志
// @Summary 获取游戏日志
// @Description 返回存档目录 Logs/ShooterGame.log 最后 tail 行（默认 1000，最多 5000）。直接读数据目录，服务器停止后也能看
// @Tags 服务器管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param tail query int false "行数"
// @Success 200 {object} map[string]server.GameLog "游戏日志"
// @Router /servers/{id}/logs/game [get]
func GetGameLog(c *gin.Context) {
	tail, _ := strconv.Atoi(c.DefaultQuery("tail", "1000"))
	if tail <= 0 || tail > 5000 {
		tail = 1000
	}
	log, err := serverService.GetGameLog(c.GetUint("user_id"), c.Param("id"), tail)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": log})
}

// ListServerMods 服务器的 Mod 及更新状态
// @Summary 服务器的 Mod 及更新状态
// @Description 按加载顺序返回服务器配置的 Mod，附带创意工坊信息、是否已安装、安装时间和是否有更新。refresh=true 时忽略缓存重新查询 Steam
// @Tags 服务器管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param refresh query bool false "重新查询 Steam"
// @Success 200 {object} map[string][]server.ServerMod "Mod 列表"
// @Router /servers/{id}/mods [get]
func ListServerMods(c *gin.Context) {
	mods, err := serverService.GetServerMods(c.GetUint("user_id"), c.Param("id"), c.Query("refresh") == "true")
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": mods})
}

// ListPlayers 在线玩家
// @Summary 在线玩家列表
// @Description 通过 RCON ListPlayers 获取在线玩家（名称和 SteamID）
// @Tags 服务器管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string][]server.Player "玩家列表"
// @Router /servers/{id}/players [get]
func ListPlayers(c *gin.Context) {
	players, err := serverService.ListPlayers(c.GetUint("user_id"), c.Param("id"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": players})
}

// ListBackups 备份列表
// @Summary 存档备份列表
// @Tags 存档备份
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string][]server.BackupInfo "备份列表，最新的在前"
// @Router /servers/{id}/backups [get]
func ListBackups(c *gin.Context) {
	list, err := serverService.ListBackups(c.GetUint("user_id"), c.Param("id"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// CreateBackup 立即备份
// @Summary 立即备份存档
// @Description 备份数据卷里的存档和配置文件（不含日志）；运行中的服务器会先保存世界
// @Tags 存档备份
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]server.BackupInfo "备份信息"
// @Failure 409 {object} map[string]string "正在备份或启停"
// @Router /servers/{id}/backups [post]
func CreateBackup(c *gin.Context) {
	info, err := serverService.CreateBackup(c.GetUint("user_id"), c.Param("id"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "备份完成", "data": info})
}

// UploadBackup 上传备份
// @Summary 上传备份文件
// @Description 上传本系统导出的 .tar.gz 备份（multipart 字段名 file），校验后加入备份列表
// @Tags 存档备份
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param file formData file true "备份文件"
// @Success 200 {object} map[string]server.BackupInfo "备份信息"
// @Failure 400 {object} map[string]string "格式无效"
// @Router /servers/{id}/backups/upload [post]
func UploadBackup(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		respondBadRequest(c)
		return
	}
	f, err := file.Open()
	if err != nil {
		respondBadRequest(c)
		return
	}
	defer f.Close()
	info, err := serverService.UploadBackup(c.GetUint("user_id"), c.Param("id"), f)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "上传成功", "data": info})
}

// DownloadBackup 下载备份
// @Summary 下载备份文件
// @Tags 存档备份
// @Produce application/gzip
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param name path string true "备份文件名"
// @Success 200 {file} file "备份文件"
// @Failure 404 {object} map[string]string "备份不存在"
// @Router /servers/{id}/backups/{name} [get]
func DownloadBackup(c *gin.Context) {
	path, err := serverService.BackupPath(c.GetUint("user_id"), c.Param("id"), c.Param("name"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.FileAttachment(path, fmt.Sprintf("server-%s_%s", c.Param("id"), c.Param("name")))
}

// DeleteBackup 删除备份
// @Summary 删除备份文件
// @Tags 存档备份
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param name path string true "备份文件名"
// @Success 200 {object} map[string]string "删除成功"
// @Router /servers/{id}/backups/{name} [delete]
func DeleteBackup(c *gin.Context) {
	if err := serverService.DeleteBackup(c.GetUint("user_id"), c.Param("id"), c.Param("name")); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "备份已删除"})
}

// RestoreBackup 恢复备份
// @Summary 从备份恢复存档
// @Description 只能在服务器停止时执行：先自动备份当前存档（pre-restore），再用所选备份覆盖存档和配置文件
// @Tags 存档备份
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param name path string true "备份文件名"
// @Success 200 {object} map[string]string "恢复成功"
// @Failure 400 {object} map[string]string "服务器未停止"
// @Router /servers/{id}/backups/{name}/restore [post]
func RestoreBackup(c *gin.Context) {
	if err := serverService.RestoreBackup(c.GetUint("user_id"), c.Param("id"), c.Param("name")); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "存档已恢复"})
}

// ListSchedules 定时任务列表
// @Summary 定时任务列表
// @Tags 定时任务
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string][]models.Schedule "定时任务"
// @Router /servers/{id}/schedules [get]
func ListSchedules(c *gin.Context) {
	list, err := serverService.ListSchedules(c.GetUint("user_id"), c.Param("id"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// CreateSchedule 新建定时任务
// @Summary 新建定时任务
// @Description action: restart / backup / update / command / start / stop；cron 为 5 段表达式，时区为 CONTAINER_TZ
// @Tags 定时任务
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param body body models.ScheduleRequest true "定时任务"
// @Success 200 {object} map[string]models.Schedule "定时任务"
// @Router /servers/{id}/schedules [post]
func CreateSchedule(c *gin.Context) {
	var req models.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	sch, err := serverService.CreateSchedule(c.GetUint("user_id"), c.Param("id"), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sch})
}

// UpdateSchedule 修改定时任务
// @Summary 修改定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param scheduleId path int true "定时任务ID"
// @Param body body models.ScheduleRequest true "定时任务"
// @Success 200 {object} map[string]models.Schedule "定时任务"
// @Router /servers/{id}/schedules/{scheduleId} [put]
func UpdateSchedule(c *gin.Context) {
	var req models.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	sch, err := serverService.UpdateSchedule(c.GetUint("user_id"), c.Param("id"), c.Param("scheduleId"), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sch})
}

// DeleteSchedule 删除定时任务
// @Summary 删除定时任务
// @Tags 定时任务
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param scheduleId path int true "定时任务ID"
// @Success 200 {object} map[string]string "删除成功"
// @Router /servers/{id}/schedules/{scheduleId} [delete]
func DeleteSchedule(c *gin.Context) {
	if err := serverService.DeleteSchedule(c.GetUint("user_id"), c.Param("id"), c.Param("scheduleId")); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "定时任务已删除"})
}

// RunSchedule 立即执行定时任务
// @Summary 立即执行一次定时任务
// @Tags 定时任务
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param scheduleId path int true "定时任务ID"
// @Success 200 {object} map[string]string "已开始执行"
// @Router /servers/{id}/schedules/{scheduleId}/run [post]
func RunSchedule(c *gin.Context) {
	if err := serverService.RunScheduleNow(c.GetUint("user_id"), c.Param("id"), c.Param("scheduleId")); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已开始执行"})
}

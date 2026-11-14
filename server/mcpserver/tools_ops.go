package mcpserver

import (
	"context"
	"fmt"

	"ark-server-commander/models"
	"ark-server-commander/service/server"

	"github.com/gin-gonic/gin/binding"
)

type logsInput struct {
	ServerID uint `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Tail     int  `json:"tail,omitempty" jsonschema:"Number of lines from the end, 1-5000 (default 200)"`
}

type rconInput struct {
	ServerID uint   `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Command  string `json:"command" jsonschema:"RCON command, e.g. ServerChat Hello, Broadcast Restart in 5 minutes, SaveWorld, KickPlayer <SteamID>, DestroyWildDinos"`
}

type serverModsInput struct {
	ServerID uint `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Refresh  bool `json:"refresh,omitempty" jsonschema:"Query Steam Workshop again instead of using the cache"`
}

type backupInput struct {
	ServerID uint   `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Name     string `json:"name" jsonschema:"Backup file name (from list_backups)"`
}

type scheduleIDInput struct {
	ServerID   uint `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	ScheduleID uint `json:"schedule_id" jsonschema:"Schedule ID (from list_schedules)"`
}

type createScheduleInput struct {
	ServerID    uint   `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Action      string `json:"action" jsonschema:"One of: restart, stop, start, backup, update (update game and mods, then restart), command (run an RCON command)"`
	Cron        string `json:"cron" jsonschema:"5-field cron expression in the server time zone (CONTAINER_TZ), e.g. 0 4 * * * for 04:00 daily, or @daily / @hourly"`
	Enabled     *bool  `json:"enabled,omitempty" jsonschema:"Default true"`
	WarnMinutes int    `json:"warn_minutes,omitempty" jsonschema:"For restart/stop/update: broadcast a countdown at N, 10, 5 and 1 minutes before running, 0-60 (default 0)"`
	Message     string `json:"message,omitempty" jsonschema:"Countdown broadcast text; {minutes} is replaced with the minutes left"`
	Command     string `json:"command,omitempty" jsonschema:"RCON command to run (action=command only)"`
	Keep        int    `json:"keep,omitempty" jsonschema:"For backup: keep only the latest N scheduled backups, 0-1000 (0 keeps all)"`
}

type updateScheduleInput struct {
	ServerID    uint    `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	ScheduleID  uint    `json:"schedule_id" jsonschema:"Schedule ID (from list_schedules)"`
	Action      *string `json:"action,omitempty" jsonschema:"One of: restart, stop, start, backup, update, command"`
	Cron        *string `json:"cron,omitempty" jsonschema:"5-field cron expression or @daily / @hourly"`
	Enabled     *bool   `json:"enabled,omitempty" jsonschema:"Enable or disable the schedule"`
	WarnMinutes *int    `json:"warn_minutes,omitempty" jsonschema:"Countdown minutes, 0-60"`
	Message     *string `json:"message,omitempty" jsonschema:"Countdown broadcast text; {minutes} is replaced with the minutes left"`
	Command     *string `json:"command,omitempty" jsonschema:"RCON command (action=command only)"`
	Keep        *int    `json:"keep,omitempty" jsonschema:"For backup: number of scheduled backups to keep, 0 keeps all"`
}

func registerOperationTools(r *registry) {
	// ---------- 日志、RCON、玩家、Mod ----------

	query(r, "get_logs", "Get container logs",
		"Get the last lines of the server container's output: download and mod install progress, startup messages and errors. Empty if the container has not been created yet.",
		func(_ context.Context, userID uint, in logsInput) (any, error) {
			lines, err := service.GetLogs(userID, sid(in.ServerID), clampTail(in.Tail, 200))
			if err != nil {
				return nil, err
			}
			return map[string]any{"lines": lines}, nil
		})

	query(r, "get_game_log", "Get game log",
		"Get the last lines of ShooterGame.log (the game's own log: player joins, tames, deaths, chat when enabled, mod errors). Available even when the server is stopped.",
		func(_ context.Context, userID uint, in logsInput) (any, error) {
			return service.GetGameLog(userID, sid(in.ServerID), clampTail(in.Tail, 200))
		})

	query(r, "list_players", "List online players",
		"List players currently online (name and Steam ID) via RCON ListPlayers. The server must be running and ready.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			players, err := service.ListPlayers(userID, sid(in.ServerID))
			if err != nil {
				return nil, err
			}
			return map[string]any{"players": players}, nil
		})

	query(r, "list_server_mods", "List server mods",
		"List the server's configured mods with Workshop title, install status and whether an update is available on Steam Workshop. Update outdated mods with update_game (mod_ids).",
		func(_ context.Context, userID uint, in serverModsInput) (any, error) {
			mods, err := service.GetServerMods(userID, sid(in.ServerID), in.Refresh)
			if err != nil {
				return nil, err
			}
			return map[string]any{"mods": mods}, nil
		})

	action(r, destructive, "rcon_command", "Run RCON command",
		"Run one RCON command on a running server and return its output. Commands can change the game world (e.g. DestroyWildDinos, KillPlayer), so use with care. Common commands: ServerChat <msg>, Broadcast <msg>, SaveWorld, ListPlayers, KickPlayer <SteamID>, BanPlayer <SteamID>, UnbanPlayer <SteamID>, SetTimeOfDay <hh:mm>.",
		func(_ context.Context, userID uint, in rconInput) (any, error) {
			out, err := service.RconExec(userID, sid(in.ServerID), in.Command)
			if err != nil {
				return nil, err
			}
			return map[string]any{"output": out}, nil
		})

	// ---------- 备份 ----------

	query(r, "list_backups", "List backups",
		"List save backups of a server, newest first. kind is manual, scheduled, pre-restore or uploaded.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			list, err := service.ListBackups(userID, sid(in.ServerID))
			if err != nil {
				return nil, err
			}
			return map[string]any{"backups": list}, nil
		})

	action(r, additive, "create_backup", "Create backup",
		"Back up the server's saves and config (SaveWorld first if running). Can take a while for large saves.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			return service.CreateBackup(userID, sid(in.ServerID))
		})

	action(r, destructive, "restore_backup", "Restore backup",
		"Replace the server's saves and config with a backup. The server must be stopped. A pre-restore backup of the current state is made automatically first.",
		func(_ context.Context, userID uint, in backupInput) (any, error) {
			if err := service.RestoreBackup(userID, sid(in.ServerID), in.Name); err != nil {
				return nil, err
			}
			return message("Backup restored. Start the server to use it."), nil
		})

	action(r, destructive, "delete_backup", "Delete backup",
		"Permanently delete a backup file.",
		func(_ context.Context, userID uint, in backupInput) (any, error) {
			if err := service.DeleteBackup(userID, sid(in.ServerID), in.Name); err != nil {
				return nil, err
			}
			return message("Backup deleted."), nil
		})

	// ---------- 定时任务 ----------

	query(r, "list_schedules", "List schedules",
		"List the server's scheduled tasks with next run time and last result.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			list, err := service.ListSchedules(userID, sid(in.ServerID))
			if err != nil {
				return nil, err
			}
			return map[string]any{"schedules": list}, nil
		})

	action(r, additive, "create_schedule", "Create schedule",
		"Create a scheduled task for a server (restart, stop, start, backup, update or RCON command) on a cron schedule. Scheduled restarts only affect running servers.",
		func(_ context.Context, userID uint, in createScheduleInput) (any, error) {
			req := models.ScheduleRequest{
				Action: in.Action, Cron: in.Cron, Enabled: in.Enabled, WarnMinutes: in.WarnMinutes,
				Message: in.Message, Command: in.Command, Keep: in.Keep,
			}
			if err := binding.Validator.ValidateStruct(&req); err != nil {
				return nil, fmt.Errorf("[bad_request] 参数无效: %w", err)
			}
			return service.CreateSchedule(userID, sid(in.ServerID), req)
		})

	action(r, idempotent, "update_schedule", "Update schedule",
		"Change a scheduled task. Only the fields you pass are changed.",
		func(_ context.Context, userID uint, in updateScheduleInput) (any, error) {
			list, err := service.ListSchedules(userID, sid(in.ServerID))
			if err != nil {
				return nil, err
			}
			var current *models.Schedule
			for i := range list {
				if list[i].ID == in.ScheduleID {
					current = &list[i]
					break
				}
			}
			if current == nil {
				return nil, server.ErrScheduleNotFound
			}

			req := models.ScheduleRequest{
				Action: current.Action, Cron: current.Cron, Enabled: in.Enabled, WarnMinutes: current.WarnMinutes,
				Message: current.Message, Command: current.Command, Keep: current.Keep,
			}
			setIf(&req.Action, in.Action)
			setIf(&req.Cron, in.Cron)
			setIf(&req.WarnMinutes, in.WarnMinutes)
			setIf(&req.Message, in.Message)
			setIf(&req.Command, in.Command)
			setIf(&req.Keep, in.Keep)
			if err := binding.Validator.ValidateStruct(&req); err != nil {
				return nil, fmt.Errorf("[bad_request] 参数无效: %w", err)
			}
			return service.UpdateSchedule(userID, sid(in.ServerID), sid(in.ScheduleID), req)
		})

	action(r, destructive, "delete_schedule", "Delete schedule",
		"Delete a scheduled task.",
		func(_ context.Context, userID uint, in scheduleIDInput) (any, error) {
			if err := service.DeleteSchedule(userID, sid(in.ServerID), sid(in.ScheduleID)); err != nil {
				return nil, err
			}
			return message("Schedule deleted."), nil
		})

	action(r, additive, "run_schedule", "Run schedule now",
		"Run a scheduled task immediately in the background (countdown broadcasts included). Check list_schedules for last_result afterwards.",
		func(_ context.Context, userID uint, in scheduleIDInput) (any, error) {
			if err := service.RunScheduleNow(userID, sid(in.ServerID), sid(in.ScheduleID)); err != nil {
				return nil, err
			}
			return message("Schedule started."), nil
		})
}

func clampTail(tail, def int) int {
	if tail <= 0 {
		return def
	}
	return min(tail, 5000)
}

func setIf[T any](dst *T, src *T) {
	if src != nil {
		*dst = *src
	}
}

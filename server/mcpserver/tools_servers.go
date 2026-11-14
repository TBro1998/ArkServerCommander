package mcpserver

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"ark-server-commander/models"
	"ark-server-commander/service/server"

	"github.com/gin-gonic/gin/binding"
)

// ---------- 输入 ----------

type serverInput struct {
	ServerID uint `json:"server_id" jsonschema:"Server ID (from list_servers)"`
}

type listServersInput struct {
	ServerID uint `json:"server_id,omitempty" jsonschema:"Only return this server. Omit to list all servers"`
}

// serverArgsInput 启动参数，三项都可省略
type serverArgsInput struct {
	QueryParams     map[string]string `json:"query_params,omitempty" jsonschema:"Extra ?Key=Value map URL options (e.g. AllowCrateSpawnsOnTopOfStructures=True). Port, QueryPort, MaxPlayers, RCONPort, SessionName and GameModIds are managed by their own fields and ignored here. Empty or False values are omitted"`
	CommandLineArgs map[string]any    `json:"command_line_args,omitempty" jsonschema:"-Flag arguments: true adds -Flag, false/empty/0 omits it, other values add -Flag=Value (e.g. {\"NoBattlEye\": true, \"MinimumTimeBetweenInventoryRetrieval\": 3600})"`
	CustomArgs      []string          `json:"custom_args,omitempty" jsonschema:"Additional raw arguments appended at the end"`
}

func (a *serverArgsInput) toRequest() *models.ServerArgsRequest {
	if a == nil {
		return nil
	}
	req := &models.ServerArgsRequest{QueryParams: a.QueryParams, CommandLineArgs: a.CommandLineArgs, CustomArgs: a.CustomArgs}
	if req.QueryParams == nil {
		req.QueryParams = map[string]string{}
	}
	if req.CommandLineArgs == nil {
		req.CommandLineArgs = map[string]any{}
	}
	if req.CustomArgs == nil {
		req.CustomArgs = []string{}
	}
	return req
}

type createServerInput struct {
	Identifier       string           `json:"identifier" jsonschema:"Unique short identifier for the server, e.g. island-pve"`
	SessionName      string           `json:"session_name,omitempty" jsonschema:"Name shown in the in-game server browser"`
	Map              string           `json:"map,omitempty" jsonschema:"Map name: TheIsland (default), TheCenter, ScorchedEarth_P, Aberration_P, Extinction, Ragnarok, Valguero_P, Genesis, Genesis2, CrystalIsles, LostIsland, Fjordur, or a mod map name"`
	MaxPlayers       int              `json:"max_players,omitempty" jsonschema:"Player slots, 1-200 (default 70)"`
	Port             int              `json:"port,omitempty" jsonschema:"Game port (UDP, port+1 is also used). Omit to pick free ports automatically"`
	QueryPort        int              `json:"query_port,omitempty" jsonschema:"Steam query port. Omit to pick automatically"`
	RCONPort         int              `json:"rcon_port,omitempty" jsonschema:"RCON port. Omit to pick automatically"`
	AdminPassword    string           `json:"admin_password,omitempty" jsonschema:"Server admin / RCON password. Omit to generate a random one (returned in the result)"`
	ClusterID        string           `json:"cluster_id,omitempty" jsonschema:"Cluster ID; servers with the same cluster ID share character and item transfers"`
	GameModIDs       string           `json:"game_mod_ids,omitempty" jsonschema:"Comma-separated Steam Workshop mod IDs in load order"`
	AutoRestart      *bool            `json:"auto_restart,omitempty" jsonschema:"Restart the container automatically if it crashes (default true)"`
	AutoUpdate       *bool            `json:"auto_update,omitempty" jsonschema:"Check for game and mod updates on every start (default false)"`
	ArkApi           *bool            `json:"arkapi,omitempty" jsonschema:"Enable ArkApi plugins. Runs the Windows server under Proton instead of the native Linux server (default false)"`
	GameUserSettings string           `json:"game_user_settings,omitempty" jsonschema:"Full GameUserSettings.ini content. Omit to use the default template"`
	GameIni          string           `json:"game_ini,omitempty" jsonschema:"Full Game.ini content. Omit to use the default template"`
	ServerArgs       *serverArgsInput `json:"server_args,omitempty" jsonschema:"Launch arguments. Omit to use the recommended defaults; if given, it replaces the defaults entirely (see server_args of an existing server via get_server)"`
}

type updateServerInput struct {
	ServerID         uint             `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Identifier       string           `json:"identifier,omitempty" jsonschema:"New identifier"`
	SessionName      string           `json:"session_name,omitempty" jsonschema:"New name shown in the server browser"`
	Map              string           `json:"map,omitempty" jsonschema:"New map name"`
	MaxPlayers       int              `json:"max_players,omitempty" jsonschema:"Player slots, 1-200"`
	Port             int              `json:"port,omitempty" jsonschema:"Game port"`
	QueryPort        int              `json:"query_port,omitempty" jsonschema:"Steam query port"`
	RCONPort         int              `json:"rcon_port,omitempty" jsonschema:"RCON port"`
	AdminPassword    string           `json:"admin_password,omitempty" jsonschema:"New admin / RCON password (applied on next start)"`
	ClusterID        *string          `json:"cluster_id,omitempty" jsonschema:"Cluster ID; empty string removes the server from its cluster"`
	GameModIDs       *string          `json:"game_mod_ids,omitempty" jsonschema:"Comma-separated Steam Workshop mod IDs in load order, replaces the whole list; empty string removes all mods"`
	AutoRestart      *bool            `json:"auto_restart,omitempty" jsonschema:"Restart the container automatically if it crashes"`
	AutoUpdate       *bool            `json:"auto_update,omitempty" jsonschema:"Check for game and mod updates on every start"`
	ArkApi           *bool            `json:"arkapi,omitempty" jsonschema:"Enable ArkApi (switches the server image; the container is recreated on next start)"`
	GameUserSettings string           `json:"game_user_settings,omitempty" jsonschema:"Full new GameUserSettings.ini content (replaces the file). To change a few keys use set_config_values instead"`
	GameIni          string           `json:"game_ini,omitempty" jsonschema:"Full new Game.ini content (replaces the file). To change a few keys use set_config_values instead"`
	ServerArgs       *serverArgsInput `json:"server_args,omitempty" jsonschema:"Launch arguments; replaces all existing launch arguments, so read them with get_server first"`
}

type setConfigValuesInput struct {
	ServerID uint              `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	File     string            `json:"file" jsonschema:"GameUserSettings.ini or Game.ini"`
	Section  string            `json:"section" jsonschema:"INI section without brackets, e.g. ServerSettings (GameUserSettings.ini) or /script/shootergame.shootergamemode (Game.ini)"`
	Values   map[string]string `json:"values" jsonschema:"Keys to set and their values, e.g. {\"XPMultiplier\": \"2.0\", \"TamingSpeedMultiplier\": \"3.0\"}"`
}

type deleteServerInput struct {
	ServerID   uint `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	DeleteData bool `json:"delete_data,omitempty" jsonschema:"Also delete the server's files, saves and backups permanently. Default false keeps them on disk"`
}

type updateGameInput struct {
	ServerID uint     `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Game     bool     `json:"game,omitempty" jsonschema:"Validate and update the game server files"`
	Mods     bool     `json:"mods,omitempty" jsonschema:"Re-download all mods"`
	ModIDs   []string `json:"mod_ids,omitempty" jsonschema:"Only re-download these mods (used when mods is false)"`
}

// ---------- 工具 ----------

func registerServerTools(r *registry) {
	query(r, "list_servers", "List servers",
		"List all ARK servers with status (stopped/starting/running/stopping), ports, map, mods, last_error, and for running servers the live runtime info (phase, players, CPU, memory). Cheap enough to poll.",
		func(_ context.Context, userID uint, in listServersInput) (any, error) {
			list, err := service.GetServers(userID)
			if err != nil {
				return nil, err
			}
			if in.ServerID == 0 {
				return map[string]any{"servers": list}, nil
			}
			for _, s := range list {
				if s.ID == in.ServerID {
					return map[string]any{"servers": []models.ServerResponse{s}}, nil
				}
			}
			return nil, server.ErrServerNotFound
		})

	query(r, "get_server", "Get server details",
		"Get one server's full configuration: settings, GameUserSettings.ini and Game.ini contents (pending edits if any), launch arguments and the generated command line. Reads files from disk, so prefer list_servers for status polling.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			return service.GetServer(userID, sid(in.ServerID))
		})

	secret(r, "get_admin_password", "Get admin password",
		"Get the server's admin / RCON password and RCON port, e.g. to give an admin the in-game password.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			return service.GetServerRCON(userID, sid(in.ServerID))
		})

	action(r, additive, "create_server", "Create server",
		"Create a new ARK server. Ports and admin password are filled in automatically when omitted. The container is created on first start; call start_server afterwards. The first start downloads about 20 GB unless another installed server can be copied.",
		func(_ context.Context, userID uint, in createServerInput) (any, error) {
			req := models.ServerRequest{
				Identifier:       in.Identifier,
				SessionName:      in.SessionName,
				ClusterID:        in.ClusterID,
				Port:             in.Port,
				QueryPort:        in.QueryPort,
				RCONPort:         in.RCONPort,
				AdminPassword:    in.AdminPassword,
				Map:              in.Map,
				MaxPlayers:       in.MaxPlayers,
				GameModIds:       in.GameModIDs,
				AutoRestart:      in.AutoRestart,
				AutoUpdate:       in.AutoUpdate,
				ArkApi:           in.ArkApi,
				GameUserSettings: in.GameUserSettings,
				GameIni:          in.GameIni,
				ServerArgs:       in.ServerArgs.toRequest(),
			}
			if req.SessionName == "" {
				req.SessionName = req.Identifier
			}
			if req.MaxPlayers == 0 {
				req.MaxPlayers = 70
			}
			if req.Port == 0 || req.QueryPort == 0 || req.RCONPort == 0 {
				p := service.SuggestPorts()
				if req.Port == 0 {
					req.Port = p.Port
				}
				if req.QueryPort == 0 {
					req.QueryPort = p.QueryPort
				}
				if req.RCONPort == 0 {
					req.RCONPort = p.RCONPort
				}
			}
			generated := req.AdminPassword == ""
			if generated {
				pw, err := randomPassword(16)
				if err != nil {
					return nil, err
				}
				req.AdminPassword = pw
			}
			if err := binding.Validator.ValidateStruct(&req); err != nil {
				return nil, fmt.Errorf("[bad_request] 参数无效: %w", err)
			}

			created, err := service.CreateServer(userID, req)
			if err != nil {
				return nil, err
			}
			out := map[string]any{"server": created}
			if generated {
				out["admin_password"] = req.AdminPassword
			}
			return out, nil
		})

	action(r, idempotent, "update_server", "Update server settings",
		"Change server settings. Only the fields you pass are changed. Changes to ports, map, mods, name or launch arguments take effect after a restart; config file edits on a running server are stored as pending and applied on next start.",
		func(_ context.Context, userID uint, in updateServerInput) (any, error) {
			req := models.ServerUpdateRequest{
				Identifier:       in.Identifier,
				SessionName:      in.SessionName,
				ClusterID:        in.ClusterID,
				Port:             in.Port,
				QueryPort:        in.QueryPort,
				RCONPort:         in.RCONPort,
				AdminPassword:    in.AdminPassword,
				Map:              in.Map,
				MaxPlayers:       in.MaxPlayers,
				GameModIds:       in.GameModIDs,
				AutoRestart:      in.AutoRestart,
				AutoUpdate:       in.AutoUpdate,
				ArkApi:           in.ArkApi,
				GameUserSettings: in.GameUserSettings,
				GameIni:          in.GameIni,
				ServerArgs:       in.ServerArgs.toRequest(),
			}
			if err := binding.Validator.ValidateStruct(&req); err != nil {
				return nil, fmt.Errorf("[bad_request] 参数无效: %w", err)
			}
			updated, argsChanged, err := service.UpdateServer(userID, sid(in.ServerID), req)
			if err != nil {
				return nil, err
			}
			// 响应里的配置文件内容很长，Agent 刚提交过，不再原样返回
			updated.GameUserSettings, updated.GameIni = "", ""
			out := map[string]any{"server": updated, "args_changed": argsChanged}
			if argsChanged && updated.Status == "running" {
				out["message"] = "Launch arguments changed; restart the server to apply them."
			}
			return out, nil
		})

	action(r, idempotent, "set_config_values", "Set config values",
		"Set individual keys in GameUserSettings.ini or Game.ini without touching the rest of the file (comments, unknown keys and other sections are kept). Duplicate keys are merged into one line, so use update_server with full content for keys that must repeat (e.g. OverrideNamedEngramEntries). Server rates, difficulty and day/night speed live in GameUserSettings.ini [ServerSettings]. Edits on a running server are applied on next start.",
		func(_ context.Context, userID uint, in setConfigValuesInput) (any, error) {
			updated, err := service.SetConfigValues(userID, sid(in.ServerID), in.File, in.Section, in.Values)
			if err != nil {
				return nil, err
			}
			return map[string]any{
				"file":           in.File,
				"section":        in.Section,
				"values":         in.Values,
				"config_pending": updated.ConfigPending,
				"status":         updated.Status,
			}, nil
		})

	action(r, destructive, "delete_server", "Delete server",
		"Delete a stopped server and its container. By default files, saves and backups stay on disk; with delete_data=true they are deleted permanently (shared cluster data is kept).",
		func(_ context.Context, userID uint, in deleteServerInput) (any, error) {
			if err := service.DeleteServer(userID, sid(in.ServerID), in.DeleteData); err != nil {
				return nil, err
			}
			return message("Server deleted."), nil
		})

	action(r, additive, "start_server", "Start server",
		"Start a stopped server. Asynchronous: poll list_servers until status is running and runtime.phase is ready. Recreates the container automatically if its configuration or image is out of date.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			if err := service.StartServer(userID, sid(in.ServerID)); err != nil {
				return nil, err
			}
			return message("Starting. Poll list_servers for status and runtime.phase."), nil
		})

	action(r, additive, "stop_server", "Stop server",
		"Stop a running server. Saves the world via RCON first. Asynchronous: poll list_servers until status is stopped.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			if err := service.StopServer(userID, sid(in.ServerID)); err != nil {
				return nil, err
			}
			return message("Stopping. Poll list_servers until status is stopped."), nil
		})

	action(r, additive, "restart_server", "Restart server",
		"Restart a server (stop if running, then start). Players are disconnected without warning; broadcast a message with rcon_command first if players are online. Asynchronous.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			if err := service.RestartServer(userID, sid(in.ServerID)); err != nil {
				return nil, err
			}
			return message("Restarting. Poll list_servers for status."), nil
		})

	action(r, additive, "recreate_container", "Recreate container",
		"Recreate the server's Docker container from the current configuration (stops it first if running). The server ends up stopped. Usually unnecessary: start_server recreates the container when needed. Asynchronous.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			if err := service.RecreateContainer(userID, sid(in.ServerID)); err != nil {
				return nil, err
			}
			return message("Recreating container. Poll list_servers until status is stopped."), nil
		})

	action(r, additive, "update_game", "Update game / mods",
		"Update the game server files and/or mods, then start the server (stopping it first if running). Use when players get a version mismatch or list_servers shows mods_outdated. Asynchronous.",
		func(_ context.Context, userID uint, in updateGameInput) (any, error) {
			if err := service.UpdateGame(userID, sid(in.ServerID), in.Game, in.Mods, in.ModIDs); err != nil {
				return nil, err
			}
			return message("Update started; the server will start afterwards. Poll list_servers for status."), nil
		})
}

func randomPassword(n int) (string, error) {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	b := make([]byte, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		b[i] = chars[idx.Int64()]
	}
	return string(b), nil
}

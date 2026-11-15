package mcpserver

import (
	"context"

	"ark-server-commander/models"
)

const configPartsDesc = "Parts to apply: config_files (GameUserSettings.ini and Game.ini; the target keeps its own SessionName, ports, MaxPlayers and ServerAdminPassword lines), server_args (launch arguments), mods (mod list), map, cluster (cluster ID), options (max_players, arkapi, auto_restart, auto_update). Omit to apply all parts"

type importConfigInput struct {
	ServerID uint                `json:"server_id" jsonschema:"Target server ID (from list_servers)"`
	Config   models.ServerConfig `json:"config" jsonschema:"Server config object exactly as returned by export_server_config (format and version are required); you may edit its fields before importing"`
	Parts    []string            `json:"parts,omitempty" jsonschema:"Parts to apply: config_files, server_args, mods, map, cluster, options. Omit to apply all parts"`
}

type copyConfigInput struct {
	ServerID  uint     `json:"server_id" jsonschema:"Source server ID (from list_servers)"`
	TargetIDs []uint   `json:"target_ids" jsonschema:"IDs of the servers to copy the config to; the source server is skipped"`
	Parts     []string `json:"parts,omitempty" jsonschema:"Parts to apply: config_files, server_args, mods, map, cluster, options. Omit to apply all parts"`
}

func registerConfigTools(r *registry) {
	query(r, "export_server_config", "Export server config",
		"Export a server's portable configuration: map, mods, cluster ID, options, GameUserSettings.ini, Game.ini and launch arguments. Ports and the admin password are not included (ServerAdminPassword is removed from GameUserSettings.ini). Pass the result to import_server_config to apply it to another server.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			return service.ExportConfig(userID, sid(in.ServerID))
		})

	action(r, idempotent, "import_server_config", "Import server config",
		"Apply a config from export_server_config to an existing server. "+configPartsDesc+". Name, ports and admin password stay unchanged. On a server that is not stopped, config files are stored as pending and everything takes effect after a restart.",
		func(_ context.Context, userID uint, in importConfigInput) (any, error) {
			updated, err := service.ImportConfig(userID, sid(in.ServerID), &in.Config, in.Parts)
			if err != nil {
				return nil, err
			}
			updated.GameUserSettings, updated.GameIni = "", ""
			return map[string]any{"server": updated}, nil
		})

	action(r, idempotent, "copy_server_config", "Copy config to other servers",
		"Copy one server's configuration to other servers in one call. "+configPartsDesc+". Each target is updated independently; check every entry in results for success and error. Targets that are not stopped get config_pending=true and need a restart.",
		func(_ context.Context, userID uint, in copyConfigInput) (any, error) {
			results, err := service.CopyConfig(userID, sid(in.ServerID), models.CopyConfigRequest{TargetIDs: in.TargetIDs, Parts: in.Parts})
			if err != nil {
				return nil, err
			}
			return map[string]any{"results": results}, nil
		})
}

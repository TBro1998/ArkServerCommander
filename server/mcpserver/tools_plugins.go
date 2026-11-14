package mcpserver

import (
	"context"
)

type pluginInput struct {
	ServerID uint   `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Name     string `json:"name" jsonschema:"Plugin name (from list_plugins)"`
}

type pluginStateInput struct {
	ServerID uint   `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Name     string `json:"name" jsonschema:"Plugin name (from list_plugins)"`
	Enabled  bool   `json:"enabled" jsonschema:"true to enable, false to disable"`
}

type pluginConfigInput struct {
	ServerID uint   `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Plugin   string `json:"plugin,omitempty" jsonschema:"Plugin name (from list_plugins). Omit for ArkApi's own config.json"`
}

type savePluginConfigInput struct {
	ServerID uint   `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Plugin   string `json:"plugin,omitempty" jsonschema:"Plugin name (from list_plugins). Omit for ArkApi's own config.json"`
	Content  string `json:"content" jsonschema:"Full config.json content; must be valid JSON"`
}

type installArkApiInput struct {
	ServerID uint   `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	Version  string `json:"version,omitempty" jsonschema:"ArkApi version (from list_arkapi_versions). Omit for the latest"`
}

type arkApiVersionsInput struct {
	Refresh bool `json:"refresh,omitempty" jsonschema:"Query GitHub again instead of using the cache"`
}

type arkApiVersionInput struct {
	Version string `json:"version" jsonschema:"Downloaded ArkApi version to delete"`
}

func registerPluginTools(r *registry) {
	query(r, "list_plugins", "List ArkApi plugins",
		"List the server's ArkApi plugins (enabled and disabled) and whether / which ArkApi version is installed. Plugins only load on servers with arkapi enabled. Installing new plugins requires uploading a zip in the web UI.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			return service.ListPlugins(userID, sid(in.ServerID))
		})

	action(r, idempotent, "set_plugin_enabled", "Enable / disable plugin",
		"Enable or disable an ArkApi plugin. The server must be stopped.",
		func(_ context.Context, userID uint, in pluginStateInput) (any, error) {
			if err := service.SetPluginEnabled(userID, sid(in.ServerID), in.Name, in.Enabled); err != nil {
				return nil, err
			}
			return message("Plugin state changed; takes effect on next start."), nil
		})

	action(r, destructive, "delete_plugin", "Delete plugin",
		"Permanently delete an ArkApi plugin and its files, including its config. The server must be stopped.",
		func(_ context.Context, userID uint, in pluginInput) (any, error) {
			if err := service.DeletePlugin(userID, sid(in.ServerID), in.Name); err != nil {
				return nil, err
			}
			return message("Plugin deleted."), nil
		})

	query(r, "get_plugin_config", "Get plugin config",
		"Read a plugin's config.json, or ArkApi's own config.json when plugin is omitted. Empty if the file does not exist.",
		func(_ context.Context, userID uint, in pluginConfigInput) (any, error) {
			content, err := service.GetPluginConfig(userID, sid(in.ServerID), in.Plugin)
			if err != nil {
				return nil, err
			}
			return map[string]any{"content": content}, nil
		})

	action(r, idempotent, "save_plugin_config", "Save plugin config",
		"Overwrite a plugin's config.json (or ArkApi's own when plugin is omitted). Allowed while running; takes effect when the plugin reloads or the server restarts.",
		func(_ context.Context, userID uint, in savePluginConfigInput) (any, error) {
			if err := service.SavePluginConfig(userID, sid(in.ServerID), in.Plugin, in.Content); err != nil {
				return nil, err
			}
			return message("Config saved."), nil
		})

	action(r, idempotent, "install_arkapi", "Install ArkApi",
		"Install or switch the ArkApi version in the server files (keeps the existing ArkApi config.json). The server must be stopped. May download from GitHub.",
		func(_ context.Context, userID uint, in installArkApiInput) (any, error) {
			version, err := service.InstallArkApi(userID, sid(in.ServerID), in.Version)
			if err != nil {
				return nil, err
			}
			return map[string]any{"ok": true, "version": version}, nil
		})

	query(r, "list_arkapi_versions", "List ArkApi versions",
		"List ArkApi versions available to install: built-in, already downloaded, and releases on GitHub not yet downloaded.",
		func(_ context.Context, _ uint, in arkApiVersionsInput) (any, error) {
			return service.ListArkApiVersions(in.Refresh), nil
		})

	action(r, destructive, "delete_arkapi_version", "Delete ArkApi version",
		"Delete a downloaded or uploaded ArkApi version from the local library (the built-in version cannot be deleted). Servers that already installed it are unaffected.",
		func(_ context.Context, _ uint, in arkApiVersionInput) (any, error) {
			if err := service.DeleteArkApiVersion(in.Version); err != nil {
				return nil, err
			}
			return message("Version deleted."), nil
		})
}

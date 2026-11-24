package mcpserver

import (
	"context"
	"strings"
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

type pluginCatalogInput struct {
	Search          string `json:"search,omitempty" jsonschema:"Only entries whose title, tagline, author or plugin names contain this text (case-insensitive)"`
	Category        string `json:"category,omitempty" jsonschema:"official, free or paid. Omit for all"`
	InstallableOnly bool   `json:"installable_only,omitempty" jsonschema:"Only entries that can be installed with install_catalog_plugin"`
	Refresh         bool   `json:"refresh,omitempty" jsonschema:"Fetch the catalog again instead of using the 1-hour cache"`
}

type installCatalogPluginInput struct {
	ServerID   uint `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	ResourceID uint `json:"resource_id" jsonschema:"Catalog entry id (from list_plugin_catalog)"`
}

type pluginLogInput struct {
	ServerID uint   `json:"server_id" jsonschema:"Server ID (from list_servers)"`
	File     string `json:"file,omitempty" jsonschema:"Log file name (from the files list). Omit to only list the files"`
	Tail     int    `json:"tail,omitempty" jsonschema:"Number of lines from the end, 1-5000 (default 200)"`
}

type arkApiVersionInput struct {
	Version string `json:"version" jsonschema:"Downloaded ArkApi version to delete"`
}

func registerPluginTools(r *registry) {
	query(r, "list_plugins", "List ArkApi plugins",
		"List the server's ArkApi plugins (enabled and disabled) and whether / which ArkApi version is installed. Plugins only load on servers with arkapi enabled. Plugins installed from the catalog have a source with the catalog resource_id and version. New plugins can be installed with install_catalog_plugin; uploading a zip requires the web UI.",
		func(_ context.Context, userID uint, in serverInput) (any, error) {
			return service.ListPlugins(userID, sid(in.ServerID))
		})

	query(r, "get_plugin_logs", "Get ArkApi / plugin logs",
		"Without file: list the .log files ArkApi and its plugins write to ShooterGame/Binaries/Win64/logs (newest first, at most 200; total is the full count). With file: get the last lines of that log file. Available even when the server is stopped.",
		func(_ context.Context, userID uint, in pluginLogInput) (any, error) {
			if in.File == "" {
				return service.ListPluginLogs(userID, sid(in.ServerID))
			}
			return service.GetPluginLog(userID, sid(in.ServerID), in.File, clampTail(in.Tail, 200))
		})

	query(r, "list_plugin_catalog", "List plugin catalog",
		"List ARK: Survival Evolved ArkApi plugins published on ark-server-api.com, as indexed by arkfgf.com (no account binding needed). Entries with download_url can be installed with install_catalog_plugin; paid ones and those that must be downloaded on the site cannot. plugins lists the plugin folder names inside the archive with their version and dependencies; compare with list_plugins to find updates. Fails when the manager has ARKFGF_API disabled.",
		func(_ context.Context, _ uint, in pluginCatalogInput) (any, error) {
			catalog, err := service.GetPluginCatalog(in.Refresh)
			if err != nil {
				return nil, err
			}
			search := strings.ToLower(strings.TrimSpace(in.Search))
			entries := []any{}
			for _, p := range catalog.Plugins {
				if in.Category != "" && p.Category != in.Category {
					continue
				}
				if in.InstallableOnly && p.DownloadURL == "" {
					continue
				}
				if search != "" {
					text := strings.ToLower(p.Title + " " + p.Tagline + " " + p.Author)
					for _, info := range p.Plugins {
						text += " " + strings.ToLower(info.Name)
					}
					if !strings.Contains(text, search) {
						continue
					}
				}
				entries = append(entries, p)
			}
			return map[string]any{"source": catalog.Source, "generated_at": catalog.GeneratedAt, "plugins": entries}, nil
		})

	action(r, idempotent, "install_catalog_plugin", "Install plugin from catalog",
		"Download a catalog entry (from list_plugin_catalog) from ark-server-api.com and install it; if already installed it is updated, keeping its config.json and enabled state. The server must be stopped. Dependencies are not installed automatically: check the entry's plugins[].dependencies.",
		func(_ context.Context, userID uint, in installCatalogPluginInput) (any, error) {
			names, err := service.InstallCatalogPlugin(userID, sid(in.ServerID), in.ResourceID)
			if err != nil {
				return nil, err
			}
			return map[string]any{"ok": true, "installed": names}, nil
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
		"List ArkApi versions available to install: built-in, already downloaded, releases on GitHub, and the version in the arkfgf.com plugin catalog (on_catalog, downloaded from ark-server-api.com on install).",
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

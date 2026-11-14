package mcpserver

import (
	"context"

	"ark-server-commander/service/server"
)

type overviewInput struct {
	RefreshDisk bool `json:"refresh_disk,omitempty" jsonschema:"Recalculate disk usage instead of using the cache (slow)"`
}

type modInfoInput struct {
	IDs string `json:"ids" jsonschema:"Comma-separated Steam Workshop mod IDs (up to 100)"`
}

type clusterInput struct {
	ClusterID string `json:"cluster_id" jsonschema:"Cluster ID (from get_system_overview unused_clusters)"`
}

const imageNames = "tbro98/ase-server:linux (native Linux server, used when arkapi is off) or tbro98/ase-server:proton (Windows server under Proton, used when arkapi is on)"

type imageInput struct {
	ImageName string `json:"image_name" jsonschema:"One of: tbro98/ase-server:linux (native Linux server, used when arkapi is off) or tbro98/ase-server:proton (Windows server under Proton, used when arkapi is on)"`
}

type updateImageInput struct {
	ImageName      string `json:"image_name" jsonschema:"One of: tbro98/ase-server:linux or tbro98/ase-server:proton"`
	RestartRunning bool   `json:"restart_running,omitempty" jsonschema:"Also restart running servers that use this image right away (players are disconnected). Default false: they switch on their next start"`
}

func registerSystemTools(r *registry) {
	query(r, "get_system_overview", "System overview",
		"Host overview: public address, server counts, players online, disk usage per server and unused cluster data.",
		func(_ context.Context, userID uint, in overviewInput) (any, error) {
			return service.GetOverview(userID, in.RefreshDisk)
		})

	query(r, "get_system_info", "System info",
		"Public address players use to connect (public_host:port) and the manager's version info.",
		func(_ context.Context, _ uint, _ noInput) (any, error) {
			return map[string]any{
				"public_host": server.PublicHost(),
				"version":     service.GetVersionInfo(false),
			}, nil
		})

	query(r, "suggest_ports", "Suggest free ports",
		"Suggest a game / query / RCON port set that no existing server uses.",
		func(_ context.Context, _ uint, _ noInput) (any, error) {
			return service.SuggestPorts(), nil
		})

	query(r, "get_mod_info", "Look up Workshop mods",
		"Look up Steam Workshop mods by ID: title, size, last update time. found=false means the ID does not exist; checked=false means Steam could not be reached.",
		func(_ context.Context, _ uint, in modInfoInput) (any, error) {
			mods, err := service.GetModInfo(in.IDs)
			if err != nil {
				return nil, err
			}
			return map[string]any{"mods": mods}, nil
		})

	action(r, destructive, "delete_cluster", "Delete unused cluster data",
		"Permanently delete an unused cluster's shared data (transferred characters and items). Refused while any server still uses the cluster ID.",
		func(_ context.Context, _ uint, in clusterInput) (any, error) {
			if err := service.DeleteCluster(in.ClusterID); err != nil {
				return nil, err
			}
			return message("Cluster data deleted."), nil
		})

	// ---------- 镜像 ----------

	query(r, "get_image_status", "Docker image status",
		"Status of the Docker images the manager uses: downloaded or not, pull progress, whether an update was found, and whether each image is required by current servers.",
		func(_ context.Context, _ uint, _ noInput) (any, error) {
			return service.GetImageStatus()
		})

	query(r, "check_image_updates", "Check image updates",
		"Compare local images with the registry (no layers are downloaded) and report which have updates. Results also show up in get_image_status.",
		func(_ context.Context, _ uint, _ noInput) (any, error) {
			updates, err := service.CheckImageUpdates()
			if err != nil {
				return nil, err
			}
			return map[string]any{"has_update": updates}, nil
		})

	action(r, additive, "pull_image", "Pull image",
		"Download a missing image in the background. Allowed images: "+imageNames+". Poll get_image_status for progress. Images are also pulled automatically when a server needs them.",
		func(_ context.Context, _ uint, in imageInput) (any, error) {
			if err := service.PullImage(in.ImageName); err != nil {
				return nil, err
			}
			return message("Pull started. Poll get_image_status for progress."), nil
		})

	action(r, additive, "update_image", "Update image",
		"Pull the latest version of an image in the background. Stopped servers switch to it on their next start; with restart_running=true running servers using it are restarted right away. Returns the affected servers.",
		func(_ context.Context, userID uint, in updateImageInput) (any, error) {
			affected, err := service.UpdateImage(in.ImageName, userID, in.RestartRunning)
			if err != nil {
				return nil, err
			}
			return map[string]any{"ok": true, "image_name": in.ImageName, "affected_servers": affected}, nil
		})

}

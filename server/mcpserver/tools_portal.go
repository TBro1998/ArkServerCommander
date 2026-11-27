package mcpserver

import (
	"context"
	"strconv"

	"ark-server-commander/models"
)

type portalSettingsInput struct {
	Enabled       bool                `json:"enabled" jsonschema:"Run the player portal container"`
	Port          int                 `json:"port,omitempty" jsonschema:"Host TCP port for the portal (1024-65535, default 8090). Must not clash with server ports"`
	PublicURL     string              `json:"public_url,omitempty" jsonschema:"URL players use when the owner has a domain or reverse proxy, e.g. https://ark.example.com. Empty: http://<public_host>:<port>"`
	Name          string              `json:"name,omitempty" jsonschema:"Site name (up to 64 characters). Empty: the first server's name"`
	Description   string              `json:"description,omitempty" jsonschema:"Site description, plain text (up to 2000 characters)"`
	QQGroup       string              `json:"qq_group,omitempty" jsonschema:"QQ group number (5-12 digits)"`
	Discord       string              `json:"discord,omitempty" jsonschema:"Discord invite URL"`
	Links         []models.PortalLink `json:"links,omitempty" jsonschema:"Up to 10 links {title, url}, http(s) only"`
	HiddenServers []uint              `json:"hidden_servers,omitempty" jsonschema:"Server IDs not shown on the portal"`
}

type announcementInput struct {
	Title  string `json:"title" jsonschema:"Title, plain text, up to 100 characters, no line breaks"`
	Body   string `json:"body,omitempty" jsonschema:"Content, plain text (line breaks are kept), up to 10000 characters"`
	Pinned bool   `json:"pinned,omitempty" jsonschema:"Show at the top"`
}

type updateAnnouncementInput struct {
	ID uint `json:"id" jsonschema:"Announcement ID from list_announcements"`
	announcementInput
}

type announcementIDInput struct {
	ID uint `json:"id" jsonschema:"Announcement ID from list_announcements"`
}

func registerPortalTools(r *registry) {
	query(r, "get_portal", "Player portal",
		"Player portal settings and status: a read-only web page for players (announcements, server status, ArkShop store) served by a separate container. Returns settings, links, status (disabled / starting / running / error, url players open), and each server's hidden flag and whether it has an ArkShop store (shop_error explains a broken ArkShop config.json).",
		func(_ context.Context, _ uint, _ noInput) (any, error) {
			return service.GetPortal()
		})

	action(r, idempotent, "update_portal", "Update player portal settings",
		"Replace all player portal settings at once (fields left out are cleared, so call get_portal first and send everything back). Enabling starts the portal container in the background (the image is pulled first if missing); poll get_portal until status.state is running. The owner must allow the port in the firewall.",
		func(_ context.Context, _ uint, in portalSettingsInput) (any, error) {
			return service.UpdatePortal(models.PortalSettingsRequest{
				Enabled: in.Enabled, Port: in.Port, PublicURL: in.PublicURL, Name: in.Name, Description: in.Description,
				QQGroup: in.QQGroup, Discord: in.Discord, Links: in.Links, HiddenServers: in.HiddenServers,
			})
		})

	query(r, "list_announcements", "List announcements",
		"Player portal announcements, newest first.",
		func(_ context.Context, _ uint, _ noInput) (any, error) {
			list, err := service.ListAnnouncements()
			return map[string]any{"announcements": list}, err
		})

	action(r, additive, "create_announcement", "Publish announcement",
		"Publish an announcement on the player portal (plain text). Open portal pages update within seconds. At most 200 announcements.",
		func(_ context.Context, _ uint, in announcementInput) (any, error) {
			return service.CreateAnnouncement(models.PortalAnnouncementRequest{Title: in.Title, Body: in.Body, Pinned: in.Pinned})
		})

	action(r, idempotent, "update_announcement", "Edit announcement",
		"Replace an announcement's title, body and pinned flag.",
		func(_ context.Context, _ uint, in updateAnnouncementInput) (any, error) {
			return service.UpdateAnnouncement(strconv.FormatUint(uint64(in.ID), 10),
				models.PortalAnnouncementRequest{Title: in.Title, Body: in.Body, Pinned: in.Pinned})
		})

	action(r, destructive, "delete_announcement", "Delete announcement",
		"Delete an announcement from the player portal.",
		func(_ context.Context, _ uint, in announcementIDInput) (any, error) {
			if err := service.DeleteAnnouncement(strconv.FormatUint(uint64(in.ID), 10)); err != nil {
				return nil, err
			}
			return map[string]any{"deleted": in.ID}, nil
		})
}

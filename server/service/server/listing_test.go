package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
)

func TestNormalizeListing(t *testing.T) {
	p, err := normalizeListing(models.ListingProfile{
		Description: "  简介\r\n第二行  ", Tags: []string{" PvE ", "pve", "新手,友好", ""}, QQGroup: " 123456789 ", Discord: "https://discord.gg/abc",
	})
	if err != nil || p.Description != "简介\n第二行" || strings.Join(p.Tags, "|") != "PvE|新手 友好" || p.QQGroup != "123456789" {
		t.Fatalf("normalize: %+v %v", p, err)
	}
	if p, _ := normalizeListing(models.ListingProfile{}); p.Tags == nil {
		t.Fatal("tags should be an empty slice, not nil (nil would be sent as null and ignored by the listing site)")
	}
	bad := []models.ListingProfile{
		{Description: strings.Repeat("字", 501)},
		{Tags: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}},
		{Tags: []string{strings.Repeat("标", 17)}},
		{QQGroup: "12ab"},
		{Discord: "https://example.com/discord"},
	}
	for _, b := range bad {
		if _, err := normalizeListing(b); !errors.Is(err, ErrInvalidListing) {
			t.Errorf("expected ErrInvalidListing for %+v, got %v", b, err)
		}
	}
}

func TestListingSentAndRepublished(t *testing.T) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "status": "running", "listing_description": "旧简介", "listing_tags": "PvE,新手"})

	ready := true
	now := time.Now()
	var bodies []publishProfile
	f := &fakeListing{handler: func(w http.ResponseWriter, r *http.Request) {
		var body publishProfile
		json.NewDecoder(r.Body).Decode(&body)
		bodies = append(bodies, body)
		okHandler(w, r)
	}}
	p := newTestPublisher(t, f, &ready, &now)
	p.profile = buildPublishProfile
	old := publisher
	publisher = p
	t.Cleanup(func() { publisher = old })

	p.tick()
	expectCalls(t, f, "/publish")
	if bodies[0].Description != "旧简介" || strings.Join(bodies[0].Tags, ",") != "PvE,新手" {
		t.Fatalf("listing not sent: %+v", bodies[0])
	}

	// 运行中修改资料：重新发布，但距上次发布至少 listingRepublishGap
	svc := &ServerService{}
	if _, _, err := svc.UpdateServer(s.UserID, fmt.Sprint(s.ID), models.ServerUpdateRequest{
		Listing: &models.ListingProfile{Description: "新简介", QQGroup: "12345"},
	}); err != nil {
		t.Fatal(err)
	}
	p.tick()
	expectCalls(t, f, "/publish")
	now = now.Add(listingRepublishGap)
	p.tick()
	expectCalls(t, f, "/publish", "/publish")
	if b := bodies[1]; b.Description != "新简介" || b.QQGroup != "12345" || b.Tags == nil || len(b.Tags) != 0 {
		t.Fatalf("republished listing: %+v", b)
	}

	// 资料没变：不重新发布
	svc.UpdateServer(s.UserID, fmt.Sprint(s.ID), models.ServerUpdateRequest{
		Listing: &models.ListingProfile{Description: "新简介", QQGroup: "12345"},
	})
	now = now.Add(time.Hour)
	p.tick()
	expectCalls(t, f, "/publish", "/publish")

	// 资料不合法：保存失败
	if _, _, err := svc.UpdateServer(s.UserID, fmt.Sprint(s.ID), models.ServerUpdateRequest{
		Listing: &models.ListingProfile{QQGroup: "abc"},
	}); !errors.Is(err, ErrInvalidListing) {
		t.Fatalf("expected ErrInvalidListing, got %v", err)
	}
}

func TestFetchPublishStats(t *testing.T) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)
	ready := true
	now := time.Now()
	f := &fakeListing{handler: func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer afp_tok" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"id": 9, "url": "https://arkfgf.com/servers/9", "badge": "https://arkfgf.com/badge/9.png",
			"week": map[string]int{"views": 12, "connects": 3, "copies": 1}})
	}}
	p := newTestPublisher(t, f, &ready, &now)
	statsCache = map[uint]cachedStats{}

	if _, err := p.fetchStats(s); !errors.Is(err, ErrNotPublished) {
		t.Fatalf("no token: %v", err)
	}
	s.PublishToken = "afp_tok"
	stats, err := p.fetchStats(s)
	if err != nil || stats.Week.Views != 12 || stats.Badge == "" {
		t.Fatalf("stats: %+v %v", stats, err)
	}
	// 缓存期内不再请求
	p.fetchStats(s)
	expectCalls(t, f, "/publish/stats")
	now = now.Add(publishStatsTTL)
	p.fetchStats(s)
	expectCalls(t, f, "/publish/stats", "/publish/stats")

	statsCache = map[uint]cachedStats{}
	s.PublishToken = "afp_revoked"
	if _, err := p.fetchStats(s); !errors.Is(err, ErrNotPublished) {
		t.Fatalf("revoked token: %v", err)
	}
}

package models

import (
	"slices"
	"strings"
	"testing"
)

func testServer() Server {
	return Server{
		Identifier:    "survival",
		SessionName:   "My Ark Server",
		ClusterID:     "cluster-a",
		Port:          7777,
		QueryPort:     27015,
		RCONPort:      32330,
		AdminPassword: "hunter2",
		Map:           "TheIsland",
		MaxPlayers:    70,
		GameModIds:    "123,456",
	}
}

// GenerateArgs 必须对相同配置生成完全相同的结果，否则启动时会误判参数变化而重建容器
func TestGenerateArgsIsDeterministic(t *testing.T) {
	server := testServer()
	args := FromServer(server)
	want := args.GenerateArgsString(server)

	for i := 0; i < 50; i++ {
		if got := args.GenerateArgsString(server); got != want {
			t.Fatalf("第 %d 次生成的结果不同:\n want: %s\n  got: %s", i+1, want, got)
		}
	}
}

func TestGenerateArgsBasics(t *testing.T) {
	server := testServer()
	got := FromServer(server).GenerateArgs(server)

	wantURL := "TheIsland?listen?Port=7777?QueryPort=27015?MaxPlayers=70?RCONEnabled=True" +
		"?RCONPort=32330?SessionName=My Ark Server?GameModIds=123,456"
	if len(got) == 0 || got[0] != wantURL {
		t.Fatalf("地图 URL 不符合预期:\n want: %s\n  got: %v", wantURL, got)
	}

	// 集群 ID 必须作为单独一个参数
	if !slices.Contains(got, "-clusterid=cluster-a") {
		t.Errorf("缺少集群参数，实际: %v", got)
	}
}

func TestGenerateArgsSkipsEmptyAndZeroValues(t *testing.T) {
	server := testServer()
	args := NewServerArgs()
	args.QueryParams["ServerPassword"] = ""
	args.QueryParams["bRawUdp"] = "false"
	args.CommandLineArgs["MinimumTimeBetweenInventoryRetrieval"] = float64(0)
	args.CommandLineArgs["MaxTamedDinos"] = 1000000.0
	args.CustomArgs = []string{"-NoTransferFromFiltering -clusterdiroverride=/ark/cluster"}

	got := args.GenerateArgs(server)

	for _, want := range []string{"-MaxTamedDinos=1000000", "-NoTransferFromFiltering", "-clusterdiroverride=/ark/cluster"} {
		if !slices.Contains(got, want) {
			t.Errorf("缺少参数 %q，实际: %v", want, got)
		}
	}

	for _, unwanted := range []string{
		"-MinimumTimeBetweenInventoryRetrieval=0",
		"-MinimumTimeBetweenInventoryRetrieval=0.000000",
		"?ServerPassword=",
		"?ServerPassword",
		"?bRawUdp=false",
	} {
		if slices.Contains(got, unwanted) {
			t.Errorf("不应包含值为空的参数: %q", unwanted)
		}
	}
}

// 管理员密码不能出现在启动参数里，否则 docker inspect 能直接看到
func TestGenerateArgsOmitsAdminPassword(t *testing.T) {
	server := testServer()
	args := FromServer(server)
	args.QueryParams["ServerAdminPassword"] = "from-query-params"

	for _, a := range args.GenerateArgs(server) {
		if strings.Contains(a, "hunter2") || strings.Contains(a, "from-query-params") || strings.Contains(a, "ServerAdminPassword") {
			t.Fatalf("启动参数不应包含管理员密码: %q", a)
		}
	}
	if strings.Contains(args.GenerateArgsString(server), "hunter2") {
		t.Fatal("启动参数字符串不应包含管理员密码")
	}
}

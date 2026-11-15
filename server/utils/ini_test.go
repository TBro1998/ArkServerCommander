package utils

import "testing"

func TestSetIniValue(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"空内容", "", "[ServerSettings]\nServerAdminPassword=pw\n"},
		{"替换已有值", "[ServerSettings]\nServerAdminPassword=old\nMaxPlayers=70\n", "[ServerSettings]\nServerAdminPassword=pw\nMaxPlayers=70\n"},
		{"键名不区分大小写并去重", "[serversettings]\nserveradminpassword=a\nX=1\nServerAdminPassword=b\n", "[serversettings]\nServerAdminPassword=pw\nX=1\n"},
		{"section 存在但没有这个键", "[ServerSettings]\nX=1\n\n[Other]\nY=2\n", "[ServerSettings]\nServerAdminPassword=pw\nX=1\n\n[Other]\nY=2\n"},
		{"没有 section", "[Other]\nY=2\n", "[Other]\nY=2\n\n[ServerSettings]\nServerAdminPassword=pw\n"},
		{"其他 section 里的同名键不动", "[Other]\nServerAdminPassword=keep\n\n[ServerSettings]\nServerAdminPassword=old\n", "[Other]\nServerAdminPassword=keep\n\n[ServerSettings]\nServerAdminPassword=pw\n"},
		{"CRLF 换行", "[ServerSettings]\r\nServerAdminPassword=old\r\n", "[ServerSettings]\nServerAdminPassword=pw\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetIniValue(tt.in, "ServerSettings", "ServerAdminPassword", "pw")
			if got != tt.want {
				t.Fatalf("期望:\n%q\n实际:\n%q", tt.want, got)
			}
			// 再设置一次结果不变，启动时据此跳过不必要的写入
			if again := SetIniValue(got, "ServerSettings", "ServerAdminPassword", "pw"); again != got {
				t.Fatalf("重复设置结果不同:\n%q\n%q", got, again)
			}
		})
	}
}

func TestKeepIniKeys(t *testing.T) {
	source := "[ServerSettings]\nSessionName=A\nServerAdminPassword=pw1\nXPMultiplier=3\n\n[SessionSettings]\nSessionName=A\nPort=7777\n\n[MessageOfTheDay]\nMessage=hi\n"
	target := "[ServerSettings]\nSessionName=B\nServerAdminPassword=pw2\nXPMultiplier=1\n\n[SessionSettings]\nsessionname=B\nPort=7779\n"
	want := "[ServerSettings]\nSessionName=B\nServerAdminPassword=pw2\nXPMultiplier=3\n\n[SessionSettings]\nsessionname=B\nPort=7779\n\n[MessageOfTheDay]\nMessage=hi\n"
	got := KeepIniKeys(source, target, "SessionName", "ServerAdminPassword", "Port")
	if got != want {
		t.Fatalf("期望:\n%q\n实际:\n%q", want, got)
	}

	// 目标里没有的键直接去掉；注释里的同名内容不动
	got = KeepIniKeys("[ServerSettings]\n;ServerAdminPassword=x\nServerAdminPassword=pw1\nX=1\n", "", "ServerAdminPassword")
	if got != "[ServerSettings]\n;ServerAdminPassword=x\nX=1\n" {
		t.Fatalf("实际: %q", got)
	}

	if got := RemoveIniKeys("[ServerSettings]\nServerAdminPassword=pw\n", "ServerAdminPassword"); got != "[ServerSettings]\n" {
		t.Fatalf("实际: %q", got)
	}
	if got := RemoveIniKeys("ServerAdminPassword=pw\n", "ServerAdminPassword"); got != "" {
		t.Fatalf("全部删除后应为空，实际: %q", got)
	}
}

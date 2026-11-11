package auth

import (
	"testing"
	"time"
)

func TestLoginLimiter(t *testing.T) {
	l := &loginLimiter{byIP: map[string]*failureRecord{}, byUser: map[string]*failureRecord{}}

	for i := 0; i < maxFailuresPerIP-1; i++ {
		l.fail("1.1.1.1", "admin")
	}
	if l.retryAfter("1.1.1.1", "admin") != 0 {
		t.Fatal("should still allow attempts below the limit")
	}
	l.fail("1.1.1.1", "admin")
	if l.retryAfter("1.1.1.1", "admin") <= 0 {
		t.Fatal("should block the IP after reaching the limit")
	}
	if l.retryAfter("2.2.2.2", "other") != 0 {
		t.Fatal("other IPs and users should not be blocked")
	}

	// 成功登录后清零
	l.success("1.1.1.1", "admin")
	if l.retryAfter("1.1.1.1", "admin") != 0 {
		t.Fatal("success should reset the counter")
	}

	// 用户名上限：换 IP 也会被拦住
	for i := 0; i < maxFailuresByUser; i++ {
		l.fail("10.0.0."+string(rune('a'+i%26)), "target")
	}
	if l.retryAfter("9.9.9.9", "target") <= 0 {
		t.Fatal("should block the username after too many failures from any IP")
	}

	// 过期后自动解除
	l.byIP["3.3.3.3"] = &failureRecord{count: maxFailuresPerIP, first: time.Now().Add(-loginWindow - time.Second)}
	if l.retryAfter("3.3.3.3", "x") != 0 {
		t.Fatal("expired records should not block")
	}
}

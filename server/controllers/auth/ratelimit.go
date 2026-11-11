package auth

import (
	"sync"
	"time"
)

// 登录失败限流：同一来源 IP 在窗口期内失败次数过多时暂时拒绝登录；
// 同一用户名另有一个更宽的上限，防止换 IP 暴力破解单个账户
const (
	loginWindow       = 15 * time.Minute
	maxFailuresPerIP  = 10
	maxFailuresByUser = 30
)

type failureRecord struct {
	count int
	first time.Time
}

type loginLimiter struct {
	mu     sync.Mutex
	byIP   map[string]*failureRecord
	byUser map[string]*failureRecord
}

var limiter = &loginLimiter{byIP: map[string]*failureRecord{}, byUser: map[string]*failureRecord{}}

func (l *loginLimiter) active(m map[string]*failureRecord, key string, now time.Time) *failureRecord {
	r, ok := m[key]
	if !ok {
		return nil
	}
	if now.Sub(r.first) > loginWindow {
		delete(m, key)
		return nil
	}
	return r
}

// retryAfter 返回还需要等待的时间，0 表示允许尝试
func (l *loginLimiter) retryAfter(ip, username string) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	var wait time.Duration
	if r := l.active(l.byIP, ip, now); r != nil && r.count >= maxFailuresPerIP {
		wait = loginWindow - now.Sub(r.first)
	}
	if r := l.active(l.byUser, username, now); r != nil && r.count >= maxFailuresByUser {
		if w := loginWindow - now.Sub(r.first); w > wait {
			wait = w
		}
	}
	return wait
}

func (l *loginLimiter) fail(ip, username string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	for _, pair := range []struct {
		m   map[string]*failureRecord
		key string
	}{{l.byIP, ip}, {l.byUser, username}} {
		if r := l.active(pair.m, pair.key, now); r != nil {
			r.count++
		} else {
			pair.m[pair.key] = &failureRecord{count: 1, first: now}
		}
	}
	// 防止记录无限增长：顺带清理过期项
	if len(l.byIP) > 10000 {
		for k := range l.byIP {
			l.active(l.byIP, k, now)
		}
	}
}

func (l *loginLimiter) success(ip, username string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.byIP, ip)
	delete(l.byUser, username)
}

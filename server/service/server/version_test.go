package server

import "testing"

func TestNewerVersion(t *testing.T) {
	cases := []struct {
		latest, current string
		want            bool
	}{
		{"0.0.3", "0.0.2", true},
		{"v0.1.0", "0.0.9", true},
		{"0.0.10", "0.0.9", true},
		{"0.0.2", "0.0.2", false},
		{"0.0.2", "0.0.3", false},
		{"0.0.2", "0.0.2-b1", true},
		{"0.0.2-b2", "0.0.2-b1", true},
		{"0.0.2-b1", "0.0.2", false},
		{"0.0.3", "dev", false},
		{"latest", "0.0.1", false},
	}
	for _, c := range cases {
		if got := newerVersion(c.latest, c.current); got != c.want {
			t.Errorf("newerVersion(%q, %q) = %v, want %v", c.latest, c.current, got, c.want)
		}
	}
}

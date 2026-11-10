package docker_manager

import "testing"

func TestSummaryStatus(t *testing.T) {
	tests := map[string]string{
		"running":    "running",
		"paused":     "running",
		"exited":     "stopped",
		"created":    "stopped",
		"restarting": "starting",
		"dead":       "unknown",
		"removing":   "unknown",
		"":           "unknown",
	}
	for state, want := range tests {
		if got := summaryStatus(state); got != want {
			t.Errorf("summaryStatus(%q) = %q，期望 %q", state, got, want)
		}
	}
}

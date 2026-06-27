package end

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sceredi/co-type/common/domain"
)

func TestViewFormatsStatisticsCorrectly(t *testing.T) {
	tests := []struct {
		name       string
		stats      domain.GameStats
		checkFuncs []func(string) bool
	}{
		{
			name: "displays_header",
			stats: domain.GameStats{
				TotalTime: 10 * time.Second,
				Lobby: domain.Lobby{
					Game: domain.GameInfo{Snippet: "hello world"},
				},
			},
			checkFuncs: []func(string) bool{
				func(v string) bool { return strings.Contains(v, "Game statistics:") },
			},
		},
		{
			name: "displays_footer",
			stats: domain.GameStats{
				TotalTime: 10 * time.Second,
				Lobby: domain.Lobby{
					Game: domain.GameInfo{Snippet: "hello world"},
				},
			},
			checkFuncs: []func(string) bool{
				func(v string) bool { return strings.Contains(v, "ctrl+c") },
			},
		},
		{
			name: "formats_duration_correctly",
			stats: domain.GameStats{
				TotalTime: 65 * time.Second, // 1 minute 5 seconds
				Lobby: domain.Lobby{
					Game: domain.GameInfo{Snippet: "test"},
				},
			},
			checkFuncs: []func(string) bool{
				func(v string) bool { return strings.Contains(v, "01:05") },
			},
		},
		{
			name: "formats_duration_zero_padded",
			stats: domain.GameStats{
				TotalTime: 5 * time.Second, // 0 minutes 5 seconds
				Lobby: domain.Lobby{
					Game: domain.GameInfo{Snippet: "test"},
				},
			},
			checkFuncs: []func(string) bool{
				func(v string) bool { return strings.Contains(v, "00:05") },
			},
		},
		{
			name: "formats_large_duration",
			stats: domain.GameStats{
				TotalTime: 125 * time.Second, // 2 minutes 5 seconds
				Lobby: domain.Lobby{
					Game: domain.GameInfo{Snippet: "test"},
				},
			},
			checkFuncs: []func(string) bool{
				func(v string) bool { return strings.Contains(v, "02:05") },
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.stats)
			view := m.View()

			for _, checkFunc := range tt.checkFuncs {
				if !checkFunc(view) {
					t.Fatalf("view does not satisfy check: %s", view)
				}
			}
		})
	}
}

func TestCPSCalculation(t *testing.T) {
	tests := []struct {
		name      string
		snippet   string
		totalTime time.Duration
		wantCPS   string
	}{
		{
			name:      "calculates_cps_correctly",
			snippet:   "hello world",
			totalTime: 10 * time.Second,
			wantCPS:   "1.10", // 11 chars / 10 seconds
		},
		{
			name:      "calculates_cps_for_one_second",
			snippet:   "typing",
			totalTime: 1 * time.Second,
			wantCPS:   "6.00", // 6 chars / 1 second
		},
		{
			name:      "calculates_cps_for_multiple_minutes",
			snippet:   "the quick brown fox",
			totalTime: 2 * time.Minute, // 120 seconds
			wantCPS:   "0.16",          // 19 chars / 120 seconds ≈ 0.158...
		},
		{
			name:      "handles_zero_duration",
			snippet:   "test",
			totalTime: 0,
			wantCPS:   "0.00", // Should show 0.00 to avoid divide by zero
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := domain.GameStats{
				TotalTime: tt.totalTime,
				Lobby: domain.Lobby{
					Game: domain.GameInfo{Snippet: tt.snippet},
				},
			}
			m := New(stats)
			view := m.View()

			expectedCPS := fmt.Sprintf("Characters per second -> %s", tt.wantCPS)
			if !strings.Contains(view, expectedCPS) {
				t.Fatalf("expected view to contain %q, got: %s", expectedCPS, view)
			}
		})
	}
}

func TestViewRendersAllComponents(t *testing.T) {
	stats := domain.GameStats{
		TotalTime: 30 * time.Second,
		Lobby: domain.Lobby{
			Game: domain.GameInfo{Snippet: "hello"},
		},
	}

	m := New(stats)
	view := m.View()

	requiredStrings := []string{
		"Game statistics:",
		"Duration ->",
		"mm:ss",
		"Characters per second ->",
		"ctrl+c",
	}

	for _, required := range requiredStrings {
		if !strings.Contains(view, required) {
			t.Fatalf("expected view to contain %q, got: %s", required, view)
		}
	}
}

func TestInit(t *testing.T) {
	stats := domain.GameStats{
		TotalTime: 10 * time.Second,
		Lobby: domain.Lobby{
			Game: domain.GameInfo{Snippet: "test"},
		},
	}
	m := New(stats)
	cmd := m.Init()

	if cmd != nil {
		t.Fatalf("expected Init to return nil command, got %v", cmd)
	}
}

func TestUpdate(t *testing.T) {
	stats := domain.GameStats{
		TotalTime: 10 * time.Second,
		Lobby: domain.Lobby{
			Game: domain.GameInfo{Snippet: "test"},
		},
	}
	m := New(stats)

	// Update should be a no-op currently
	m2, _ := m.Update(nil)

	if m2.stats.TotalTime != m.stats.TotalTime {
		t.Fatalf("expected Update to not change stats")
	}
}

package playerinfo

import (
	"strings"
	"testing"

	"github.com/sceredi/co-type/common/domain"
)

func TestRenderWithConfiguredPlayer(t *testing.T) {
	player := domain.Player{
		AllowedCharacters: "[a-z]",
		BlockedCharacters: "[0-9]",
		CanDelete:         true,
	}

	rendered := Render(player)

	for _, want := range []string{"Allowed characters: \"[a-z]\"", "Blocked characters: \"[0-9]\"", "Backspace:", "enabled"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("expected rendered player info to include %q, got: %q", want, rendered)
		}
	}
}

func TestRenderWithUnsetPlayerSettings(t *testing.T) {
	player := domain.Player{CanDelete: false}

	rendered := Render(player)

	for _, want := range []string{"Allowed characters not yet set", "Blocked characters not yet set", "Backspace:", "disabled"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("expected rendered player info to include %q, got: %q", want, rendered)
		}
	}
}

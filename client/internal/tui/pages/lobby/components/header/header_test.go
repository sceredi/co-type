package header

import (
	"strings"
	"testing"

	"github.com/sceredi/co-type/common/domain"
)

func TestRenderIncludesLobbyMetadata(t *testing.T) {
	player := &domain.Player{Name: "TestUser"}
	lobby := &domain.Lobby{
		ID: "ABCD",
		Players: []*domain.Player{
			{Name: "Host", IsReady: true},
			{Name: "Guest", IsReady: false},
			{Name: "TestUser", IsReady: true},
		},
	}

	rendered := Render(player, lobby)

	for _, want := range []string{"co-type", "ABCD", "TestUser", "2/3"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("expected rendered header to include %q, got: %q", want, rendered)
		}
	}
}

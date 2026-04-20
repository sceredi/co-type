package playerlist

import (
	"strings"
	"testing"

	"github.com/sceredi/co-type/common/domain"
)

func testLobby() *domain.Lobby {
	host := &domain.Player{Name: "Ana", IsReady: true}
	guest := &domain.Player{Name: "Bo", IsReady: false}
	third := &domain.Player{Name: "Ci", IsReady: true}

	return &domain.Lobby{
		Host:    host,
		Players: []*domain.Player{host, guest, third},
	}
}

func TestRenderIncludesAllPlayersAndSingleSelectionPointer(t *testing.T) {
	rendered := Render(1, testLobby())

	for _, want := range []string{"Ana", "Bo", "Ci"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("expected rendered player list to include player %q, got: %q", want, rendered)
		}
	}

	if strings.Count(rendered, "→") != 1 {
		t.Fatalf("expected exactly one selection pointer, got: %q", rendered)
	}
}

func TestRenderIncludesReadyMarkerForReadyPlayers(t *testing.T) {
	rendered := Render(0, testLobby())

	if strings.Count(rendered, "r") < 2 {
		t.Fatalf("expected ready marker for two ready players, got: %q", rendered)
	}
}

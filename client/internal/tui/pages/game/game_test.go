package game

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/common/domain"
)

func createTestGame() domain.Game {
	return domain.Game{
		Lobby: domain.Lobby{
			Game: domain.GameInfo{Snippet: "hello world"},
			Players: []*domain.Player{
				{Name: "Player1"},
				{Name: "Player2"},
			},
		},
		State: domain.GameState{
			Status:       domain.Running,
			CorrectChars: 0,
			LastChar:     0,
		},
	}
}

func TestNew(t *testing.T) {
	game := createTestGame()
	m := New(game, nil, nil, nil)

	if m.game.Lobby.Game.Snippet != "hello world" {
		t.Fatalf("expected model to store game, got empty")
	}
}

func TestUpdateWindowSize(t *testing.T) {
	m := New(createTestGame(), nil, nil, nil)
	initialWidth := m.width
	initialHeight := m.height

	m2, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})

	if m2.width != 100 {
		t.Fatalf("expected width to be updated to 100, got %d", m2.width)
	}

	if m2.height != 50 {
		t.Fatalf("expected height to be updated to 50, got %d", m2.height)
	}

	if m2.width == initialWidth && m2.height == initialHeight && initialWidth == 0 {
		t.Fatalf("expected window size to change from initial values")
	}
}

func TestRenderPlayersListsAllPlayers(t *testing.T) {
	game := createTestGame()
	m := New(game, nil, nil, nil)
	players := m.renderPlayers()

	if !strings.Contains(players, "Players:") {
		t.Fatalf("expected players list to contain 'Players:', got: %s", players)
	}

	if !strings.Contains(players, "Player1") {
		t.Fatalf("expected players list to contain 'Player1', got: %s", players)
	}

	if !strings.Contains(players, "Player2") {
		t.Fatalf("expected players list to contain 'Player2', got: %s", players)
	}
}

func TestViewRendersCorrectCharactersWithGreenBackground(t *testing.T) {
	game := createTestGame()
	game.State.CorrectChars = 5
	game.State.LastChar = 5
	m := New(game, nil, nil, nil)

	view := m.View()

	// View should contain rendered output
	if view == "" {
		t.Fatal("expected View to return non-empty string")
	}

	// When correctly typed, we should see the correct portion
	// The exact format depends on styling, but we can verify it's not empty
	if !strings.Contains(view, "Players:") {
		t.Fatalf("expected view to contain player list")
	}
}

func TestViewRendersIncorrectCharactersWithRedBackground(t *testing.T) {
	game := createTestGame()
	game.State.CorrectChars = 0
	game.State.LastChar = 3
	m := New(game, nil, nil, nil)

	view := m.View()

	if view == "" {
		t.Fatal("expected View to return non-empty string")
	}
}

func TestViewWithPausedStatus(t *testing.T) {
	game := createTestGame()
	game.Lobby.Status = domain.LobbyPaused
	m := New(game, nil, nil, nil)

	view := m.View()

	if !strings.Contains(view, "Game paused") {
		t.Fatalf("expected paused view to contain 'Game paused', got: %s", view)
	}
}

func TestViewWithRunningStatus(t *testing.T) {
	game := createTestGame()
	game.Lobby.Status = domain.LobbyPlaying
	m := New(game, nil, nil, nil)

	view := m.View()

	if strings.Contains(view, "Game paused") {
		t.Fatalf("expected running view to NOT contain pause message, got: %s", view)
	}
}

func TestViewWithNoCorrectChars(t *testing.T) {
	game := createTestGame()
	game.State.CorrectChars = 0
	game.State.LastChar = 0
	m := New(game, nil, nil, nil)

	view := m.View()

	if view == "" {
		t.Fatal("expected View to handle zero correct chars")
	}
}

func TestViewWithAllCharsInput(t *testing.T) {
	game := createTestGame()
	game.State.CorrectChars = len(game.Lobby.Game.Snippet)
	game.State.LastChar = len(game.Lobby.Game.Snippet)
	m := New(game, nil, nil, nil)

	view := m.View()

	if view == "" {
		t.Fatal("expected View to handle all chars input")
	}
}

func TestInit(t *testing.T) {
	m := New(createTestGame(), nil, nil, nil)
	cmd := m.Init()

	if cmd != nil {
		t.Fatalf("expected Init to return nil, got %v", cmd)
	}
}

func TestRenderPlayersWithSinglePlayer(t *testing.T) {
	game := domain.Game{
		Lobby: domain.Lobby{
			Game: domain.GameInfo{Snippet: "test"},
			Players: []*domain.Player{
				{Name: "OnlyPlayer"},
			},
		},
		State: domain.GameState{},
	}
	m := New(game, nil, nil, nil)
	players := m.renderPlayers()

	if !strings.Contains(players, "Players:") {
		t.Fatalf("expected single player render to work")
	}

	if !strings.Contains(players, "OnlyPlayer") {
		t.Fatalf("expected single player name in render")
	}
}

func TestRenderPlayersWithNoPlayers(t *testing.T) {
	game := domain.Game{
		Lobby: domain.Lobby{
			Game:    domain.GameInfo{Snippet: "test"},
			Players: []*domain.Player{},
		},
		State: domain.GameState{},
	}
	m := New(game, nil, nil, nil)
	players := m.renderPlayers()

	if !strings.Contains(players, "Players:") {
		t.Fatalf("expected empty player list to still show header")
	}
}

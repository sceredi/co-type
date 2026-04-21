package welcome

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	welcome_messages "github.com/sceredi/co-type/client/internal/tui/pages/welcome/messages"
)

func TestTabCyclesThroughFocusSlots(t *testing.T) {
	m := New()

	if m.focus != focusName {
		t.Fatalf("expected initial focusName, got %d", m.focus)
	}
	if !m.nameTi.Focused() {
		t.Fatal("expected name input to start focused")
	}

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusCode {
		t.Fatalf("expected focusCode after first tab, got %d", m.focus)
	}
	if !m.codeTi.Focused() {
		t.Fatal("expected code input to be focused after first tab")
	}

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusCreate {
		t.Fatalf("expected focusCreate after second tab, got %d", m.focus)
	}

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusJoin {
		t.Fatalf("expected focusJoin after third tab, got %d", m.focus)
	}

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusName {
		t.Fatalf("expected focusName after wrapping, got %d", m.focus)
	}
	if !m.nameTi.Focused() {
		t.Fatal("expected name input to be focused after wrapping")
	}
}

func TestEnterOnNameAdvancesToCode(t *testing.T) {
	m := New()

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if m.focus != focusCode {
		t.Fatalf("expected focusCode after enter on name, got %d", m.focus)
	}
}

func TestJoinLobbyErrorMsgUpdatesViewError(t *testing.T) {
	m := New()
	m, _ = m.Update(welcome_messages.JoinLobbyErrorMsg{Error: "Unable to join lobby \"ABCD\""})

	if m.error != "Unable to join lobby \"ABCD\"" {
		t.Fatalf("expected error message to be stored, got %q", m.error)
	}
	if !strings.Contains(m.View(), m.error) {
		t.Fatalf("expected rendered view to include error %q", m.error)
	}
}

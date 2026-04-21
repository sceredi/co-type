package settings

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/sceredi/co-type/common/domain"
)

func newTestPlayer() *domain.Player {
	return &domain.Player{
		Name:              "TestUser",
		AllowedCharacters: "[a-z]",
		BlockedCharacters: "[0-9]",
		CanDelete:         true,
	}
}

func TestTabCyclesFocus(t *testing.T) {
	m := New(newTestPlayer())

	// Tab once: focusAllowed -> focusBlocked
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if !m.blockedList.Focused() {
		t.Error("expected blockedList to be focused after first tab")
	}

	// Tab again: focusBlocked -> focusBackspace
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusBackspace {
		t.Errorf("expected focusBackspace, got %d", m.focus)
	}
}

func TestSpaceTogglesBackspace(t *testing.T) {
	m := New(newTestPlayer()) // CanDelete starts true

	// Tab twice to reach focusBackspace
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})

	// Space should toggle off
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if m.backspaceAllowed {
		t.Error("expected backspaceAllowed=false after space")
	}

	// Space again should toggle back on
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if !m.backspaceAllowed {
		t.Error("expected backspaceAllowed=true after second space")
	}
}

func TestConfirmEmitsCommand(t *testing.T) {
	m := New(newTestPlayer())

	// 3 tabs: allowed -> blocked -> backspace -> confirm
	for range 3 {
		m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	}
	if m.focus != focusConfirm {
		t.Fatalf("expected focusConfim, got %d", m.focus)
	}

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a batch command from confirm, got nil")
	}
}

func TestCancelEmitsCommand(t *testing.T) {
	m := New(newTestPlayer())

	// 4 tabs to reach focusCancel
	for range 4 {
		m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	}
	if m.focus != focusCancel {
		t.Fatalf("expected focusCancel, got %d", m.focus)
	}

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command from cancel, got nil")
	}
}

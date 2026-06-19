package service

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/common/domain"
)

type mockLobbyRepository struct {
	lobbyID    repository.LobbyID
	serverName repository.ServerName
	err        error
}

func (m *mockLobbyRepository) Create(lobbyID repository.LobbyID, serverName repository.ServerName) error {
	m.lobbyID = lobbyID
	m.serverName = serverName
	return m.err
}

func (m *mockLobbyRepository) Get(lobbyID repository.LobbyID) (repository.ServerName, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.serverName, nil
}

type mockServerService struct {
	server *domain.Server
	err    error
}

func (m *mockServerService) Create(name, addr string, port int) (*domain.Server, error) {
	return nil, nil
}

func (m *mockServerService) LowestLoad() (*domain.Server, error) {
	return nil, nil
}

func (m *mockServerService) GetByName(name string) (*domain.Server, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.server, nil
}

func TestLobbyService_Create(t *testing.T) {
	repo := &mockLobbyRepository{}
	svc := NewLobbyService(repo, &mockServerService{})

	err := svc.Create("lobby1", "server1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.lobbyID != "lobby1" || repo.serverName != "server1" {
		t.Fatalf("expected lobby1 and server1, got %s and %s", repo.lobbyID, repo.serverName)
	}
}

func TestLobbyService_Get(t *testing.T) {
	expectedServer := &domain.Server{Name: "server1"}
	repo := &mockLobbyRepository{serverName: "server1"}
	serverSvc := &mockServerService{server: expectedServer}
	svc := NewLobbyService(repo, serverSvc)

	server, err := svc.Get("lobby1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if server != expectedServer {
		t.Fatalf("expected %+v, got %+v", expectedServer, server)
	}
}

func TestLobbyService_Get_Error(t *testing.T) {
	repoErr := errors.New("not found")
	repo := &mockLobbyRepository{err: repoErr}
	svc := NewLobbyService(repo, &mockServerService{})

	_, err := svc.Get("lobby1")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}

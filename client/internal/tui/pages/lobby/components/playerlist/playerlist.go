package playerlist

import "github.com/sceredi/co-type/common/domain"

type Model struct {
	lobby domain.Lobby
}

func New(lobby domain.Lobby) Model {
	return Model{lobby: lobby}
}

package lobby

func handleReadyCmd(m *Model) {
	m.player.IsReady = !m.player.IsReady
}

func handlePlayerSelect(selectedPlayer int, m *Model) {
	if selectedPlayer < 0 {
		m.selectedPlayer = len(m.lobby.Players) - 1
	} else {
		m.selectedPlayer = selectedPlayer % len(m.lobby.Players)
	}
}

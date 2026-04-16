package lobby

import "log"

func handleReadyCmd(m *Model) {
	m.player.IsReady = !m.player.IsReady
	log.Printf("Setting ready to %t\n", m.player.IsReady)
}

func handlePlayerSelect(selectedPlayer int, m *Model) {
	if selectedPlayer < 0 {
		m.selectedPlayer = len(m.lobby.Players) - 1
	} else {
		m.selectedPlayer = selectedPlayer % len(m.lobby.Players)
	}
}

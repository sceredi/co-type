package lobby

import "log"

// --- Message handlers ---
func handleReadyCmd(m *Model) {
	m.player.IsReady = !m.player.IsReady
	log.Printf("Setting ready to %t\n", m.player.IsReady)
}

package main

import (
	"sync"
)

// Player is used to hold player data.
type Player struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Position Vector `json:"position"`
}

// PlayerManager tracks the player data.
type PlayerManager struct {
	players map[int64]*Player
	mutex   sync.Mutex
}

// CreatePlayerManager will create and return a Manager instance
func CreatePlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[int64]*Player),
	}
}

// Register will add the player to the players map
func (m *PlayerManager) Register(player *Player) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.players[player.ID] = player
}

// Unregister will remove the player from the player map
func (m *PlayerManager) Unregister(playerID int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.players, playerID)
}

// HasPlayer will return true if the player is in the players map.
func (m *PlayerManager) HasPlayer(playerID int64) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.players[playerID]
	return ok
}

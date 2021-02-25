package game

import (
	"context"
	"encoding/binary"
	"log"
	"math"
	"sync"

	"github.com/cdrpl/idlemon/pkg/ws"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Player is used to hold player data.
type Player struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Position    Vector `json:"position"`
	Destination Vector `json:"-"` // The position that the player is moving towards.
}

// Pos will return the player position as a byte slice.
// The player ID is included at the start of the slice.
func (p *Player) Pos() []byte {
	buf := make([]byte, 12)
	binary.LittleEndian.PutUint32(buf[:4], p.ID)
	binary.LittleEndian.PutUint32(buf[4:8], math.Float32bits(float32(p.Position.X)))
	binary.LittleEndian.PutUint32(buf[8:], math.Float32bits(float32(p.Position.Y)))
	return buf
}

// PlayerManager tracks the player data.
type PlayerManager struct {
	players map[uint32]*Player
	mutex   sync.Mutex
}

// CreatePlayerManager will create and return a Manager instance
func CreatePlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[uint32]*Player),
	}
}

// Register will add the player to the players map
func (m *PlayerManager) Register(player *Player) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.players[player.ID] = player
}

// Unregister will remove the player from the player map
func (m *PlayerManager) Unregister(playerID uint32) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.players, playerID)
}

// HasPlayer will return true if the player is in the players map.
func (m *PlayerManager) HasPlayer(playerID uint32) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.players[playerID]
	return ok
}

// GetPlayerCopy will return a copy of the player, ok will be false if player can't be found.
func (m *PlayerManager) GetPlayerCopy(playerID uint32) (Player, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if player, ok := m.players[playerID]; ok {
		return *player, ok
	}

	return Player{}, false
}

// ClearPlayerDestination will set the player destination to the current position.
func (m *PlayerManager) ClearPlayerDestination(playerID uint32) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if player, ok := m.players[playerID]; ok {
		player.Destination = player.Position
	}
}

// UpdatePlayerPositions will move the players towards their destinations by 1 tick.
func (m *PlayerManager) UpdatePlayerPositions(server *ws.Server, moveSpeed float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, player := range m.players {
		isOnline := server.PlayerOnline(player.ID)

		if isOnline && player.Position != player.Destination {
			player.Position.MoveTowards(player.Destination, moveSpeed)

			// Broadcast the updated position
			message := &ws.Message{Channel: ws.Position, Data: player.Pos()}
			go server.BroadcastAll(message)
		}
	}
}

// SendPlayerPositions will send all online player's positions to the specified player.
func (m *PlayerManager) SendPlayerPositions(server *ws.Server, playerID uint32) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !server.PlayerOnline(playerID) {
		return
	}

	for _, player := range m.players {
		isOnline := server.PlayerOnline(player.ID)

		if isOnline {
			// Broadcast the player position
			message := &ws.Message{Channel: ws.Position, Data: player.Pos()}
			go server.BroadcastSingle(message, playerID)
		}
	}
}

// SetPlayerDestination will set the player destination if player exists.
func (m *PlayerManager) SetPlayerDestination(playerID uint32, destination Vector) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if player, ok := m.players[playerID]; ok {
		player.Destination = destination
	}
}

// PlayersToSlice creates a slice using the players map, the players are copies.
func (m *PlayerManager) PlayersToSlice() []Player {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	players := make([]Player, 0)
	for _, player := range m.players {
		players = append(players, *player)
	}
	return players
}

// SavePlayerData saves the data for each player to the database. Make sure to run in a goroutine.
func (m *PlayerManager) SavePlayerData(server *ws.Server, pgPool *pgxpool.Pool) {
	players := m.PlayersToSlice()

	if len(players) == 0 {
		return
	}

	if len(players) == 1 {
		log.Printf("Saving data for %v player\n", len(players))
	} else {
		log.Printf("Saving data for %v players\n", len(players))
	}

	for _, player := range players {
		// Save player position
		_, err := pgPool.Exec(context.Background(), "UPDATE positions SET x = $1, y = $2 WHERE id = $3", player.Position.X, player.Position.Y, player.ID)
		if err != nil {
			log.Printf("Failed to save data for player %d: %v\n", player.ID, err)
			continue
		}

		// delete player data if player is offline
		if !server.PlayerOnline(player.ID) {
			m.Unregister(player.ID)
		}
	}
}

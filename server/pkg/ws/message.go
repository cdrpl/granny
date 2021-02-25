package ws

import (
	"encoding/binary"
	"errors"
)

const (
	// Chat channel for receiving chat messages.
	Chat = Channel(iota)
	// Position channel for sending player positions.
	Position
	// PlayerData channel for sending player data on login.
	PlayerData
	// Destination channel for receiving player destinations.
	Destination
	// PlayerConnected channel for sending player connected events.
	PlayerConnected
	// PlayerDisconnected channel for sending player disconnected events.
	PlayerDisconnected
	// Used to check if a number is a valid channel. (n < end)
	end
)

// Channel is used to determine how a WebSocket message will be handled.
type Channel uint16

// Message is used to easily parse data received and sent over the WebSocket connections.
type Message struct {
	PlayerID uint32 // The player ID of the message sender.
	Channel  Channel
	Data     []byte
}

// Parse will extract the channel from the data.
// The channel is represented by the first 2 bytes of data.
// An error will be returned if the data is less than 2 bytes.
func (m *Message) parse() error {
	if len(m.Data) < 2 {
		return errors.New("Message data must contain at least 2 bytes")
	}

	channelB := m.Data[0:2] // channel is the first 2 bytes
	m.Channel = Channel(binary.LittleEndian.Uint16(channelB))
	m.Data = m.Data[2:] // remove channel bytes from the data
	return nil
}

// Merge channel and data into a single byte slice, channel is placed at index 0.
// Used for broadcasting the message to clients as a single byte slice.
func (m *Message) bytes() []byte {
	// bytes will hold the entire message and be returned
	bytes := make([]byte, 0)

	// append channel to bytes slice
	channelB := make([]byte, 2)
	binary.LittleEndian.PutUint16(channelB, uint16(m.Channel))
	bytes = append(bytes, channelB...)

	// append data
	bytes = append(bytes, m.Data...)
	return bytes
}

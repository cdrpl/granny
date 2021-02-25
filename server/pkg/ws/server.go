package ws

import (
	"encoding/binary"
	"log"
	"sync"
)

// Server tracks the clients
type Server struct {
	clients  map[uint32]*Client
	Incoming chan *Message // Incoming data is sent to this channel.
	mutex    sync.Mutex    // Mutex used for the clients map.
}

// CreateServer will return a Server instance.
func CreateServer() *Server {
	return &Server{
		clients:  make(map[uint32]*Client),
		Incoming: make(chan *Message),
		mutex:    sync.Mutex{},
	}
}

// Register adds the client to the server's client map.
func (s *Server) Register(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Printf("Player '%v' connected\n", client.id)
	s.clients[client.id] = client
}

// Unregister removes the client from the server's client map.
func (s *Server) Unregister(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.clients[client.id]; ok {
		log.Printf("Player '%v' disconnected\n", client.id)
		delete(s.clients, client.id)
		close(client.send)

		// send player disconnected event
		var idB [4]byte
		binary.LittleEndian.PutUint32(idB[:], client.id)
		message := &Message{Channel: PlayerDisconnected, Data: idB[:]}
		go s.BroadcastAll(message)
	}
}

// broadcast the message to the given client. Do not call this function without locking the server mutex.
func (s *Server) broadcast(message *Message, client *Client) {
	select {
	case client.send <- message.bytes():

	default: // assume the client is dead if the send channel is full
		close(client.send)
		delete(s.clients, client.id)
	}
}

// Broadcast will send the message to the given targets.
func (s *Server) Broadcast(message *Message, targets []uint32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, id := range targets {
		if client, ok := s.clients[id]; ok {
			s.broadcast(message, client)
		}
	}
}

// BroadcastSingle will send the message to the specified client.
func (s *Server) BroadcastSingle(message *Message, target uint32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, ok := s.clients[target]; ok {
		s.broadcast(message, client)
	}
}

// BroadcastAll will send the message to all clients.
func (s *Server) BroadcastAll(message *Message) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, client := range s.clients {
		s.broadcast(message, client)
	}
}

// BroadcastAllExclude will broadcast the message to every client except the specified one.
func (s *Server) BroadcastAllExclude(message *Message, exclude uint32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, client := range s.clients {
		if id == exclude {
			continue
		}
		s.broadcast(message, client)
	}
}

// PlayerOnline will return true if the player has an active connection.
// Can be safely called from other goroutines.
func (s *Server) PlayerOnline(id uint32) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, isOnline := s.clients[id]
	return isOnline
}

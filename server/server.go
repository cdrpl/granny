package main

import (
	"log"
	"sync"
)

// Server tracks the clients
type Server struct {
	clients  map[int64]*Client
	Incoming chan []byte // Incoming data is sent to this channel.
	mutex    sync.Mutex  // Mutex used for the clients map.
}

// CreateServer will return a Server instance.
func CreateServer() *Server {
	return &Server{
		clients:  make(map[int64]*Client),
		Incoming: make(chan []byte),
		mutex:    sync.Mutex{},
	}
}

// Register adds the client to the server's client map.
func (s *Server) Register(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	log.Printf("User '%v' connected\n", client.id)
	s.clients[client.id] = client
}

// Unregister removes the client from the server's client map.
func (s *Server) Unregister(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.clients[client.id]; ok {
		log.Printf("User '%v' disconnected\n", client.id)
		delete(s.clients, client.id)
		close(client.send)
	}
}

// broadcast the message to the given client. Do not call this function without locking the server mutex.
func (s *Server) broadcast(data []byte, client *Client) {
	select {
	case client.send <- data:

	default: // assume the client is dead if the send channel is full
		close(client.send)
		delete(s.clients, client.id)
	}
}

// Broadcast will send the message to the given targets.
func (s *Server) Broadcast(data []byte, targets []int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, id := range targets {
		if client, ok := s.clients[id]; ok {
			s.broadcast(data, client)
		}
	}
}

// BroadcastSingle will send the message to the specified client.
func (s *Server) BroadcastSingle(data []byte, target int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, ok := s.clients[target]; ok {
		s.broadcast(data, client)
	}
}

// BroadcastAll will send the message to all clients.
func (s *Server) BroadcastAll(data []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, client := range s.clients {
		s.broadcast(data, client)
	}
}

// BroadcastAllExclude will broadcast the message to every client except the specified one.
func (s *Server) BroadcastAllExclude(data []byte, exclude int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, client := range s.clients {
		if id == exclude {
			continue
		}
		s.broadcast(data, client)
	}
}

// IsUserOnline will return true if the user has an active connection.
// Can be safely called from other goroutines.
func (s *Server) IsUserOnline(id int64) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, isOnline := s.clients[id]
	return isOnline
}

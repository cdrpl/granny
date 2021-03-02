package main

import (
	"fmt"
	"log"
	"sync"
)

const roomSize = 5 // Max number of users allowed in a room

// User account data.
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Room represents a game room.
type Room struct {
	Users []User `json:"users"`
}

func (r *Room) addUser(user User) {
	r.Users = append(r.Users, user)
}

func (r *Room) isFull() bool {
	return len(r.Users) >= roomSize
}

func (r *Room) hasUser(id int64) bool {
	for _, user := range r.Users {
		if user.ID == id {
			return true
		}
	}
	return false
}

// Server tracks the clients
type Server struct {
	clients    map[int64]*Client
	clientsMut sync.Mutex
	users      map[int64]User
	usersMut   sync.Mutex
	room       Room
	roomMut    sync.Mutex
	Incoming   chan Message // Incoming data is sent to this channel.
}

// CreateServer will return a Server instance.
func CreateServer() *Server {
	return &Server{
		clients:  make(map[int64]*Client),
		users:    make(map[int64]User),
		room:     Room{Users: make([]User, 0)},
		Incoming: make(chan Message),
	}
}

// Register adds the client to the server's client map.
func (s *Server) Register(client *Client, user User) {
	s.clientsMut.Lock()
	s.usersMut.Lock()

	s.clients[client.id] = client
	s.users[user.ID] = user

	s.clientsMut.Unlock()
	s.usersMut.Unlock()
}

// Unregister removes the client from the server's client map.
func (s *Server) Unregister(client *Client) {
	s.clientsMut.Lock()
	defer s.clientsMut.Unlock()

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
	s.clientsMut.Lock()
	defer s.clientsMut.Unlock()

	for _, id := range targets {
		if client, ok := s.clients[id]; ok {
			s.broadcast(data, client)
		}
	}
}

// BroadcastSingle will send the message to the specified client.
func (s *Server) BroadcastSingle(data []byte, target int64) {
	s.clientsMut.Lock()
	defer s.clientsMut.Unlock()

	if client, ok := s.clients[target]; ok {
		s.broadcast(data, client)
	}
}

// BroadcastAll will send the message to all clients.
func (s *Server) BroadcastAll(data []byte) {
	s.clientsMut.Lock()
	defer s.clientsMut.Unlock()

	for _, client := range s.clients {
		s.broadcast(data, client)
	}
}

// BroadcastAllExclude will broadcast the message to every client except the specified one.
func (s *Server) BroadcastAllExclude(data []byte, exclude int64) {
	s.clientsMut.Lock()
	defer s.clientsMut.Unlock()

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
	s.clientsMut.Lock()
	defer s.clientsMut.Unlock()

	_, isOnline := s.clients[id]
	return isOnline
}

func (s *Server) run() {
	for {
		message := <-s.Incoming

		switch message.channel {
		case JoinRoom:
			s.joinRoomHandler(message.client.id)
			break

		default:
			fmt.Printf("Recv message invalid channel %d\n:", message.channel)
		}
	}
}

// True will be returned if the user was added to the room.
func (s *Server) joinRoomHandler(userID int64) {
	s.usersMut.Lock()
	s.roomMut.Lock()
	defer s.usersMut.Unlock()
	defer s.roomMut.Unlock()

	roomIsFull := s.room.isFull()
	if roomIsFull {
		return
	}

	if user, ok := s.users[userID]; ok {
		if !s.room.hasUser(user.ID) {
			s.room.addUser(user)
		}
	}
}

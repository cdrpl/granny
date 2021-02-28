package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

const roomSize = 5 // Max number of users allowed in a room

const (
	joinRoom = iota
)

// Message represents a socket message and links the client(sender) to the data
type Message struct {
	client *Client
	data   []byte
}

func (m *Message) channel() byte {
	return m.data[0]
}

// User account data.
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Room represents a game room.
type Room struct {
	Users []User `json:"users"`
}

// If room is full, user will not be added and false will be returned
func (r *Room) addUser(user User) bool {
	isFull := r.isFull()

	if isFull {
		return false
	}

	r.Users = append(r.Users, user)
	return true
}

func (r Room) isFull() bool {
	return len(r.Users) >= roomSize
}

// Server tracks the clients
type Server struct {
	clients    map[int64]*Client
	clientsMut sync.Mutex
	users      map[int64]*User
	usersMut   sync.Mutex
	room       Room
	roomMut    sync.Mutex
	Incoming   chan Message // Incoming data is sent to this channel.
}

// CreateServer will return a Server instance.
func CreateServer() *Server {
	return &Server{
		clients:  make(map[int64]*Client),
		users:    make(map[int64]*User),
		room:     Room{Users: make([]User, 0)},
		Incoming: make(chan Message),
	}
}

// Register adds the client to the server's client map.
func (s *Server) Register(client *Client, user *User) {
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

func (s *Server) getRoom() Room {
	s.roomMut.Lock()
	defer s.roomMut.Unlock()

	return s.room
}

func (s *Server) getUser(id int64) (User, bool) {
	s.usersMut.Lock()
	defer s.usersMut.Unlock()

	if user, ok := s.users[id]; ok {
		return *user, true
	}

	return User{}, false
}

// True will be returned if the user was added to the room.
func (s *Server) addUserToRoom(userID int64) (bool, error) {
	s.roomMut.Lock()
	defer s.roomMut.Unlock()

	if user, ok := s.getUser(userID); ok {
		if ok2 := s.room.addUser(user); ok2 {
			return true, nil
		}
		return false, nil
	}

	return false, errors.New("addUserToRoom() invalid userID received")
}

func (s *Server) run() {
	for {
		message := <-s.Incoming

		switch message.channel() {
		case joinRoom:
			isJoined, err := s.addUserToRoom(message.client.id)
			if err != nil {
				log.Println(err)
				message.client.send <- []byte{2}
			}

			if isJoined {
				log.Println("user", message.client.id, "joined the room")
				message.client.send <- []byte{0}
			} else {
				log.Println("user", message.client.id, "could not join room, already full")
				message.client.send <- []byte{1}
			}
			break

		default:
			fmt.Printf("Recv message invalid channel %d\n:", message.channel())
		}
	}
}

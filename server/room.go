package main

import (
	"errors"
	"sync"
)

// Room represents a game room.
type Room struct {
	users map[int]*RoomUser
	mut   sync.Mutex
}

func newRoom() Room {
	return Room{
		users: make(map[int]*RoomUser),
	}
}

// Check if room is full, mut must be locked first.
func (r *Room) isFull() bool {
	return len(r.users) >= roomSize
}

func (r *Room) joinRoom(user *RoomUser) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	roomIsFull := r.isFull()
	if roomIsFull {
		return errors.New("Room is full")
	}

	if _, ok := r.users[user.id]; ok {
		return errors.New("User is already in the room")
	}

	r.users[user.id] = user

	// Broadcast user joined to every other user
	for _, ru := range r.users {
		ru.joined <- user
	}

	return nil
}

func (r *Room) getUser(id int) *RoomUser {
	r.mut.Lock()
	defer r.mut.Unlock()

	if user, ok := r.users[id]; ok {
		return user
	}

	return nil
}

// RoomUser describes a user in a room.
type RoomUser struct {
	id     int
	name   string
	joined chan *RoomUser // Channel receives user id when a user joins
	leave  chan int
}

func newRoomUser(id int, name string) *RoomUser {
	return &RoomUser{
		id:     id,
		name:   name,
		joined: make(chan *RoomUser),
		leave:  make(chan int),
	}
}

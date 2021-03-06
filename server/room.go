package main

import (
	"context"
	"errors"
	"fmt"
)

// Room represents a game room.
type Room struct {
	Users map[int]*RoomUser
}

func newRoom() Room {
	return Room{
		Users: make(map[int]*RoomUser),
	}
}

func (r *Room) isFull() bool {
	return len(r.Users) >= roomSize
}

func (r *Room) joinRoom(user *RoomUser) error {
	roomIsFull := r.isFull()
	if roomIsFull {
		return errors.New("Room is full")
	}

	if _, ok := r.Users[user.id]; ok {
		return errors.New("User is already in the room")
	}

	r.Users[user.id] = user

	return nil
}

func (r *Room) run() {
	for {
		select {}
	}
}

// RoomUser describes a user in a room.
type RoomUser struct {
	id     int
	name   string
	joined chan int // Channel receives user id when a user joins
	leave  chan int
	cancel context.CancelFunc
}

func newRoomUser(id int, name string) *RoomUser {
	return &RoomUser{
		id:     id,
		name:   name,
		joined: make(chan int),
		leave:  make(chan int),
	}
}

func (r *RoomUser) run() {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	select {
	case id := <-r.joined:
		fmt.Println("joined", id)

	case id := <-r.leave:
		fmt.Println("leave", id)

	case <-ctx.Done():
		fmt.Println("room user closed", r.id)
		return
	}
}

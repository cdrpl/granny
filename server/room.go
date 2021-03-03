package main

import "errors"

// Room represents a game room.
type Room struct {
	Users []User `json:"users"`
}

func newRoom() Room {
	return Room{
		Users: make([]User, 0),
	}
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

func (r *Room) joinRoom(user User) error {
	roomIsFull := r.isFull()
	if roomIsFull {
		return errors.New("Room is full")
	}

	if r.hasUser(user.ID) {
		return errors.New("User is already in the room")
	}

	r.Users = append(r.Users, user)
	return nil
}

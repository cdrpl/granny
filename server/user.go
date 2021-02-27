package main

import (
	"sync"
)

// User account data.
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// UserManager caches user data.
type UserManager struct {
	users map[int64]*User
	mutex sync.Mutex
}

// CreateUserManager will create and return a Manager instance.
func CreateUserManager() *UserManager {
	return &UserManager{
		users: make(map[int64]*User),
	}
}

// Register will query the user data add it to the users map.
func (m *UserManager) Register(user *User) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.users[user.ID] = user
}

// Unregister will remove the user from the users map.
func (m *UserManager) Unregister(id int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.users, id)
}

// HasUser will return true if the user is in the users map.
func (m *UserManager) HasUser(id int64) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.users[id]
	return ok
}

package models

// User represents a user in the system
type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// UserResponse is the standard format for User responses
type UserResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
	Users   []User `json:"users,omitempty"`
}

// NewUser creates a new user with the given id and name
func NewUser(id int, name string) *User {
	return &User{
		ID:   id,
		Name: name,
	}
}

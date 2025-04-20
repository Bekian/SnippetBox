package models

import (
	"database/sql"
	"time"
)

// model struct for users table
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// wraps db connection pool
type UserModel struct {
	DB *sql.DB
}

// // methods below are placeholders for now
// insert user into users table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// validate user information to log the user in,
// returns userID if successful
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// checks if a user exists by ID
func (m *UserModel) Exitsts(id int) (bool, error) {
	return false, nil
}

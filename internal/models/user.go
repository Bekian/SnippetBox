package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// check for specific sql err type
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			// if the error relates to the email value
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				// return duplicate email error
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// validate user information to log the user in,
// returns userID if successful
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	// search for user with provided password
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// user not found with provided email
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// check if passwords match
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			// wrong pass
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// successful login
	return id, nil
}

// checks if a user exists by ID
func (m *UserModel) Exitsts(id int) (bool, error) {
	return false, nil
}

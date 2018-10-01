package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateName      = errors.New("models: name or email address already in use")
	ErrInvalidCredentials = errors.New("models: invalid user credentials")
)

type User struct {
	ID       int
	Name     string
	Password string
	Created  time.Time
}

func (user *User) Valid() error {
	if user.Name == "" || user.Password == "" {
		return errors.New(fmt.Sprintf("User data incorrect!\n=== %v ===", user))
	}
	return nil
}

func (db *Database) InsertUser(user *User) error {
	err := user.Valid()
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, password, created)
	VALUES (?, ?, UTC_TIMESTAMP())`
	_, err = db.Exec(stmt, user.Name, hashedPassword)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return ErrDuplicateName
		}
	}
	return err
}

func (db *Database) ChangeUserPassword(username string, password string) error {
	if username == "" {
		return errors.New("Empty Username")
	}

	if password == "" {
		return errors.New("Empty Password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	result, err := db.Exec(`UPDATE users SET password = ? WHERE name = ?`, hashedPassword, username)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("ChangeUserPassword %d Rows Affected", rowsAffected)

	return err
}

func (db *Database) VerifyUser(name, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := db.QueryRow("SELECT id, password FROM users WHERE name = ?", name)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) UserInfo(userID int) (*User, error) {
	user := &User{}
	row := db.QueryRow("SELECT id, name FROM users WHERE id = ?", userID)
	err := row.Scan(&user.ID, &user.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

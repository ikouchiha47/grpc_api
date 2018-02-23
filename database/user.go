package database

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
)

type User struct {
	ID          int64
	Name        string
	Email       string
	Password    string
	AccessToken string
}

type Database struct {
	Name  string
	Users []*User
}

func NewDatabase(name string) *Database {
	return &Database{
		Name:  name,
		Users: []*User{&User{ID: 0}},
	}
}

func (u *User) HasValidParams() bool {
	return u.Name != "" && u.Email != "" && u.Password != ""
}

func (d *Database) FindUser(user *User) (*User, error) {
	for _, u := range d.Users {
		if u.ID == 0 {
			continue
		}
		if u.ID == user.ID || u.Email == user.Email {
			return u, nil
		}
	}
	return nil, errors.New("NOT_FOUND")
}

func (d *Database) AddUser(user *User) (*User, error) {
	if !user.HasValidParams() {
		return nil, errors.New("PARAMS_MISSING")
	}

	if _, err := d.FindUser(user); err == nil {
		return nil, errors.New("USER_EXISTS")
	}

	hpass := GenerateHashPassword(user.Password)
	lastID := d.Users[len(d.Users)-1].ID
	user.ID = lastID + 1
	user.Password = hpass
	d.Users = append(d.Users, user)

	return user, nil
}

func (d *Database) ValidateUser(user *User) (bool, error) {
	dbUser, err := d.FindUser(user)
	if err != nil {
		return false, err
	}

	hpass := GenerateHashPassword(user.Password)
	return hpass == dbUser.Password, nil
}

func (d *Database) CreateToken(user *User) (*User, error) {
	var err error
	user.AccessToken, err = GenerateHash()
	return user, err
}

func (d *Database) RemoveToken(user *User) (*User, error) {
	user.AccessToken = ""
	return user, nil
}

func (d *Database) ListUsers() ([]*User, error) {
	return d.Users[1:], nil
}

func GenerateHash() (string, error) {
	b := make([]byte, 60)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func GenerateHashPassword(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

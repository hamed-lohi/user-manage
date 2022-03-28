package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Role uint

const (
	Unknown Role = iota
	Guest
	Member
	Moderator
	Admin
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password"`
	Bio      *string            //`bson:"bio"`
	Roles    []Role             `bson:"roles,omitempty"`
	// Image      *string
}

func (u *User) SetPassword(plain string) error {
	h, err := u.HashPassword(plain)
	u.Password = h
	return err
}

func (u *User) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}

package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}
	return m
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := make(map[string]string)
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf(
			"firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf(
			"lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf(
			"password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("invalid e-mail")
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID                string `bson:"_id,omitempty"    json:"id,omitempty"`
	FirstName         string `bson:"firstName"        json:"firstName"`
	LastName          string `bson:"lastName"         json:"lastName"`
	Email             string `bson:"email"            json:"email"`
	EncryptedPassword string `bson:"encyptedPassword" json:"-"`
	IsAdmin           bool   `bson:"isAdmin"          json:"isAdmin"`
}

func NewUserFromParams(p CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         p.FirstName,
		LastName:          p.LastName,
		Email:             p.Email,
		EncryptedPassword: string(encpw),
	}, nil
}

func IsValidPassword(encryptedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password)) == nil
}

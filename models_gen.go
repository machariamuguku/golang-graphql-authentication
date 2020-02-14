// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package golang_graphql_authentication

import (
	"time"
)

type FieldErrors struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserPayload struct {
	User        *User          `json:"user"`
	JwtToken    *string        `json:"jwtToken"`
	StatusCode  int            `json:"statusCode"`
	Message     string         `json:"message"`
	FieldErrors []*FieldErrors `json:"fieldErrors"`
}

type RegisterUserInput struct {
	FirstName                    string `json:"firstName"`
	LastName                     string `json:"lastName"`
	Email                        string `json:"email"`
	PhoneNumber                  string `json:"phoneNumber"`
	Password                     string `json:"password"`
	EmailVerificationCallBackURL string `json:"emailVerificationCallBackURL"`
}

type RegisterUserPayload struct {
	User        *User          `json:"user"`
	JwtToken    *string        `json:"jwtToken"`
	StatusCode  int            `json:"statusCode"`
	Message     string         `json:"message"`
	FieldErrors []*FieldErrors `json:"fieldErrors"`
}

type User struct {
	ID              string    `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email"`
	PhoneNumber     string    `json:"phoneNumber"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	IsPhoneVerified bool      `json:"isPhoneVerified"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

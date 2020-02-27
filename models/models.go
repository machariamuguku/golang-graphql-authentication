package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GormUser : structure for postgres user table
type GormUser struct {
	ID                           string `json:"id" gorm:"column:id;PRIMARY_KEY"`
	FirstName                    string `json:"firstName" validate:"required"`
	LastName                     string `json:"lastName" validate:"required"`
	Email                        string `json:"email" validate:"required,email" gorm:"type:varchar(100);unique_index"`
	PhoneNumber                  string `json:"phoneNumber" validate:"required"`
	Password                     string `json:"password" validate:"required,min=8"`
	EmailVerificationCallBackURL string `json:"emailVerificationCallBackURL" validate:"required"`
	IsEmailVerified              bool   `json:"isEmailVerified"`
	IsPhoneVerified              bool   `json:"isPhoneVerified"`
	EmailVerificationToken       string `json:"emailVerificationToken" gorm:"unique"`
	PhoneVerificationToken       string `json:"phoneVerificationToken"`
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	DeletedAt                    *time.Time `sql:"index"`
}

//Claims : struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserID string `json:"userid"`
	jwt.StandardClaims
}

// ValidateLoginInput : login validation struct
type ValidateLoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

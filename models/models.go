package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GormUser : for postgres
type GormUser struct {
	ID          string `gorm:"column:id; PRIMARY_KEY" json:"id"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Email       string `gorm:"type:varchar(100);unique_index" json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Password    string `json:"password" validate:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

//Claims : struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

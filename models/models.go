package models

import (
	"time"
)

// GormUser : for postgres
type GormUser struct {
	ID          string `gorm:"column:id; PRIMARY_KEY" json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `gorm:"type:varchar(100);unique_index" json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

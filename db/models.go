package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Otto struct {
	gorm.Model

	Serial string `gorm:"not null;unique_index"`

	Manufactured  time.Time
	ProductionLot string

	OTPSeed []byte
	Owner   *User
}

type User struct {
	gorm.Model

	Username string `gorm:"not null;unique_index"`

	LastSeenOnline   time.Time
	LastConversation time.Time
}

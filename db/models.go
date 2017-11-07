package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Otto struct {
	gorm.Model

	UserID    *uint
	Serial    string `gorm:"not null;unique_index"`
	OTPSecret string

	Manufactured  time.Time
	ProductionLot string
}

type User struct {
	gorm.Model

	Username string `gorm:"not null;unique_index"`
	Ottos    []Otto `gorm:"ForeignKey:UserID"`

	LastSeenOnline   time.Time
	LastConversation time.Time
}

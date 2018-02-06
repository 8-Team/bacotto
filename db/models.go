package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Otto struct {
	gorm.Model

	DeviceModel string
	Revision    int

	MACAddress []byte
	Serial     string `gorm:"not null;unique_index"`
	OTPSecret  string

	Manufactured  time.Time
	ProductionLot string

	UserID *uint
}

type User struct {
	gorm.Model

	Username       string          `gorm:"not null;unique_index"`
	Ottos          []*Otto         `gorm:"ForeignKey:UserID"`
	ProjectEntries []*ProjectEntry `gorm:"ForeignKey:UserID"`
	Projects       []*Project      `gorm:"many2many:user_projects;"`

	LastSeenOnline   time.Time
	LastConversation time.Time
}

type Project struct {
	gorm.Model

	Name      string `gorm:"not null;unique_index"`
	ShortName string
}

type ProjectEntry struct {
	gorm.Model

	UserID      uint
	ProjectID   uint
	Description string

	StartTime time.Time
	EndTime   time.Time
}

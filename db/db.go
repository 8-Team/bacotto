package db

import (
	"github.com/jinzhu/gorm"
)

type DB struct {
	*gorm.DB
}

func Open(uri string) (*DB, error) {
	db, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Otto{})

	return &DB{
		DB: db,
	}, nil
}

func (db *DB) Close() {
	db.Close()
}

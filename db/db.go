package db

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *Database

func Open(uri string) error {
	gdb, err := gorm.Open("postgres", uri)
	if err != nil {
		return err
	}

	gdb.DB().SetConnMaxLifetime(time.Minute * 5)
	gdb.DB().SetMaxIdleConns(0)
	gdb.DB().SetMaxOpenConns(5)

	gdb.AutoMigrate(&Otto{}, &User{})

	DB = &Database{DB: gdb}

	if _, y := os.LookupEnv("MOCK_DB"); y {
		fillMockData()
	}

	return nil
}

func Close() {
	DB.Close()
}

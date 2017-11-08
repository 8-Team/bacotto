package db

import (
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

	DB = &Database{DB: gdb}

	// Setup some sane DB connection parameters
	gdb.DB().SetConnMaxLifetime(time.Minute * 5)
	gdb.DB().SetMaxIdleConns(0)
	gdb.DB().SetMaxOpenConns(5)

	// Migrate schema and manually create foreign keys
	gdb.AutoMigrate(&Otto{}, &User{}, &Project{}, &ProjectEntry{})

	gdb.Model(&Otto{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	gdb.Model(&ProjectEntry{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	gdb.Model(&ProjectEntry{}).AddForeignKey("project_id", "projects(id)", "RESTRICT", "RESTRICT")

	return nil
}

func Close() {
	DB.Close()
}

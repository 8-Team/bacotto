package db

import (
	"time"

	"github.com/8-team/bacotto/conf"
	"github.com/8-team/bacotto/erp"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

type Database struct {
	*gorm.DB
}

var log *logrus.Entry
var DB *Database

func init() {
	if conf.DebugLogLevel() {
		logrus.SetLevel(logrus.DebugLevel)
	}
	log = logrus.WithField("app", "db")
}

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

	// When first starting, try to sync with the ERP DB
	projects, err := erp.ListProjects()
	if err != nil {
		log.Warnf("Could not retrieve projects from ERP, skipping sync")
	} else {
		for _, prj := range projects {
			InsertProject(&Project{
				Name:      prj.Name,
				ShortName: prj.Name,
			})
		}
	}

	return nil
}

func Close() {
	DB.Close()
}

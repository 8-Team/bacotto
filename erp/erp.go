package erp

import (
	"errors"
	"time"

	"github.com/8-team/bacotto/conf"
)

type Project struct {
	ID   uint
	Name string
}

type ProjectEntry struct {
	ID          uint
	ProjectID   uint
	Description string
	Date        time.Time
	Duration    time.Duration
}

type erper interface {
	listProjects() ([]Project, error)
	addProjectEntry(pid uint, description string, date time.Time, duration time.Duration) (ProjectEntry, error)
	listProjectEntries(pid uint, from time.Time, to time.Time) ([]ProjectEntry, error)
	removeProjectEntry(pe ProjectEntry) error
}

var erp erper

func Open() (err error) {
	if conf.UseMockERP() {
		erp, err = NewMockERP()
	} else {
		err = errors.New("only mock ERP supported for now")
	}
	return
}

func Close() {

}

func ListProjects() ([]Project, error) {
	return erp.listProjects()
}

func AddProjectEntry(pid uint, description string, date time.Time, duration time.Duration) (ProjectEntry, error) {
	return erp.addProjectEntry(pid, description, date, duration)
}

func ListProjectEntries(pid uint, from time.Time, to time.Time) ([]ProjectEntry, error) {
	return erp.listProjectEntries(pid, from, to)
}

func RemoveProjectEntry(pe ProjectEntry) error {
	return erp.removeProjectEntry(pe)
}

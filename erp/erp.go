package erp

import "time"

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

type Erper interface {
	ListProjects() ([]Project, error)
	AddProjectEntry(pid uint, description string, date time.Time, duration time.Duration) (ProjectEntry, error)
	ListProjectEntries(pid uint, from time.Time, to time.Time) ([]ProjectEntry, error)
	RemoveProjectEntry(pe ProjectEntry) error
}

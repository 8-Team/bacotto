package erp

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/8-team/bacotto/conf"
)

type MockERP struct {
	projects map[uint]Project
	entries  map[uint]ProjectEntry

	nextPID uint
	nextEID uint
}

func NewMockERP() (*MockERP, error) {
	erp := &MockERP{
		projects: make(map[uint]Project),
		entries:  make(map[uint]ProjectEntry),
	}

	return erp, erp.loadData(conf.MockERPDataPath())
}

func (erp *MockERP) loadData(path string) error {
	var mockData = struct {
		Projects []Project
		Entries  []ProjectEntry
		PID      uint
		EID      uint
	}{}

	r, err := os.Open(path)
	if err != nil {
		return err
	}

	d := json.NewDecoder(r)

	if err := d.Decode(&mockData); err != nil {
		return err
	}

	for _, p := range mockData.Projects {
		erp.projects[p.ID] = p
	}

	for _, e := range mockData.Entries {
		erp.entries[e.ID] = e
	}

	erp.nextPID = mockData.PID
	erp.nextEID = mockData.EID

	return nil
}

func (erp *MockERP) listProjects() ([]Project, error) {
	prjs := make([]Project, 0, len(erp.projects))
	for _, prj := range erp.projects {
		prjs = append(prjs, prj)
	}
	return prjs, nil
}

func (erp *MockERP) addProjectEntry(pid uint, description string, date time.Time, duration time.Duration) (entry ProjectEntry, err error) {
	if _, ok := erp.projects[pid]; !ok {
		err = errors.New("invalid project ID")
		return
	}

	if description == "" {
		err = errors.New("description cannot be empty")
		return
	}

	entry.ID = erp.nextEID
	entry.ProjectID = pid
	entry.Description = description
	entry.Date = date
	entry.Duration = duration

	erp.entries[entry.ID] = entry
	erp.nextEID++

	return
}

func (erp *MockERP) listProjectEntries(pid uint, from time.Time, to time.Time) ([]ProjectEntry, error) {
	entries := make([]ProjectEntry, 0)

	for _, e := range erp.entries {
		if e.ProjectID == pid && e.Date.After(from) && e.Date.Before(to) {
			entries = append(entries, e)
		}
	}

	return entries, nil
}

func (erp *MockERP) removeProjectEntry(pe ProjectEntry) error {
	if _, exists := erp.entries[pe.ID]; !exists {
		return errors.New("entry does not exist")
	}

	delete(erp.entries, pe.ID)
	return nil
}

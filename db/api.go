package db

import (
	"fmt"
	"time"

	"github.com/8-team/bacotto/conf"
	"github.com/pquerna/otp/totp"
)

func GetOtto(serial string) (*Otto, error) {
	otto := new(Otto)

	if err := DB.First(otto, "serial = ?", serial).Error; err != nil {
		return otto, fmt.Errorf("Could not retrieve serial (%s): %s", serial, err)
	}

	return otto, nil
}

func GetUserFromSerial(serial string) (*User, error) {
	user := new(User)

	otto, err := GetOtto(serial)
	if err != nil {
		return nil, err
	}

	if err := DB.First(user, "id = ?", otto.UserID).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserProjects(u *User) ([]*Project, error) {
	if err := DB.Preload("Projects").First(u, "id = ?", u.ID).Error; err != nil {
		return nil, err
	}
	return u.Projects, nil
}

func GetUserEntries(u *User, from time.Time, to time.Time) ([]*ProjectEntry, error) {
	if err := DB.Preload("ProjectEntries").First(u, "id = ?", u.ID).Error; err != nil {
		return nil, err
	}
	return u.ProjectEntries, nil
}

func AddProjectToUser(u *User, p *Project) error {
	u.Projects = append(u.Projects, p)
	return DB.Save(u).Error
}

func InsertProject(p *Project) error {
	return DB.FirstOrCreate(p, "name = ?", p.Name).Error
}

// Authorize authenticates a request from an Otto by using the provided OTP code.
func Authorize(serial string, otp string) (*Otto, error) {
	otto, err := GetOtto(serial)
	if err != nil {
		return nil, err
	}

	if !conf.VerifyMsg() {
		return otto, nil
	}

	if !totp.Validate(otp, otto.OTPSecret) {
		return otto, fmt.Errorf("invalid OTP code %s from %s", otp, serial)
	}

	return otto, nil
}

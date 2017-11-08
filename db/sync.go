package db

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Sync retrieves a CSV file located at the specified URI containing
// Otto serial numbers and merges it into the database.
func Sync(uri string) error {
	var r io.Reader

	pURI, err := url.Parse(uri)
	if err != nil {
		return err
	}

	switch pURI.Scheme {
	case "":
		r, err = os.Open(pURI.Path)
		if err != nil {
			return err
		}
	default:
		return errors.New("URI schema not yet supported")
	}

	records, err := readRecordsFromCSV(r)
	if err != nil {
		return err
	}

	for _, rec := range records {
		if err := DB.FirstOrCreate(&rec, "serial = ?", rec.Serial).Error; err != nil {
			return err
		}
	}

	return nil
}

func readRecordsFromCSV(r io.Reader) ([]Otto, error) {
	rd := csv.NewReader(r)
	rd.Comma = ';' // CSV separator

	// skip CSV header
	_, err := rd.Read()
	if err != nil {
		return nil, err
	}

	records, err := rd.ReadAll()
	if err != nil {
		return nil, err
	}

	ottos := make([]Otto, len(records))
	for i, rec := range records {
		rev, err := strconv.Atoi(rec[1])
		if err != nil {
			return nil, fmt.Errorf("invalid revision at record %d: %s", i, err)
		}

		mac, err := net.ParseMAC(rec[3])
		if err != nil {
			return nil, fmt.Errorf("invalid MAC at record %d: %s", i, err)
		}

		ottos[i] = Otto{
			Serial:        rec[0],
			Revision:      rev,
			DeviceModel:   rec[2],
			MACAddress:    mac,
			OTPSecret:     rec[4],
			Manufactured:  time.Now(),
			ProductionLot: rec[6],
		}
	}

	return ottos, nil
}

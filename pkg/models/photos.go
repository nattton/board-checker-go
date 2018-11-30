package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"gitlab.com/code-mobi/board-checker/pkg/forms"
)

func (db *Database) GetAutoRunningNumber(worksheetID int) (next int) {
	stmt := `SELECT MAX(running_number) FROM photos WHERE worksheet_id = ?`
	db.QueryRow(stmt, worksheetID).Scan(&next)
	next++
	return
}

func (db *Database) InsertPhoto(f *Photo) error {
	if f.RunningNumber < 1 {
		f.RunningNumber = db.GetAutoRunningNumber(f.WorksheetID)
	}

	stmt := `INSERT INTO photos (worksheet_id, running_number, filename, location, created)
	VALUES (?, ?, ?, ?, UTC_TIMESTAMP())`
	_, err := db.Exec(stmt, f.WorksheetID, f.RunningNumber, f.FileName, f.Location)
	return err
}

func (db *Database) ListPhotos(worksheetID int, q *forms.Query) (Photos, error) {
	stmt := `SELECT id, worksheet_id, running_number, filename, location, created FROM photos WHERE worksheet_id = ? ORDER BY id ASC`
	rows, err := db.Query(stmt, worksheetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photos := Photos{}
	for rows.Next() {
		f := &Photo{}
		rows.Scan(&f.ID, &f.WorksheetID, &f.RunningNumber, &f.FileName, &f.Location, &f.Created)
		if err != nil {
			return nil, err
		}
		photos = append(photos, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}
func (db *Database) ListPhotosMaps(worksheetID int, storeDir string) (Locations, error) {
	stmt := `SELECT id, worksheet_id, running_number, filename, location, created FROM photos WHERE worksheet_id = ? ORDER BY id ASC`
	rows, err := db.Query(stmt, worksheetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photos := Photos{}
	for rows.Next() {
		f := &Photo{}
		rows.Scan(&f.ID, &f.WorksheetID, &f.RunningNumber, &f.FileName, &f.Location, &f.Created)
		if err != nil {
			return nil, err
		}
		photos = append(photos, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	locations := Locations{}
	for _, photo := range photos {
		photoPath := storeDir + "/" + strconv.Itoa(photo.WorksheetID) + "/" + photo.FileName
		log.Println(photoPath)

		f, err := os.Open(photoPath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		exif.RegisterParsers(mknote.All...)
		x, err := exif.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		lat, long, err := x.LatLong()
		fmt.Printf("lat, long: %f %f\n", lat, long)
		if err != nil {
			continue
		}

		lo := &Location{lat, long}
		locations = append(locations, lo)
	}
	locationStr, err := json.Marshal(locations)
	if err != nil {
		log.Print(err.Error())
	}
	log.Print(string(locationStr))
	return locations, nil
}

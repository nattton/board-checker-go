package models

import "gitlab.com/code-mobi/board-checker/pkg/forms"

func (db *Database) GetAutoRunningNumber(projectID int) (next int) {
	stmt := `SELECT MAX(running_number) FROM photos WHERE project_id = ?`
	db.QueryRow(stmt, projectID).Scan(&next)
	next++
	return
}

func (db *Database) InsertPhoto(f *Photo) error {
	if f.RunningNumber < 1 {
		f.RunningNumber = db.GetAutoRunningNumber(f.ProjectID)
	}

	stmt := `INSERT INTO photos (project_id, running_number, filename, location, created)
	VALUES (?, ?, ?, ?, UTC_TIMESTAMP())`
	_, err := db.Exec(stmt, f.ProjectID, f.RunningNumber, f.FileName, f.Location)
	return err
}

func (db *Database) ListPhotos(projectID int, q *forms.Query) (Photos, error) {
	stmt := `SELECT id, project_id, running_number, filename, location, created FROM photos WHERE project_id = ? ORDER BY id ASC`
	rows, err := db.Query(stmt, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photos := Photos{}
	for rows.Next() {
		f := &Photo{}
		rows.Scan(&f.ID, &f.ProjectID, &f.RunningNumber, &f.FileName, &f.Location, &f.Created)
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

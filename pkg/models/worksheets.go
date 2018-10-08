package models

import "database/sql"

func (db *Database) ListWorksheets() (Worksheets, error) {
	stmt := `SELECT id, no, name, created FROM worksheets ORDER BY created DESC`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	worksheets := Worksheets{}
	for rows.Next() {
		p := &Worksheet{}
		rows.Scan(&p.ID, &p.Name, &p.Created)
		if err != nil {
			return nil, err
		}
		worksheets = append(worksheets, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return worksheets, nil
}

func (db *Database) GetWorksheet(id int) (*Worksheet, error) {
	stmt := `SELECT id, name, created FROM worksheets WHERE id = ?`
	row := db.QueryRow(stmt, id)

	p := &Worksheet{}
	err := row.Scan(&p.ID, &p.Name, &p.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *Database) InsertWorksheet(worksheet *Worksheet) error {
	stmt := `INSERT INTO worksheets (id, name, created) VALUES (?, ?, UTC_TIMESTAMP())`
	_, err := db.Exec(stmt, worksheet.ID, worksheet.Name)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateWorksheet(worksheet *Worksheet) error {
	stmt := `UPDATE worksheets SET name = ? WHERE id = ?`
	_, err := db.Exec(stmt, worksheet.Name, worksheet.ID)
	if err != nil {
		return err
	}
	return nil
}

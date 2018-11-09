package models

import "database/sql"

func (db *Database) ListDistinctDate() ([]string, error) {
	stmt := `SELECT DISTINCT date_format(created, '%Y-%m-%d') as uniquedates 
	FROM worksheets 
	GROUP BY date_format(created, '%Y-%m-%d') 
	ORDER BY uniquedates DESC`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	listDate := []string{}
	for rows.Next() {
		date := ""
		rows.Scan(&date)
		if err != nil {
			return nil, err
		}
		listDate = append(listDate, date)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listDate, nil
}

func (db *Database) ListWorksheets() (Worksheets, error) {
	stmt := `SELECT id, number, name, created 
	FROM worksheets 
	ORDER BY created DESC`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	worksheets := Worksheets{}
	for rows.Next() {
		p := &Worksheet{}
		rows.Scan(&p.ID, &p.Number, &p.Name, &p.Created)
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

func (db *Database) ListWorksheetsByDate(date string) (Worksheets, error) {
	stmt := `SELECT id, number, name, created 
	FROM worksheets 
	WHERE date_format(created, '%Y-%m-%d') = ? 
	ORDER BY created DESC`
	rows, err := db.Query(stmt, date)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	worksheets := Worksheets{}
	for rows.Next() {
		p := &Worksheet{}
		rows.Scan(&p.ID, &p.Number, &p.Name, &p.Created)
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

func (db *Database) ListWorksheetsByZone(zoneID int) (Worksheets, error) {
	stmt := `SELECT id, number, name, created FROM worksheets WHERE zone_id = ? ORDER BY created DESC`
	rows, err := db.Query(stmt, zoneID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	worksheets := Worksheets{}
	for rows.Next() {
		p := &Worksheet{}
		rows.Scan(&p.ID, &p.Number, &p.Name, &p.Created)
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

func (db *Database) ListWorksheetsByTeam(teamID int) (Worksheets, error) {
	stmt := `SELECT id, number, name, created FROM worksheets WHERE team_id = ? ORDER BY created DESC`
	rows, err := db.Query(stmt, teamID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	worksheets := Worksheets{}
	for rows.Next() {
		p := &Worksheet{}
		rows.Scan(&p.ID, &p.Number, &p.Name, &p.Created)
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
	stmt := `SELECT w.id, w.number, w.name, w.created, z.id zone_id, z.name zone_name, t.id team_id, t.name team_name FROM worksheets w 
	INNER JOIN zones z on (w.zone_id = z.id) 
	INNER JOIN teams t on (w.team_id = t.id) 
	WHERE w.id = ?`
	row := db.QueryRow(stmt, id)

	p := &Worksheet{}
	err := row.Scan(&p.ID, &p.Number, &p.Name, &p.Created, &p.ZoneID, &p.ZoneName, &p.TeamID, &p.TeamName)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *Database) InsertWorksheet(worksheet *Worksheet) error {
	stmt := `INSERT INTO worksheets (number, name, zone_id, team_id, created) VALUES (?, ?, ?, ?, UTC_TIMESTAMP())`
	result, err := db.Exec(stmt, worksheet.Number, worksheet.Name, worksheet.ZoneID, worksheet.TeamID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	worksheet.ID = int(id)
	return err
}

func (db *Database) UpdateWorksheet(worksheet *Worksheet) error {
	stmt := `UPDATE worksheets SET number = ?, name = ?, zone_id = ?, team_id = ? WHERE id = ?`
	_, err := db.Exec(stmt, worksheet.Number, worksheet.Name, worksheet.ID, worksheet.ZoneID, worksheet.TeamID)
	if err != nil {
		return err
	}
	return nil
}

package models

import "database/sql"

func (db *Database) ListZones() (Zones, error) {
	stmt := `SELECT id, name FROM zones ORDER BY name`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	zones := Zones{}
	for rows.Next() {
		t := &Zone{}
		rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		zones = append(zones, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return zones, nil
}

func (db *Database) GetZone(id int) (*Zone, error) {
	stmt := `SELECT id, name FROM zones WHERE id = ?`
	row := db.QueryRow(stmt, id)

	t := &Zone{}
	err := row.Scan(&t.ID, &t.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return t, nil
}

func (db *Database) InsertZone(zone *Zone) error {
	stmt := `INSERT INTO zones (name) VALUES (?)`
	_, err := db.Exec(stmt, zone.Name)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateZone(zone *Zone) error {
	stmt := `UPDATE zones SET name = ? WHERE id = ?`
	_, err := db.Exec(stmt, zone.Name, zone.ID)
	if err != nil {
		return err
	}
	return nil
}

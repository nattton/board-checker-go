package models

import "database/sql"

func (db *Database) ListTeams() (Teams, error) {
	stmt := `SELECT id, name FROM teams ORDER BY name`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	teams := Teams{}
	for rows.Next() {
		t := &Team{}
		rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func (db *Database) GetTeam(id int) (*Team, error) {
	stmt := `SELECT id, name FROM worksheets WHERE id = ?`
	row := db.QueryRow(stmt, id)

	t := &Team{}
	err := row.Scan(&t.ID, &t.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return t, nil
}

func (db *Database) InsertTeam(team *Team) error {
	stmt := `INSERT INTO teams (id, name) VALUES (?, ?)`
	_, err := db.Exec(stmt, team.ID, team.Name)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateTeam(team *Team) error {
	stmt := `UPDATE teams SET name = ? WHERE id = ?`
	_, err := db.Exec(stmt, team.Name, team.ID)
	if err != nil {
		return err
	}
	return nil
}

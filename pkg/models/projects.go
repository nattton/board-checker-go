package models

import "database/sql"

func (db *Database) ListProjects() (Projects, error) {
	stmt := `SELECT id, name, created FROM projects ORDER BY created DESC`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projects := Projects{}
	for rows.Next() {
		p := &Project{}
		rows.Scan(&p.ID, &p.Name, &p.Created)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (db *Database) GetProject(id int) (*Project, error) {
	stmt := `SELECT id, name, created FROM projects WHERE id = ?`
	row := db.QueryRow(stmt, id)

	p := &Project{}
	err := row.Scan(&p.ID, &p.Name, &p.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *Database) InsertProject(project *Project) error {
	stmt := `INSERT INTO projects (id, name, created) VALUES (?, ?, UTC_TIMESTAMP())`
	_, err := db.Exec(stmt, project.ID, project.Name)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateProject(project *Project) error {
	stmt := `UPDATE projects SET name = ? WHERE id = ?`
	_, err := db.Exec(stmt, project.Name, project.ID)
	if err != nil {
		return err
	}
	return nil
}

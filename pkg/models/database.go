package models

import (
	"database/sql"
)

type Database struct {
	*sql.DB
}

func (db *Database) CreateTable() error {
	_, err := db.Exec(`CREATE TABLE users (
		id int(11) NOT NULL AUTO_INCREMENT,
		name varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
		password char(60) COLLATE utf8mb4_general_ci NOT NULL,
		created datetime NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE projects (
		id int(11) NOT NULL,
		name varchar(255) NOT NULL,
		created datetime NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
	  `)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE photos (
		id int(11) NOT NULL AUTO_INCREMENT,
		project_id int(11) NOT NULL,
		running_number int(5) NOT NULL,
		filename varchar(255) NOT NULL,
		created datetime NOT NULL,
		location varchar(255) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;`)
	if err != nil {
		return err
	}

	return err
}

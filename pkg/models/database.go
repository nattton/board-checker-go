package models

import (
	"database/sql"
)

type Database struct {
	*sql.DB
}

func (db *Database) CreateTable() error {

	_, err := db.Exec(`CREATE TABLE zones (
		id int(11) NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
	
	  CREATE TABLE photos (
		id int(11) NOT NULL AUTO_INCREMENT,
		worksheet_id int(11) NOT NULL,
		running_number int(5) NOT NULL,
		filename varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
		created datetime NOT NULL,
		location varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
		photoscol varchar(45) COLLATE utf8mb4_general_ci DEFAULT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
	
	  CREATE TABLE teams (
		id int(11) NOT NULL AUTO_INCREMENT,
		name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
	
	
	CREATE TABLE worksheets (
		id int(11) NOT NULL AUTO_INCREMENT,
		number varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
		team_id int(11) NOT NULL,
		zone_id int(11) NOT NULL,
		name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
		created datetime NOT NULL,
		PRIMARY KEY (id,number),
		UNIQUE KEY number_UNIQUE (number)
	  ) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
	  
	
	CREATE TABLE users (
		id int(11) NOT NULL AUTO_INCREMENT,
		name varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
		password char(60) COLLATE utf8mb4_general_ci NOT NULL,
		created datetime NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
	`)

	if err != nil {
		return err
	}

	return err
}

package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"gitlab.com/code-mobi/board-checker/pkg/models"
)

func main() {
	dsn := flag.String("dsn", os.Getenv("BC_DSN"), "Database DSN")
	cmd := flag.String("cmd", "", `Command
	adduser -name -password
	changepwd -name -password`)

	name := flag.String("name", "", "User Name")
	password := flag.String("password", "", "User Password")

	flag.Parse()

	db := connect(*dsn)

	database := &models.Database{db}

	switch *cmd {
	case "migrate":
		err := database.CreateTable()
		if err != nil {
			log.Fatal(err)
		}
	case "adduser":
		user := &models.User{
			Name:     *name,
			Password: *password,
		}
		log.Printf("Add User %v", user)
		err := database.InsertUser(user)
		if err != nil {
			log.Fatal(err)
		}
	case "changepwd":
		err := database.ChangeUserPassword(*name, *password)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func connect(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

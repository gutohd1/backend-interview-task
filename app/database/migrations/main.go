package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	fmt.Println("Running migrations")
	db, err := sql.Open("mysql", "app:password@tcp(127.0.0.1:33306)/muzzapp?multiStatements=true")

	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to the DB")
	driver, _ := mysql.WithInstance(db, &mysql.Config{})

	m, _ := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)

	var command string
	flag.StringVar(&command, "migrate", "", "UP, DOWN. command to migrate up or down")
	flag.Parse()

	switch command {
	case "down":
		fmt.Println("Migrations will go Down")
		if err := m.Down(); err != nil {
			fmt.Printf("error: %s", err)
		}
	case "up":
		fmt.Println("Migrations will go Up")
		if err := m.Up(); err != nil {
			fmt.Print(err)
		}
	default:
		//implement migrate
	}
}

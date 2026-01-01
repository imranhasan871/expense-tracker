package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Student struct {
	id         int
	name       string
	email      string
	bloodgroup string
	contact    string
	address    string
}

func main() {

	connStr := "user=admin password=root host=localhost port=5432 dbname=expense_tracker sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}

	// err := db.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	rows, err := db.Query("SELECT * FROM students")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		log.Println(rows)
	}

	defer db.Close()

	// Now I want to insert a new student direct put values now

	insert, err := db.Prepare("INSERT INTO students(name, email, bloodgroup, contact, address) VALUES('Resun Akondo', 'resun.akondo@example.com', 'A+', '01712345678', 'Dhaka, Bangladesh')")
	if err != nil {
		log.Fatal(err)
	}
	defer insert.Close()

	insert.Exec()
	log.Println("Student inserted successfully")

}

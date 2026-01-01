package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type StudentRepo struct{}

var connStr = "user=admin password=root host=localhost port=5432 dbname=expense_tracker sslmode=disable"

func (s StudentRepo) insertStudent(student Student) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	name := student.name
	email := student.email
	bloodgroup := student.bloodgroup
	contact := student.contact
	address := student.address

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	bloodgroup = strings.TrimSpace(bloodgroup)
	contact = strings.TrimSpace(contact)
	address = strings.TrimSpace(address)

	insert, err := db.Prepare("INSERT INTO students(name, email, bloodgroup, contact, address) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		log.Fatal(err)
	}
	defer insert.Close()

	_, err = insert.Exec(name, email, bloodgroup, contact, address)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Student inserted successfully")
}

func (s StudentRepo) getStudents() []Student {

	var students []Student

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM students")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var student Student
		err := rows.Scan(&student.id, &student.name, &student.email, &student.bloodgroup, &student.contact, &student.address)
		if err != nil {
			log.Fatal(err)
		}
		students = append(students, student)
	}
	defer rows.Close()

	return students
}

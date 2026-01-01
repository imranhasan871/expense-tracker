package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type StudentRepo struct{}

var connStr = "user=admin password=root host=localhost port=5432 dbname=expense_tracker sslmode=disable"

func (s StudentRepo) insertStudent(student Student) error {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	name := student.Name
	email := student.Email
	bloodgroup := student.Bloodgroup
	contact := student.Contact
	address := student.Address

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	bloodgroup = strings.TrimSpace(bloodgroup)
	contact = strings.TrimSpace(contact)
	address = strings.TrimSpace(address)

	insert, err := db.Prepare("INSERT INTO students(name, email, bloodgroup, contact, address) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer insert.Close()

	_, err = insert.Exec(name, email, bloodgroup, contact, address)
	if err != nil {
		return err
	}
	log.Println("Student inserted successfully")
	return nil
}

func (s StudentRepo) getStudents() ([]Student, error) {

	var students []Student

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var student Student
		err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Bloodgroup, &student.Contact, &student.Address)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

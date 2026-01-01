package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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
	fmt.Println("________Studnet Information System________")
	fmt.Println("1. Add a new student")
	fmt.Println("2. View all students")
	fmt.Println("0. Exit")

	var choice int
	fmt.Scanln(&choice)

	if choice == 0 {
		return
	}

	if choice == 1 {
		createStudent()
	}

	if choice == 2 {
		getAllStudents()
	}

}

func getAllStudents() {
	connStr := "user=admin password=root host=localhost port=5432 dbname=expense_tracker sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}

	rows, err := db.Query("SELECT * FROM students")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var student Student
		rows.Scan(&student.id, &student.name, &student.email, &student.bloodgroup, &student.contact, &student.address)
		log.Println(student)
	}

	defer db.Close()
}

func createStudent() {
	connStr := "user=admin password=root host=localhost port=5432 dbname=expense_tracker sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter student name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter student email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter student blood group: ")
	bloodgroup, _ := reader.ReadString('\n')
	bloodgroup = strings.TrimSpace(bloodgroup)

	fmt.Print("Enter student contact: ")
	contact, _ := reader.ReadString('\n')
	contact = strings.TrimSpace(contact)

	fmt.Print("Enter student address: ")
	address, _ := reader.ReadString('\n')
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
	log.Println("Student inserted successfully")
}

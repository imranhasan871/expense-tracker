package main

import (
	"fmt"
)

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
		insertStudent()
	}

	if choice == 2 {
		getStudents()
	}
}

func insertStudent() {

	var name string
	var email string
	var bloodgroup string
	var contact string
	var address string

	fmt.Println("Enter student name: ")
	fmt.Scanln(&name)

	fmt.Println("Enter student email: ")
	fmt.Scanln(&email)

	fmt.Println("Enter student blood group: ")
	fmt.Scanln(&bloodgroup)

	fmt.Println("Enter student contact: ")
	fmt.Scanln(&contact)

	fmt.Println("Enter student address: ")
	fmt.Scanln(&address)

	repo := StudentRepo{}

	student := Student{
		name:       name,
		email:      email,
		bloodgroup: bloodgroup,
		contact:    contact,
		address:    address,
	}

	repo.insertStudent(student)

	fmt.Println("Student added successfully")
}

func getStudents() {
	repo := StudentRepo{}
	students := repo.getStudents()

	for _, student := range students {
		fmt.Println("Name: ", student.name, "Email: ", student.email, "Blood Group: ", student.bloodgroup, "Phone: ", student.contact, "Address: ", student.address)
	}
}

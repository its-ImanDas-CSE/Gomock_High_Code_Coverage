package main

import (
	"fmt"

	"gorm.io/driver/postgres" // Enables communication with a PostgreSQL database, // go get -u gorm.io/driver/postgres
	"gorm.io/gorm"            // GORM library for ORM functionality, // go get -u gorm.io/gorm
)

// Define the Student struct to represent the table schema.
type Student struct {
	ID     int    `gorm:"primaryKey"` // Primary key
	Name   string `gorm:"size:100"`   // String column with max size
	Age    int
	DOB    string
	Course string
	City   string
}

// It tells GORM how to map the Student struct to a specific table name in the database.
// TableName method explicitly specifies the table name for the Student struct.
func (Student) TableName() string {
	return "student" // This function explicitly specifies that the table name in the database is student.
}

//-----------------------------------------------START---------------------------------------------------------------------
// *** This Part of the code is written because we want to create a fake database using gomock.

/*
In the provided code, you are directly using the gorm.DB and Student struct to interact with the database,
and gorm.DB and Student is not an interface. Since mockgen generates mocks for interfaces,
 it will not work directly with gorm.DB and Student because gorm.DB and Student is a struct and not an interface.
*/
// Without using an interface like DBInterface, you cannot mock the database interactions with gomock effectively.
/*
 the function GetStudentNameByID() interacts with the database through the gorm.DB (which has the First() method).

 In your function GetStudentNameByID(), you're calling db.First(&student, id), where db is the GORM database connection.
 First() is a method on the gorm.DB type, which is the real connection to the PostgreSQL database. You're using it to query the database and get the first result that matches the ID.
 If you want to mock the database connection, you need to mock the actual method that's used for querying, which is First() in this case.
*/
type DBInterface interface {
	First(out interface{}, where ...interface{}) *gorm.DB
	// First() is a method that belongs to gorm.DB
	// out represents the variable where the result will be stored and interface{} represent the variable type.
	// interface{} is a special type in Go that can hold any type of value (like a student, a string, or an int).
	// The ...interface{} means that this method can take one or more conditions to filter the data you're looking for.
	// The method returns a pointer to *gorm.DB.
	// *gorm.DB is the type that GORM uses to interact with Database.
}

// -----------------------------------------------*** Ended***------------------------------------------------------------------------

// GetStudentNameByID fetches the student name by ID.
// It takes two parameters: the database connection and the student ID.
func GetStudentNameByID(db DBInterface, id int) (string, error) {
	var student Student              // Variable to hold the fetched student name
	result := db.First(&student, id) // Fetch the first record where ID matches and stored in student variable.
	if result.Error != nil {
		return "", result.Error // Return an error if something goes wrong
	}
	return student.Name, nil // Return the student's name
}

func main() {
	// Define PostgreSQL connection string
	dsn := "host=localhost user=postgres password=Virat@2#Virat@2# dbname=test port=8899 sslmode=disable"

	// Open a connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	fmt.Println("Connected to PostgreSQL!")

	// Fetch and print the name of the student with ID 1
	id := 101                               // Example ID
	name, err := GetStudentNameByID(db, id) // Pass the database connection and the ID
	if err != nil {
		fmt.Printf("Failed to get student name for ID %d: %v\n", id, err)
	} else {
		fmt.Printf("Student name for ID %d: %s\n", id, name)
	}
}

// to create a mock_DBinterface.go, we need to run this command:
//  mockgen -source=C:\Users\iman.das\Desktop\GoLang\Gomock_Testing_Actuall_DB\mainfile_database.go -destination=C:\Users\iman.das\Desktop\GoLang\Gomock_Testing_Actuall_DB\mock_DBinterface.go -package=main

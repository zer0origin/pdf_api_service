package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

var (
	host          = os.Getenv("DATABASE_HOST")
	port, portErr = strconv.Atoi(os.Getenv("DATABASE_PORT"))
	user          = os.Getenv("DATABASE_USER")
	password      = os.Getenv("DATABASE_PASSWORD")
	dbname        = os.Getenv("DATABASE_DB")
)

func getHost() string {
	if host == "" {
		fmt.Println("DATABASE_HOST environment variable not set. Using default: localhost")
		return "localhost"
	}

	return host
}

func getPort() int {
	if portErr != nil {
		fmt.Println("DATABASE_PORT environment variable not set. Using default: 5432")
		return 5432
	}

	return port
}

func getUser() string {
	if user == "" {
		fmt.Println("DATABASE_USER environment variable not set. Using default: user")
		return "user"
	}

	return user
}

func getPassword() string {
	if password == "" {
		fmt.Println("DATABASE_PASSWORD environment variable not set. Using default: password")
		return "password"
	}

	return password
}

func getDatabase() string {
	if dbname == "" {
		fmt.Println("DATABASE_DB environment variable not set. Using default: job_store")
		return "job_store"
	}

	return dbname
}

func getPsqlInfo() string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		getHost(), getPort(), getUser(), getPassword(), getDatabase())

	return psqlInfo
}

type createdCallback func(db *sql.DB) error

func WithConnection(callback createdCallback) error {
	db, err := sql.Open("postgres", getPsqlInfo()) //Create connection string
	if err != nil {
		return err
	}

	err = db.Ping() //open up a connection to the database
	if err != nil {
		return err
	}

	err = callback(db)
	if err != nil {
		return err
	}

	defer func(db *sql.DB) { // Runs once withConnection has finished execution!
		err := db.Close()
		if err != nil {
			panic(err)
		}

		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v (type: %T)\n", r, r)
		}
	}(db)

	return nil
} //Defer runs here.

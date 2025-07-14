package postgres

import (
	"database/sql"
	"fmt"
)

type DatabaseHandler struct {
	DbConfig ConfigForDatabase
}

type createdCallback func(db *sql.DB) error

func (t *DatabaseHandler) WithConnection(callback createdCallback) error {
	str := t.DbConfig.GetPsqlInfo()
	db, err := sql.Open("postgres", str) //Create connection string
	if err != nil {
		return err
	}

	err = db.Ping() //open up a connection to the Database
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
}

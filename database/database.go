package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type ConfigForDatabase struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func (t *ConfigForDatabase) getHost() string {
	if t.Host == "" {
		fmt.Println("DATABASE_HOST environment variable not set. Using default: localhost")
		t.Host = "localhost"
	} else {
		fmt.Println("DATABASE_HOST environment variable set to " + t.Host)
	}

	return t.Host
}

func (t *ConfigForDatabase) getPort() string {
	if t.Port == "" {
		fmt.Println("DATABASE_PORT environment variable not set. Using default: 5432")
		t.Port = "5432"
	} else {
		fmt.Println("DATABASE_PORT environment variable set to " + t.Port)
	}

	return t.Port
}

func (t *ConfigForDatabase) getUser() string {
	if t.Username == "" {
		fmt.Println("DATABASE_USER environment variable not set. Using default: user")
		t.Username = "user"
	} else {
		fmt.Println("DATABASE_USER environment variable set to " + t.Username)
	}

	return t.Username
}

func (t *ConfigForDatabase) getPassword() string {
	if t.Password == "" {
		fmt.Println("DATABASE_PASSWORD environment variable not set. Using default: Password")
		t.Password = "password"
	} else {
		fmt.Println("DATABASE_PASSWORD environment variable set to " + t.Password)
	}

	return t.Password
}

func (t *ConfigForDatabase) getDatabase() string {
	if t.Database == "" {
		fmt.Println("DATABASE_DB environment variable not set. Using default: postgres")
		t.Database = "postgres"
	} else {
		fmt.Println("DATABASE_DB environment variable set to " + t.Database)
	}

	return t.Database
}

func (t *ConfigForDatabase) getPsqlInfo() string {
	psqlInfo := fmt.Sprintf("Host=%s Port=%s user=%s Password=%s dbname=%s sslmode=disable",
		t.getHost(), t.getPort(), t.getUser(), t.getPassword(), t.getDatabase())

	return psqlInfo
}

type createdCallback func(db *sql.DB) error

func (t *ConfigForDatabase) WithConnection(callback createdCallback) error {
	fmt.Println("CALLED WITH CONNECTION")
	str := t.getPsqlInfo()
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
} //Defer runs here.

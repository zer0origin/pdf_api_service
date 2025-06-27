package database

import "os"

var (
	host     = os.Getenv("DATABASE_HOST")
	port     = os.Getenv("DATABASE_PORT")
	user     = os.Getenv("DATABASE_USER")
	password = os.Getenv("DATABASE_PASSWORD")
	dbname   = os.Getenv("DATABASE_DATABASE")
)

func getHost() string {
	if host == "" {
		return "localhost"
	}

	return host
}

func getPort() string {
	if port == "" {
		return "8080"
	}

	return port
}

func getUser() string {
	if user == "" {
		return "postgres"
	}

	return user
}

func getPassword() string {
	if password == "" {
		return "postgres"
	}

	return password
}

func getDatabase() string {
	if dbname == "" {
		return "postgres"
	}

	return dbname
}

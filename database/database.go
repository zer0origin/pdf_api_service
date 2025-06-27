package database

import "os"

var (
	host          = os.Getenv("DATABASE_HOST")
	port, portErr = strconv.Atoi(os.Getenv("DATABASE_PORT"))
	user          = os.Getenv("DATABASE_USER")
	password      = os.Getenv("DATABASE_PASSWORD")
	dbname        = os.Getenv("DATABASE_DATABASE")
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

func getPsqlInfo() string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		getHost(), getPort(), getUser(), getPassword(), getDatabase())

	return psqlInfo
}

type createdCallback func(db *sql.DB)

func withConnection(callback createdCallback) error {
	db, err := sql.Open("postgres", getPsqlInfo())
	if err != nil {
		panic(err)
	}
	callback(db)

	defer func(db *sql.DB) { // Runs once withConnection has finished execution!
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	return err
} //Defer runs here.

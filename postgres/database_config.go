package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
)

type ConfigForDatabase struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	ConUrl   string
}

func (t *ConfigForDatabase) HostOrDefault() string {
	if t.Host == "" {
		fmt.Println("DATABASE_HOST environment variable not set. Using default: localhost")
		t.Host = "localhost"
	} else {
		fmt.Println("DATABASE_HOST environment variable set to " + t.Host)
	}

	return t.Host
}

func (t *ConfigForDatabase) PortOrDefault() string {
	if t.Port == "" {
		fmt.Println("DATABASE_PORT environment variable not set. Using default: 5432")
		t.Port = "5432"
	} else {
		fmt.Println("DATABASE_PORT environment variable set to " + t.Port)
	}

	return t.Port
}

func (t *ConfigForDatabase) UserOrDefault() string {
	if t.Username == "" {
		fmt.Println("DATABASE_USER environment variable not set. Using default: user")
		t.Username = "user"
	} else {
		fmt.Println("DATABASE_USER environment variable set to " + t.Username)
	}

	return t.Username
}

func (t *ConfigForDatabase) PasswordOrDefault() string {
	if t.Password == "" {
		fmt.Println("DATABASE_PASSWORD environment variable not set. Using default: Password")
		t.Password = "password"
	} else {
		fmt.Println("DATABASE_PASSWORD environment variable set to " + t.Password)
	}

	return t.Password
}

func (t *ConfigForDatabase) DatabaseOrDefault() string {
	if t.Database == "" {
		fmt.Println("DATABASE_DB environment variable not set. Using default: postgres")
		t.Database = "postgres"
	} else {
		fmt.Println("DATABASE_DB environment variable set to " + t.Database)
	}

	return t.Database
}

func (t *ConfigForDatabase) GetPsqlInfo() string {
	if t.ConUrl == "" {
		psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			t.GetUsername(), t.GetPassword(), t.GetHost(), t.GetPort(), t.GetDatabase())

		return psqlInfo
	}

	return t.ConUrl
}

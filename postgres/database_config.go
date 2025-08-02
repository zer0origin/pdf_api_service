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

func (t *ConfigForDatabase) GetHost() string {
	if len(t.Host) == 0 {
		panic("Host variable not set for database!")
	}

	return t.Host
}

func (t *ConfigForDatabase) GetPort() string {
	if len(t.Port) == 0 {
		panic("Host variable not set for database!")
	}

	return t.Port
}

func (t *ConfigForDatabase) GetUsername() string {
	if len(t.Username) == 0 {
		panic("Host variable not set for database!")
	}

	return t.Username
}

func (t *ConfigForDatabase) GetPassword() string {
	if len(t.Password) == 0 {
		panic("Host variable not set for database!")
	}

	return t.Password
}

func (t *ConfigForDatabase) GetDatabase() string {
	if len(t.Database) == 0 {
		panic("Host variable not set for database!")
	}

	return t.Database
}

func (t *ConfigForDatabase) GetPsqlInfo() string {
	if len(t.ConUrl) == 0 {
		psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			t.GetUsername(), t.GetPassword(), t.GetHost(), t.GetPort(), t.GetDatabase())

		return psqlInfo
	}

	return t.ConUrl
}

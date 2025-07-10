package unit

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"pdf_service_api/database"
	"testing"
)

func TestDatabaseEmptyArgs(t *testing.T) {
	dbConfig := database.ConfigForDatabase{
		Host:     "",
		Port:     "",
		Username: "",
		Password: "",
	}

	err := dbConfig.WithConnection(func(db *sql.DB) error {
		return nil
	})

	assert.Error(t, err)
	assert.Equal(t, "postgres://user:password@localhost:5432/postgres?sslmode=disable", dbConfig.GetPsqlInfo())
}

func TestDatabaseEmptyCon(t *testing.T) {
	dbConfig := database.ConfigForDatabase{
		ConUrl: "",
	}

	err := dbConfig.WithConnection(func(db *sql.DB) error {
		return nil
	})
	fmt.Println(err)
	assert.Error(t, err)
	assert.Error(t, err)
	assert.Equal(t, "postgres://user:password@localhost:5432/postgres?sslmode=disable", dbConfig.GetPsqlInfo())
}

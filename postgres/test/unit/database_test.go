package unit

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	pg "pdf_service_api/postgres"
	"testing"
)

func TestDatabaseEmptyArgs(t *testing.T) {
	handler := pg.DatabaseHandler{
		DbConfig: pg.ConfigForDatabase{
			Host:     "",
			Port:     "",
			Username: "",
			Password: "",
		}}

	err := handler.WithConnection(func(db *sql.DB) error {
		return nil
	})

	assert.Error(t, err)
	assert.Equal(t, "postgres://user:password@localhost:5432/postgres?sslmode=disable", handler.DbConfig.GetPsqlInfo())
}

func TestDatabaseEmptyCon(t *testing.T) {
	handler := pg.DatabaseHandler{
		DbConfig: pg.ConfigForDatabase{ConUrl: ""},
	}

	err := handler.WithConnection(func(db *sql.DB) error {
		return nil
	})
	fmt.Println(err)
	assert.Error(t, err)
	assert.Error(t, err)
	assert.Equal(t, "postgres://user:password@localhost:5432/postgres?sslmode=disable", handler.DbConfig.GetPsqlInfo())
}

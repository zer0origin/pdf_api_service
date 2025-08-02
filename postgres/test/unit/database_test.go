package unit

import (
	"database/sql"
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

	assert.Panics(t, func() {
		_ = handler.WithConnection(func(db *sql.DB) error {
			return nil
		})
	})
}

func TestDatabaseEmptyCon(t *testing.T) {
	handler := pg.DatabaseHandler{
		DbConfig: pg.ConfigForDatabase{ConUrl: ""},
	}

	assert.Panics(t, func() {
		_ = handler.WithConnection(func(db *sql.DB) error {
			return nil
		})
	})
}

package unit

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"pdf_service_api/service/postgres"
	_ "pdf_service_api/testutil"
	"testing"
)

func TestDatabaseEmptyArgs(t *testing.T) {
	handler := postgres.DatabaseHandler{
		DbConfig: postgres.ConfigForDatabase{
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
	handler := postgres.DatabaseHandler{
		DbConfig: postgres.ConfigForDatabase{ConUrl: ""},
	}

	assert.Panics(t, func() {
		_ = handler.WithConnection(func(db *sql.DB) error {
			return nil
		})
	})
}

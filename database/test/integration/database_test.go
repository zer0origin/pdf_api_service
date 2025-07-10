package integration

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"pdf_service_api/models"
	"pdf_service_api/testutil"
	"testing"
)

var dbUser = "user"
var dbPassword = "password"

func TestGetDatabase(t *testing.T) {
	TestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "TestGetDatabase", dbUser, dbPassword)
	assert.NoError(t, err)
	t.Cleanup(testutil.CleanUp(ctx, *ctr))

	dbConfig, err := testutil.CreateDbConfig(ctx, *ctr)
	assert.NoError(t, err)

	document := &models.Document{}
	err = dbConfig.WithConnection(func(db *sql.DB) error {
		sqlStatement := `SELECT "Document_UUID", "Document_Base64" FROM document_table WHERE "Document_UUID" = $1`
		row := db.QueryRow(sqlStatement, TestUUID)
		err := row.Scan(&document.Uuid, &document.PdfBase64)
		assert.NoError(t, err)

		return nil
	})
	assert.NoError(t, err, "Error running database statement!")
	assert.Equal(t, TestUUID, document.Uuid.String())
}

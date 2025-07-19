package integration

import (
	"pdf_service_api/testutil"
	"testing"
)

func TestMetaIntegration(t *testing.T) {
	t.Run("getMetaDataFromDatabase", getMetaDataFromDatabase)
}

func getMetaDataFromDatabase(t *testing.T) {
	t.Parallel()
	router := testutil.CreateV1RouterAndPostgresContainer(t, "BasicSetupWithOneDocumentTableEntryTwoSelectionsAndMetaData", dbUser, dbPassword)

	GetMetaR
}

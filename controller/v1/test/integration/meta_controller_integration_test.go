package integration

import (
	"testing"
)

func TestMetaIntegration(t *testing.T) {
	t.Run("getMetaDataFromDatabase", getMetaDataFromDatabase)
}

func getMetaDataFromDatabase(t *testing.T) {
	t.Parallel()
	//router := testutil.CreateV1RouterAndPostgresContainer(t, "BasicSetupWithOneDocumentTableEntryTwoSelectionsAndMetaData", dbUser, dbPassword)

}

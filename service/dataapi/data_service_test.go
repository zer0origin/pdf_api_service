package dataapi

import (
	_ "embed"
	"fmt"
	"pdf_service_api/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestSendMetaRequest(t *testing.T) {
	if *testutil.SkipDataApiIntegrationTest {
		t.Skip("Skipping test due to flags")
	}

	p, ctr, err := testutil.CreateDataApiTestContainer()
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	srv := DataService{BaseUrl: fmt.Sprintf("http://localhost:%d", p.Int())}
	meta, err := srv.SendMetaRequest(testutil.HundredPagesPdfInBase64)
	require.NoError(t, err)
	assert.EqualValues(t, 792, *meta.Height)
	assert.EqualValues(t, 612, *meta.Width)
	assert.EqualValues(t, 99, *meta.NumberOfPages)
	assert.EqualValues(t, 99, len(*meta.Images))
}

func TestSendExtractionRequest(t *testing.T) {
	if *testutil.SkipDataApiIntegrationTest {
		t.Skip("Skipping test due to flags")
	}

	p, ctr, err := testutil.CreateDataApiTestContainer()
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	srv := DataService{BaseUrl: fmt.Sprintf("http://localhost:%d", p.Int())}
	meta, err := srv.SendMetaRequest(testutil.HundredPagesPdfInBase64)
	require.NoError(t, err)
	assert.EqualValues(t, 792, *meta.Height)
	assert.EqualValues(t, 612, *meta.Width)
	assert.EqualValues(t, 99, *meta.NumberOfPages)
	assert.EqualValues(t, 99, len(*meta.Images))
}

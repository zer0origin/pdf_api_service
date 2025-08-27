package dataapi

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"pdf_service_api/testutil"
	"testing"
)

//go:embed test-data/100-pages-pdf.txt
var hundredPagesPdf string

func TestSendMetaRequest(t *testing.T) {
	p, ctr, err := testutil.CreateDataApiTestContainer()
	require.NoError(t, err)

	srv := DataService{baseUrl: fmt.Sprintf("http://localhost:%d", p.Int())}
	meta, err := srv.SendMetaRequest(hundredPagesPdf)
	require.NoError(t, err)
	assert.EqualValues(t, 792, *meta.Height)
	assert.EqualValues(t, 612, *meta.Width)
	assert.EqualValues(t, 101, *meta.NumberOfPages)
	assert.EqualValues(t, 101, len(*meta.Images))
	defer testcontainers.TerminateContainer(ctr)
}

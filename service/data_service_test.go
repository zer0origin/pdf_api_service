package service

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

//go:embed test-data/100-pages-pdf.txt
var hundredPagesPdf string

func TestSendMetaRequest(t *testing.T) {
	p, ctr, err := CreateTestContainer()
	defer testcontainers.TerminateContainer(ctr)
	require.NoError(t, err)

	srv := DataService{baseUrl: fmt.Sprintf("http://localhost:%d", p.Int())}
	meta, err := srv.SendMetaRequest(hundredPagesPdf)
	require.NoError(t, err)
	assert.EqualValues(t, 792, *meta.Height)
	assert.EqualValues(t, 612, *meta.Width)
	assert.EqualValues(t, 101, *meta.NumberOfPages)
	assert.EqualValues(t, 101, len(*meta.Images))
}

func CreateTestContainer() (nat.Port, *testcontainers.DockerContainer, error) {
	ctx := context.Background()
	ctr, err := testcontainers.Run(ctx, "pdf_service_data:0.0.5",
		testcontainers.WithWaitStrategy(wait.ForHTTP("/ping").WithPort("8080/tcp")),
		testcontainers.WithExposedPorts("8080"))
	if err != nil {
		return "", ctr, nil
	}

	fmt.Printf("Container started!\n")
	p, err := ctr.MappedPort(ctx, "8080")
	fmt.Printf("Postgres container listening to: %s\n", p)
	return p, ctr, err
}

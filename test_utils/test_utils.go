package test_utils

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/myscribae/myscribae-sdk-go/gql"
	"github.com/myscribae/myscribae-sdk-go/provider"
)

const (
	TestProviderUuid = "00000000-0000-0000-0000-000000000001"
)

func PreTest() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func SetupTest(t *testing.T, tf func(ctx context.Context, prov *provider.Provider) error) {
	PreTest()
	var (
		apiUrl   = os.Getenv("MYSCRIBAE_API_URL")
		apiToken = os.Getenv("MYSCRIBAE_API_TOKEN")
	)

	client := gql.CreateGraphQLClient(
		apiUrl,
		&apiToken,
	)

	prov := provider.Provider{
		ApiUrl: apiUrl,
		Uuid:   uuid.MustParse(TestProviderUuid),
		Client: client,
	}
	ctx := context.Background()
	err := tf(ctx, &prov)
	if err != nil {
		t.Errorf("Test failed: %s", err)
	}
}

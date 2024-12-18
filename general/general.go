package general

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/hasura/go-graphql-client"
	"github.com/myscribae/myscribae-sdk-go/environment"
	"github.com/myscribae/myscribae-sdk-go/gql"
)

type MyScribae struct {
	publicKey *rsa.PublicKey
	client    *graphql.Client
}

func NewMyScribae(client *graphql.Client) *MyScribae {
	return &MyScribae{
		client: client,
	}
}

func (m *MyScribae) Client() *graphql.Client {
	if m.client == nil {
		var url = os.Getenv(environment.ApiUrlEnvVar)
		if url == "" {
			panic("api url not set")
		}

		m.client = graphql.NewClient(
			url,
			&http.Client{},
		)
	}

	return m.client
}

func (m *MyScribae) PublicKey(ctx context.Context) (*rsa.PublicKey, error) {
	if m.publicKey == nil {
		var res gql.GetMyScribaePublicKey
		if err := m.Client().Query(ctx, &res, nil); err != nil {
			log.Printf("failed to get public key: %s", err.Error())
			return nil, errors.New("failed to get public key")
		}

		// Decode the PEM-encoded public key
		block, _ := pem.Decode([]byte(res.PublicKey))
		if block == nil || block.Type != "PUBLIC KEY" {
			return nil, errors.New("failed to decode PEM-encoded public key")
		}

		// Parse the RSA public key
		myscribaePublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, errors.New("failed to parse RSA public key")
		}

		// Ensure the key is of type *rsa.PublicKey
		rsaPublicKey, ok := myscribaePublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("public key is not of type *rsa.PublicKey")
		}

		m.publicKey = rsaPublicKey
	}

	return m.publicKey, nil
}

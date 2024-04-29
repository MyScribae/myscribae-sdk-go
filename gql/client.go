package gql

import (
	"time"

	"net/http"

	"github.com/hasura/go-graphql-client"
)

type ContextKey string

const (
	GraphQLClientKey ContextKey = "gql_client"
)

func CreateGraphQLClient(
	graphqlUrl string,
	apiToken *string,
) *graphql.Client {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	httpClient.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}

	client := graphql.NewClient(
		graphqlUrl,
		httpClient,
	)

	if apiToken != nil {
		client = client.WithRequestModifier(
			func(r *http.Request) {
				r.Header.Set("X-MyScribae-ApiToken", *apiToken)
			},
		)
	}

	return client
}

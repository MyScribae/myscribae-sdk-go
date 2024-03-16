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
) *graphql.Client {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	httpClient.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}

	return graphql.NewClient(
		graphqlUrl,
		httpClient,
	)
}

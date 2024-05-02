package myscribae_sdk

import (
	"context"

	"github.com/myscribae/myscribae-sdk-go/provider"
)

func NewProvider(config provider.ProviderConfig) (*provider.Provider, error) {
	return provider.InitializeProvider(
		context.Background(),
		config,
	)
}
package myscribae_sdk

import (
	"context"

	"github.com/Pritch009/myscribae-sdk-go/provider"
)

func NewProvider(config provider.ProviderConfig) (*provider.Provider, error) {
	return provider.InitializeProvider(
		context.Background(),
		config,
	)
}

type Provider = provider.Provider
type ProviderConfig = provider.ProviderConfig
type ProviderProfileInput = provider.ProviderProfileInput
type ScriptGroup = provider.ScriptGroup
type ScriptGroupInput = provider.ScriptGroupInput
type Script = provider.Script
type ScriptInput = provider.ScriptInput

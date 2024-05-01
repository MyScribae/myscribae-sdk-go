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

type Provider = provider.Provider
type ProviderConfig = provider.ProviderConfig
type ProviderProfileInput = provider.CreateProviderProfileInput
type ScriptGroup = provider.ScriptGroup
type ScriptGroupInput = provider.CreateScriptGroupInput
type Script = provider.Script
type ScriptInput = provider.CreateScriptInput

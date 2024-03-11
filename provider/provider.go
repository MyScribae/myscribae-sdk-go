package provider

import (
	"errors"
	"os"

	"github.com/Pritch009/myscribae-sdk-go/environment"
	"github.com/google/uuid"
)

type Provider struct {
	Uuid        uuid.UUID
	SecretKey   *string
	ApiKey      *string
	ApiUrl      *string
	initialized bool
}

var (
	ErrProviderAlreadyInitialized = errors.New("provider already initialized")
	ErrProviderNotInitialized     = errors.New("provider not initialized")
	ErrMissingApiUrl              = errors.New("missing myscribae api url")
	ErrMissingApiKey              = errors.New("missing myscribae api key")
	ErrMissingSecretKey           = errors.New("missing myscribae secret key")
)

func (p *Provider) connect() error {
	if p.initialized {
		return ErrProviderAlreadyInitialized
	}

	// Attempt to connect to backend services
	

	return nil
}

// / Loads the provider secrets if they are not provided
func (p *Provider) pre_initialize() error {
	if p.initialized {
		return ErrProviderAlreadyInitialized
	}

	if p.ApiUrl == nil {
		apiUrlEnv, success := os.LookupEnv(environment.ApiUrlEnvVar)
		if !success {
			return ErrMissingApiUrl
		}
		p.ApiUrl = &apiUrlEnv
	}

	if p.ApiKey == nil {
		apiKeyEnv, success := os.LookupEnv(environment.ApiKeyEnvVar)
		if !success {
			return ErrMissingApiKey
		}
		p.ApiKey = &apiKeyEnv
	}

	if p.SecretKey == nil {
		secretKeyEnv, success := os.LookupEnv(environment.SecretKeyEnvVar)
		if !success {
			return ErrMissingSecretKey
		}
		p.SecretKey = &secretKeyEnv
	}

	p.initialized = true

	return nil
}

func InitializeProvider(
	apiKey *string,
	secretKey *string,
) (*Provider, error) {
	// load from env
	var (
		_apiKey    *string = apiKey
		_secretKey *string = secretKey
	)

	p := &Provider{
		ApiKey:    _apiKey,
		SecretKey: _secretKey,
	}

	// Attempt to load from env if not provided
	if err := p.pre_initialize(); err != nil {
		return nil, err
	}

	// Attempt to connect to backend service
	if err := p.connect(); err != nil {
		return nil, err
	}

	p.initialized = true

	return p, nil
}

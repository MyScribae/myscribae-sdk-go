package provider

import (
	"context"
	"errors"
	"os"

	"github.com/Pritch009/myscribae-sdk-go/environment"
	"github.com/Pritch009/myscribae-sdk-go/gql"
	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
)

type Provider struct {
	Uuid uuid.UUID
	ProviderProfileInput
	SecretKey     *string
	ApiKey        *string
	ApiUrl        *string
	initialized   bool
	RemoteProfile *gql.ProviderProfile

	client *graphql.Client

	ScriptGroups []*ScriptGroup
}

type ProviderConfig struct {
	ApiKey    string
	SecretKey string
	Url       string
}

type ProviderProfileInput struct {
	AltID       string `json:"alt_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Url         string `json:"url"`
	Color       string `json:"color"`
	Public      bool   `json:"public"`
}

type ScriptGroupInput struct {
	AltID       string        `json:"alt_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Public      bool          `json:"public"`
	Scripts     []ScriptInput `json:"scripts"`
}

type ScriptInput struct {
	AltID            string `json:"alt_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Recurrence       string `json:"recurrence"`
	PriceInCents     int    `json:"price_in_cents"`
	SlaSec           int    `json:"sla_sec"`
	TokenLifetimeSec int    `json:"token_lifetime_sec"`
	Public           bool   `json:"public"`
}

var (
	ErrProviderAlreadyInitialized = errors.New("provider already initialized")
	ErrProviderNotInitialized     = errors.New("provider not initialized")
	ErrMissingApiUrl              = errors.New("missing myscribae api url")
	ErrMissingApiKey              = errors.New("missing myscribae api key")
	ErrMissingSecretKey           = errors.New("missing myscribae secret key")
	ErrFailedToCreateClient       = errors.New("failed to create graphql client")
)

func (p *Provider) Sync() error {
	if p.initialized {
		return ErrProviderAlreadyInitialized
	}

	// Attempt to connect to backend services
	p.client = gql.CreateGraphQLClient(*p.ApiUrl)
	if p.client == nil {
		return ErrFailedToCreateClient
	}

	// Ask client for provider profile
	var query struct {
		ProviderSelf struct {
			Uuid    uuid.UUID
			Profile gql.ProviderProfile `graphql:"... on ProviderProfile"`
		}
	}

	err := p.client.Query(
		context.Background(),
		&query,
		nil,
	)
	if err != nil {
		return err
	}

	p.Uuid = query.ProviderSelf.Uuid

	return nil
}

// / Loads the provider secrets if they are not provided
func (p *Provider) local_initialize(
	scriptGroups []ScriptGroupInput,
) error {
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

	p.ScriptGroups = make([]*ScriptGroup, 0)
	for _, scriptGroup := range scriptGroups {
		scripts := make([]*Script, 0)
		for _, script := range scriptGroup.Scripts {
			scripts = append(scripts, &Script{ScriptInput: script})
		}
		p.ScriptGroups = append(p.ScriptGroups, &ScriptGroup{
			ScriptGroupInput: scriptGroup,
			Scripts:          scripts,
		})
	}

	p.initialized = true

	return nil
}

func InitializeProvider(
	config ProviderConfig,
	providerProfile ProviderProfileInput,
	scriptGroups []ScriptGroupInput,
) (*Provider, error) {
	// load from env
	p := &Provider{
		ApiKey:               &config.ApiKey,
		SecretKey:            &config.SecretKey,
		ApiUrl:               &config.Url,
		ProviderProfileInput: providerProfile,
	}

	// Attempt to load from env if not provided
	if err := p.local_initialize(scriptGroups); err != nil {
		return nil, err
	}

	// Attempt to connect to backend service
	if err := p.Sync(); err != nil {
		return nil, err
	}

	return p, nil
}

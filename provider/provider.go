package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
	"github.com/myscribae/myscribae-sdk-go/environment"
	"github.com/myscribae/myscribae-sdk-go/gql"
)

type Provider struct {
	Uuid      uuid.UUID
	SecretKey string
	ApiKey    string
	ApiUrl    string

	publicKey *string

	Client *graphql.Client
}

type ProviderConfig struct {
	ApiKey    *string
	SecretKey *string
	Url       *string
}

type ProviderProfileInput struct {
	AltID          string  `json:"alt_id"`
	Name           string  `json:"name"`
	Category       *string `json:"category"`
	Description    string  `json:"description"`
	LogoUrl        *string `json:"logo_url"`
	BannerUrl      *string `json:"banner_url"`
	Url            *string `json:"url"`
	Color          *string `json:"color"`
	Public         bool    `json:"public"`
	AccountService bool    `json:"account_service"`
}

var (
	ErrProviderAlreadyInitialized = errors.New("provider already initialized")
	ErrProviderNotInitialized     = errors.New("provider not initialized")
	ErrMissingApiUrl              = errors.New("missing myscribae api url")
	ErrMissingApiKey              = errors.New("missing myscribae api key")
	ErrMissingSecretKey           = errors.New("missing myscribae secret key")
	ErrFailedToCreateClient       = errors.New("failed to create graphql client")
)

// ValidateSubscriberToken validates a subscriber token
func (p *Provider) ValidateSubscriberToken(
	token string,
) (*SubscriberToken, error) {
	publicKey, err := p.GetPublicKey()
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(*publicKey), nil
	})
	if err != nil {
		return nil, err
	}

	return NewSubscriberToken(parsedToken)
}

// IssueSubscriberToken issues a subscriber token
func (p *Provider) IssueSubscriberToken(
	subscriberId string,
) (*string, error) {
	var mutation gql.IssueSubscriberToken
	err := p.secretClient().Mutate(
		context.Background(),
		&mutation,
		map[string]interface{}{
			"subscriberId": subscriberId,
		},
	)
	if err != nil {
		return nil, err
	}

	return &mutation.Provider.Tokens.Issue, nil
}

// Sync syncs the provider with the backend
func (p *Provider) Update(ctx context.Context, profile ProviderProfileInput) (*uuid.UUID, error) {
	var query gql.GetProviderProfile
	// Ask client for provider profile
	err := p.Client.Query(
		context.Background(),
		&query,
		nil,
	)
	if err != nil {
		return nil, err
	}

	p.Uuid = query.ProviderSelf.Uuid

	// update provider
	var mutation gql.EditProviderProfile
	if err := p.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"alt_id":          profile.AltID,
		"category":        profile.Category,
		"name":            profile.Name,
		"description":     profile.Description,
		"logo":            profile.LogoUrl,
		"url":             profile.Url,
		"color":           profile.Color,
		"public":          profile.Public,
		"account_service": profile.AccountService,
	}); err != nil {
		log.Panicf("failed to update provider: %s", err.Error())
	}

	return &mutation.Provider.Edit.Uuid, nil
}

// / Read reads the provider profile
func (p *Provider) Read(ctx context.Context) (*gql.ProviderProfile, error) {
	var query gql.GetProviderProfile
	err := p.Client.Query(
		ctx,
		&query,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &query.ProviderSelf, nil
}

// / SetPublic sets the provider to public or private
func (p *Provider) SetPublic(ctx context.Context, public bool) error {
	var mutation gql.EditProviderProfile
	err := p.Client.Mutate(
		ctx,
		&mutation,
		map[string]interface{}{
			"public": public,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetPublicKey returns the public key of the provider
func (p *Provider) GetPublicKey() (*string, error) {
	if p.publicKey != nil {
		return p.publicKey, nil
	}

	var query gql.GetPublicKey
	err := p.secretClient().Query(
		context.Background(),
		&query,
		nil,
	)
	if err != nil {
		return nil, err
	}

	p.publicKey = &query.ProviderSelf.Keys.PublicKey
	return p.publicKey, nil
}

func (p *ProviderProfileInput) Printf(format string, a ...interface{}) {
	log.Printf(fmt.Sprintf("[%s] %s", p.AltID, format), a...)
}

func (p *ProviderProfileInput) Println(a ...interface{}) {
	log.Println(fmt.Sprintf("[%s]", p.AltID), a)
}

func InitializeProvider(
	ctx context.Context,
	config ProviderConfig,
) (*Provider, error) {
	if config.Url == nil {
		apiUrlEnv, success := os.LookupEnv(environment.ApiUrlEnvVar)
		if !success {
			return nil, ErrMissingApiUrl
		}
		config.Url = &apiUrlEnv
	}

	if config.ApiKey == nil {
		apiKeyEnv, success := os.LookupEnv(environment.ApiKeyEnvVar)
		if !success {
			return nil, ErrMissingApiKey
		}
		config.ApiKey = &apiKeyEnv
	}

	if config.SecretKey == nil {
		secretKeyEnv, success := os.LookupEnv(environment.SecretKeyEnvVar)
		if !success {
			return nil, ErrMissingSecretKey
		}
		config.SecretKey = &secretKeyEnv
	}

	// Attempt to connect to backend services
	client := gql.CreateGraphQLClient(*config.Url)
	if client == nil {
		return nil, ErrFailedToCreateClient
	}

	return &Provider{
		ApiKey:    *config.ApiKey,
		SecretKey: *config.SecretKey,
		ApiUrl:    *config.Url,
		Client:    client,
	}, nil
}

// secretClient returns a client with the provider's secret key
func (p *Provider) secretClient() *graphql.Client {
	return p.Client.WithRequestModifier(
		func(r *http.Request) {
			r.Header.Set("X-MyScribae-ApiKey", p.ApiKey)
		},
	)

}

func (p *Provider) ScriptGroup(alt_id string) *ScriptGroup {
	return &ScriptGroup{
		AltID:    alt_id,
		Provider: p,
	}
}

func (p *Provider) Script(script_group_uuid string, script_alt_id string) *Script {
	return &Script{}
}

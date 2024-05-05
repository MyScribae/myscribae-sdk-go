package provider

import (
	"context"
	"encoding/json"
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
	"github.com/myscribae/myscribae-sdk-go/utilities"
)

type Provider struct {
	ApiUrl string

	Uuid      uuid.UUID
	altId     *utilities.AltUuid
	SecretKey *string
	ApiKey    *string

	publicKey *string

	Client *graphql.Client
}

func (p *Provider) ID() utilities.AltUuid {
	if p.altId == nil {
		id, err := utilities.NewAltUuid(p.Uuid.String())
		if err != nil {
			log.Panicf("failed to create alt id: %s", err.Error())
		}
		p.altId = &id
	}

	return *p.altId
}

type ProviderConfig struct {
	ApiKey    *string
	SecretKey *string
	ApiToken  *string
	ApiUrl    *string
}

type CreateProviderProfileInput struct {
	AltID          *string `json:"alt_id"`
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

type UpdateProviderProfileInput struct {
	AltID          *string `json:"alt_id"`
	Name           *string `json:"name"`
	Category       *string `json:"category"`
	Description    *string `json:"description"`
	LogoUrl        *string `json:"logo_url"`
	BannerUrl      *string `json:"banner_url"`
	Url            *string `json:"url"`
	Color          *string `json:"color"`
	Public         *bool   `json:"public"`
	AccountService *bool   `json:"account_service"`
}

var (
	ErrProviderAlreadyInitialized = errors.New("provider already initialized")
	ErrProviderNotInitialized     = errors.New("provider not initialized")
	ErrMissingApiUrl              = errors.New("missing myscribae api url")
	ErrMissingApiKey              = errors.New("missing myscribae api key")
	ErrMissingSecretKey           = errors.New("missing myscribae secret key")
	ErrFailedToCreateClient       = errors.New("failed to create graphql client")
)

func CreateNewProvider(ctx context.Context, client *graphql.Client, input *CreateProviderProfileInput) (*Provider, error) {
	if client == nil {
		return nil, ErrFailedToCreateClient
	}

	var mutation gql.CreateNewProvider
	err := client.Mutate(
		ctx,
		&mutation,
		map[string]interface{}{
			"alt_id":          input.AltID,
			"category":        input.Category,
			"name":            input.Name,
			"description":     input.Description,
			"logo":            input.LogoUrl,
			"url":             input.Url,
			"color":           input.Color,
			"public":          input.Public,
			"account_service": input.AccountService,
		},
	)
	if err != nil {
		return nil, err
	}

	// Created provider, return provider
	prov := &Provider{
		Uuid:   mutation.Providers.Create.Uuid,
		Client: client,
	}

	// get secret key and api key
	err = prov.ResetProviderKeys(ctx)
	if err != nil {
		return nil, err
	}

	return prov, nil
}

// ValidateSubscriberToken validates a subscriber token
func (p *Provider) ValidateSubscriberToken(
	ctx context.Context,
	token string,
) (*SubscriberToken, error) {
	publicKey, err := p.GetPublicKey(ctx)
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
func (p *Provider) Update(ctx context.Context, profile UpdateProviderProfileInput) (*uuid.UUID, error) {
	var changes []byte
	changes, err := profile.MarshalJSON()
	if err != nil {
		return nil, err
	}

	// update provider
	var mutation gql.EditProviderProfile
	if err := p.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"id":      p.ID(),
		"changes": string(changes),
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
		map[string]interface{}{
			"id": p.ID(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &query.ProviderSelf, nil
}

// / SetPublic sets the provider to public or private
func (p *Provider) SetPublic(ctx context.Context, public bool) error {
	_, err := p.Update(ctx, UpdateProviderProfileInput{
		Public: &public,
	})

	return err
}

// GetPublicKey returns the public key of the provider
func (p *Provider) GetPublicKey(ctx context.Context) (*string, error) {
	if p.publicKey != nil {
		return p.publicKey, nil
	}

	var query gql.GetPublicKey
	err := p.secretClient().Query(
		ctx,
		&query,
		map[string]interface{}{
			"provider_id": p.Uuid.ID(),
		},
	)
	if err != nil {
		return nil, err
	}

	p.publicKey = &query.ProviderSelf.Keys.PublicKey
	return p.publicKey, nil
}

func (p *CreateProviderProfileInput) Printf(format string, a ...interface{}) {
	log.Printf(fmt.Sprintf("[%v] %s", p.AltID, format), a...)
}

func (p *CreateProviderProfileInput) Println(a ...interface{}) {
	log.Println(fmt.Sprintf("[%v]", p.AltID), a)
}

func InitializeProvider(
	ctx context.Context,
	config ProviderConfig,
) (*Provider, error) {
	if config.ApiUrl == nil {
		apiUrlEnv, success := os.LookupEnv(environment.ApiUrlEnvVar)
		if !success {
			return nil, ErrMissingApiUrl
		}
		config.ApiUrl = &apiUrlEnv
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
	client := gql.CreateGraphQLClient(*config.ApiUrl, nil)
	if client == nil {
		return nil, ErrFailedToCreateClient
	}

	return &Provider{
		ApiKey:    config.ApiKey,
		SecretKey: config.SecretKey,
		ApiUrl:    *config.ApiUrl,
		Client: client.WithRequestModifier(
			func(r *http.Request) {
				if config.ApiKey != nil {
					r.Header.Set("X-MyScribae-ApiKey", *config.ApiKey)
				}
				if config.ApiToken != nil {
					r.Header.Set("X-MyScribae-ApiToken", *config.ApiToken)
				}
			},
		),
	}, nil
}

// secretClient returns a client with the provider's secret key
func (p *Provider) secretClient() *graphql.Client {
	client := p.Client
	if p.SecretKey != nil {
		client = p.Client.WithRequestModifier(
			func(r *http.Request) {
				r.Header.Set("X-MyScribae-SecretKey", *p.SecretKey)
			},
		)
	}

	return client
}

func (p *Provider) ScriptGroup(alt_id string) (*ScriptGroup, error) {
	id, err := utilities.NewAltUuid(alt_id)
	if err != nil {
		return nil, err
	}
	return &ScriptGroup{
		AltID:    id,
		Provider: p,
	}, nil
}

func (p *Provider) Script(script_group_id utilities.AltUuid, alt_id string) (*Script, error) {
	id, err := utilities.NewAltUuid(alt_id)
	if err != nil {
		return nil, err
	}
	return &Script{
		AltID:         id,
		ScriptGroupID: script_group_id,
		Provider:      p,
	}, nil
}

func (p *Provider) ResetProviderKeys(ctx context.Context) error {
	var mutation gql.ResetProviderKeys
	err := p.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"provider_id": p.ID(),
	})

	if err != nil {
		return err
	}

	p.ApiKey = &mutation.Provider.Keys.Reset.ApiKey
	p.SecretKey = &mutation.Provider.Keys.Reset.SecretKey

	return nil
}

func (pi *UpdateProviderProfileInput) MarshalJSON() ([]byte, error) {
	// only marshal provided fields, ignore nil fields
	data := make(map[string]interface{})
	if pi.AltID != nil {
		data["alt_id"] = *pi.AltID
	}

	if pi.Name != nil {
		data["name"] = *pi.Name
	}

	if pi.Category != nil {
		data["category"] = *pi.Category
	}

	if pi.Description != nil {
		data["description"] = *pi.Description
	}

	if pi.LogoUrl != nil {
		data["logo"] = *pi.LogoUrl
	}

	if pi.BannerUrl != nil {
		data["banner_url"] = *pi.BannerUrl
	}

	if pi.Url != nil {
		data["url"] = *pi.Url
	}

	if pi.Color != nil {
		data["color"] = *pi.Color
	}

	if pi.Public != nil {
		data["public"] = *pi.Public
	}

	if pi.AccountService != nil {
		data["account_service"] = *pi.AccountService
	}

	return json.Marshal(data)
}

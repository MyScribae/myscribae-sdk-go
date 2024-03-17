package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Pritch009/myscribae-sdk-go/environment"
	"github.com/Pritch009/myscribae-sdk-go/gql"
	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
)

type Provider struct {
	Uuid          uuid.UUID
	Profile       ProviderProfileInput
	SecretKey     string
	ApiKey        string
	ApiUrl        string
	initialized   bool
	RemoteProfile *gql.GetProviderProfile

	Client *graphql.Client

	ScriptGroups []*ScriptGroup
}

type ProviderConfig struct {
	ApiKey    *string
	SecretKey *string
	Url       *string
}

type ProviderProfileInput struct {
	AltID        string             `json:"alt_id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Logo         string             `json:"logo"`
	Url          string             `json:"url"`
	Color        string             `json:"color"`
	Public       bool               `json:"public"`
	ScriptGroups []ScriptGroupInput `json:"script_groups"`
}

var (
	ErrProviderAlreadyInitialized = errors.New("provider already initialized")
	ErrProviderNotInitialized     = errors.New("provider not initialized")
	ErrMissingApiUrl              = errors.New("missing myscribae api url")
	ErrMissingApiKey              = errors.New("missing myscribae api key")
	ErrMissingSecretKey           = errors.New("missing myscribae secret key")
	ErrFailedToCreateClient       = errors.New("failed to create graphql client")
)

func (p *Provider) Sync(ctx context.Context) error {
	p.Println("Syncing provider")

	if p.initialized {
		p.Println("Provider already initialized")
		return ErrProviderAlreadyInitialized
	}

	// Attempt to connect to backend services
	p.Client = gql.CreateGraphQLClient(p.ApiUrl)
	if p.Client == nil {
		return ErrFailedToCreateClient
	}

	var query gql.GetProviderProfile
	// Ask client for provider profile
	err := p.Client.Query(
		context.Background(),
		&query,
		nil,
	)
	if err != nil {
		return err
	}

	p.RemoteProfile = &query
	p.Uuid = query.ProviderSelf.Uuid

	// Sync script groups
	existingScriptGroups := make(map[string]*gql.RemoteScriptGroup)
	for _, remoteScriptGroup := range query.ProviderSelf.ScriptGroups {
		existingScriptGroups[remoteScriptGroup.AltID] = &remoteScriptGroup
	}

	scriptGroupsSeen := make(map[string]bool)
	for _, sg := range p.ScriptGroups {
		if _, ok := scriptGroupsSeen[sg.Profile.AltID]; ok {
			sg.Println("Script group exists and has already been seen")
			continue
		} else {
			sg.Println("Syncing script group")
		}

		sg.Provider = p
		scriptGroupsSeen[sg.Profile.AltID] = true

		sg.Sync(ctx, existingScriptGroups[sg.Profile.AltID])
	}

	for _, rsg := range query.ProviderSelf.ScriptGroups {
		if _, ok := scriptGroupsSeen[rsg.AltID]; !ok {
			p.Printf("Script group %s has been abandoned", rsg.AltID)

			var mutation gql.EditScriptGroup
			err := p.Client.Mutate(ctx, &mutation, map[string]interface{}{
				"id":     rsg.Uuid,
				"public": false,
			})

			if err != nil {
				panic("failed to abandon script group: " + err.Error())
			}
		}
	}

	p.Println("Provider synced")

	return nil
}

func (p *Provider) Printf(format string, a ...interface{}) {
	log.Printf(fmt.Sprintf("[%s] %s", p.Profile.AltID, format), a...)
}

func (p *Provider) Println(a ...interface{}) {
	log.Println(fmt.Sprintf("[%s]", p.Profile.AltID), a)
}

func (p *Provider) SecretClient() *graphql.Client {
	return p.Client.WithRequestModifier(
		func(r *http.Request) {
			r.Header.Set("X-MyScribae-ApiKey", p.ApiKey)
		},
	)

}

// / Loads the provider secrets if they are not provided
func (p *Provider) local_initialize() error {
	p.Println("Initializing provider")

	if p.initialized {
		p.Println("Provider already initialized")
		return ErrProviderAlreadyInitialized
	}

	p.ScriptGroups = make([]*ScriptGroup, 0)
	for _, scriptGroupInput := range p.Profile.ScriptGroups {
		scriptGroup := &ScriptGroup{
			Provider: p,
			Profile:  scriptGroupInput,
			Scripts:  make([]Script, 0),
		}

		for _, script := range scriptGroupInput.Scripts {
			scriptGroup.Scripts = append(scriptGroup.Scripts, Script{
				ScriptGroup: scriptGroup,
				Profile:     script,
			})
		}

		p.ScriptGroups = append(p.ScriptGroups, scriptGroup)
	}

	p.initialized = true
	p.Println("Provider initialized")

	return nil
}

func InitializeProvider(
	ctx context.Context,
	config ProviderConfig,
	providerProfile ProviderProfileInput,
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

	// load from env
	p := &Provider{
		ApiKey:    *config.ApiKey,
		SecretKey: *config.SecretKey,
		ApiUrl:    *config.Url,
		Profile:   providerProfile,
	}

	// Attempt to load from env if not provided
	if err := p.local_initialize(); err != nil {
		return nil, err
	}

	// Attempt to connect to backend service
	if err := p.Sync(ctx); err != nil {
		return nil, err
	}

	return p, nil
}

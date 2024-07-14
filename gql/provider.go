package gql

import (
	"github.com/google/uuid"
)

type ProviderProfile struct {
	Uuid           uuid.UUID `graphql:"uuid"`
	AltID          *string   `graphql:"alt_id"`
	Category       string    `graphql:"category_id"`
	Name           string    `graphql:"name"`
	Description    string    `graphql:"description"`
	Color          *string   `graphql:"color"`
	LogoUrl        *string   `graphql:"logo_url"`
	BannerUrl      *string   `graphql:"banner_url"`
	MyRole         *string   `graphql:"my_role"`
	Url            *string   `graphql:"url"`
	AccountService struct {
		Enabled bool `graphql:"enabled"`
	} `graphql:"account_service"`
	Public bool `graphql:"public"`
}

type RemoteScript struct {
	Uuid    uuid.UUID
	AltID   string
	Version string
}

type RemoteScriptGroup struct {
	Uuid    uuid.UUID
	AltID   string
	Version string
	Scripts []RemoteScript
}

type GetProviderProfile struct {
	ProviderSelf ProviderProfile `graphql:"provider_self(id:$id)"`
}

type EditProviderProfile struct {
	Provider struct {
		Edit struct {
			Uuid uuid.UUID
		} `graphql:"edit(changes:$changes)"`
	} `graphql:"provider(id:$id)"`
}

type GetScriptGroup struct {
	ProviderSelf struct {
		ScriptGroup ScriptGroupProfile `graphql:"script_group(id:$id)"`
	} `graphql:"provider_self(id:$provider_id)"`
}

type ScriptGroupProfile struct {
	Uuid        uuid.UUID `graphql:"uuid"`
	AltID       string    `graphql:"alt_id"`
	Name        string    `graphql:"name"`
	Description string    `graphql:"description"`
	Public      bool      `graphql:"public"`
}

type EditScriptGroup struct {
	Provider struct {
		ScriptGroup struct {
			Edit struct {
				Uuid uuid.UUID
			} `graphql:"edit(changes:$changes)"`
		} `graphql:"script_group(id:$id)"`
	} `graphql:"provider(id:$provider_id)"`
}

type GetScript struct {
	ProviderSelf struct {
		ScriptGroup struct {
			Script GQLScriptProfile `graphql:"script(id:$id)"`
		} `graphql:"script_group(id:$script_group_id)"`
	} `graphql:"provider_self(id:$provider_id)"`
}

type GQLScriptProfile struct {
	Uuid             uuid.UUID `graphql:"uuid"`
	AltID            string    `graphql:"alt_id"`
	Name             string    `graphql:"name"`
	Description      string    `graphql:"description"`
	Recurrence       string    `graphql:"recurrence"`
	PriceInCents     uint      `graphql:"price_in_cents"`
	SlaSec           int       `graphql:"sla_sec"`
	TokenLifetimeSec int       `graphql:"token_lifetime_sec"`
	Public           bool      `graphql:"public"`
}
type ScriptProfile struct {
	Uuid             uuid.UUID `graphql:"uuid"`
	AltID            string    `graphql:"alt_id"`
	Name             string    `graphql:"name"`
	Description      string    `graphql:"description"`
	Recurrence       string    `graphql:"recurrence"`
	PriceInCents     uint      `graphql:"price_in_cents"`
	SlaSec           uint      `graphql:"sla_sec"`
	TokenLifetimeSec uint      `graphql:"token_lifetime_sec"`
	Public           bool      `graphql:"public"`
}

type EditScript struct {
	Provider struct {
		ScriptGroup struct {
			Script struct {
				Edit struct {
					Uuid uuid.UUID
				} `graphql:"edit(changes:$changes)"`
			} `graphql:"script(id:$id)"`
		} `graphql:"script_group(id:$script_group_id)"`
	} `graphql:"provider(id:$provider_id)"`
}

type IssueSubscriberToken struct {
	Provider struct {
		Tokens struct {
			Issue string `graphql:"issue(subscriber_id:$subscriber_id)"`
		}
	}
}

type RequestUserAssociation struct {
	Provider struct {
		Associate uuid.UUID `graphql:"associate(user_identifier:$user_identifier, user_avatar_url: $user_avatar_url, script_credits: $script_credits, redirect: $redirect)"`
	}
}

type GetPublicKey struct {
	ProviderSelf struct {
		Keys struct {
			PublicKey string
		}
	} `graphql:"provider_self(id:$provider_id)"`
}

type CreateNewScript struct {
	Provider struct {
		ScriptGroup struct {
			Scripts struct {
				Create struct {
					Uuid uuid.UUID
				} `graphql:"create(alt_id: $alt_id, name: $name, description: $description, price_in_cents: $price_in_cents, recurrence: $recurrence, sla_sec: $sla_sec, token_lifetime_sec: $token_lifetime_sec, public: $public)"`
			}
		} `graphql:"script_group(id: $script_group_id)"`
	} `graphql:"provider(id: $provider_id)"`
}

type CreateNewScriptGroup struct {
	Provider struct {
		ScriptGroups struct {
			Create struct {
				Uuid uuid.UUID
			} `graphql:"create(alt_id: $alt_id, name: $name, description: $description, public: $public)"`
		} `graphql:"script_groups"`
	} `graphql:"provider(id: $provider_id)"`
}

type ResetProviderKeys struct {
	Provider struct {
		Keys struct {
			Reset struct {
				ApiKey    string `graphql:"api_key"`
				SecretKey string `graphql:"secret_key"`
			} `graphql:"reset"`
		} `graphql:"keys"`
	} `graphql:"provider(id:$provider_id)"`
}

type CreateNewProvider struct {
	Providers struct {
		Create struct {
			Uuid uuid.UUID `graphql:"uuid"`
		} `graphql:"create(name: $name, description: $description, category_id: $category_id)"`
	} `graphql:"providers"`
}

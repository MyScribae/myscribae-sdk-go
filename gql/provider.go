package gql

import (
	"github.com/google/uuid"
)

type ProviderProfile struct {
	Uuid           uuid.UUID
	AltID          *string
	Category       string
	Name           string
	Description    string
	AccountTier    *string
	Color          *string
	LogoUrl        *string
	BannerUrl      *string
	MyRole         *string
	Url            *string
	AccountService bool
	Public         bool
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
	Provider ProviderProfile `graphql:"provider(id: $provider_id)"`
}

type EditProviderProfile struct {
	Provider struct {
		Edit struct {
			Uuid uuid.UUID
		} `graphql:"edit(alt_id: $alt_id, name: $name, category: $category, description: $description, logo: $logo, url: $url, color: $color, public: $public)"`
	} `graphql:"provider(id: $provider_id)"`
}

type GetScriptGroup struct {
	Provider struct {
		ScriptGroup ScriptGroupProfile `graphql:"script_group(id: $id)"`
	} `graphql:"provider(id: $provider_id)"`
}

type ScriptGroupProfile struct {
	Uuid        uuid.UUID
	AltID       string
	Name        string
	Description string
	Public      bool
}

type EditScriptGroup struct {
	Provider struct {
		ScriptGroup struct {
			Edit struct {
				Uuid uuid.UUID
			} `graphql:"edit(name: $name, description: $description, public: $public)"`
		} `graphql:"script_group(id: $id)"`
	} `graphql:"provider(id: $provider_id)"`
}

type GetScript struct {
	Provider struct {
		ScriptGroup struct {
			Script ScriptProfile `graphql:"script(id: $id)"`
		} `graphql:"script_group(id: $script_group_id)"`
	} `graphql:"provider(id: $provider_id)"`
}

type ScriptProfile struct {
	Uuid             uuid.UUID
	AltID            string
	Name             string
	Description      string
	Recurrence       string
	PriceInCents     int
	SlaSec           int
	TokenLifetimeSec int
	Public           bool
}

type EditScript struct {
	Provider struct {
		ScriptGroup struct {
			Script struct {
				Edit struct {
					Uuid uuid.UUID
				} `graphql:"edit(name: $name, description: $description, price_in_cents: $price_in_cents, sla_sec: $sla_sec, token_lifetime_sec: $token_lifetime_sec, public: $public)"`
			} `graphql:"script(id: $id)"`
		} `graphql:"script_group(id: $script_group_uuid)"`
	} `graphql:"provider(id: $provider_id)"`
}

type IssueSubscriberToken struct {
	Provider struct {
		Tokens struct {
			Issue string `graphql:"issue(subscriber_id: $subscriber_id)"`
		}
	}
}

type RequestUserAssociation struct {
	Provider struct {
		Associate uuid.UUID `graphql:"associate(user_identifier: $user_identifier, user_avatar_url: $user_avatar_url, script_credits: $script_credits, redirect: $redirect)"`
	}
}

type GetPublicKey struct {
	ProviderSelf struct {
		Keys struct {
			PublicKey string
		}
	}
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
		}
	} `graphql:"provider(id: $provider_id)"`
}

type ResetProviderKeys struct {
	ProviderSelf struct {
		Keys struct {
			Reset struct {
				ApiKey    string `graphql:"api_key"`
				SecretKey string `graphql:"secret_key"`
			} `graphql:"reset"`
		} `graphql:"keys"`
	} `graphql:"provider"`
}

type CreateNewProvider struct {
	Providers struct {
		Create struct {
			Uuid uuid.UUID `graphql:"uuid"`
		} `graphql:"create(name: $name, description: $description, category: $category)"`
	} `graphql:"providers"`
}

package gql

import (
	"github.com/google/uuid"
)

type ProviderProfile struct {
	AccountTier string
	AltId       *string
	Color       *string
	Description string
	Logo        *string
	MyRole      *string
	Name        string
	Url         *string
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
	ProviderSelf struct {
		Uuid         uuid.UUID
		Version      string
		ScriptGroups []RemoteScriptGroup
	}
}

type EditProviderProfile struct {
	Provider struct {
		Edit struct {
			Uuid uuid.UUID
		} `graphql:"edit(alt_id: $alt_id, name: $name, description: $description, logo: $logo, url: $url, color: $color, public: $public)"`
	}
}

type EditScriptGroup struct {
	Provider struct {
		ScriptGroup struct {
			Edit struct {
				Uuid uuid.UUID
			} `graphql:"edit(name: $name, description: $description, public: $public)"`
		} `graphql:"script_group(id: $id)"`
	}
}

type EditScript struct {
	Provider struct {
		ScriptGroup struct {
			Script struct {
				Edit struct {
					Uuid uuid.UUID
				} `graphql:"edit(name: $name, description: $description, price_in_cents: $price_in_cents, sla_sec: $sla_sec, token_lifetime_sec: $token_lifetime_sec, public: $public)"`
			} `graphql:"script(id: $id)"`
		} `graphql:"script_group(id: $script_group_id)"`
	}
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
		Associate uuid.UUID `graphql:"associate(user_identifier: $user_identifier, user_avatar_src: $user_avatar_src, script_credits: $script_credits, redirect: $redirect)"`
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
		} `graphql:"script_group(id: $scriptGroupId)"`
	} `graphql:"provider(id: $providerId)"`
}

type CreateNewScriptGroup struct {
	Provider struct {
		ScriptGroups struct {
			Create struct {
				Uuid uuid.UUID
			} `graphql:"create(alt_id: $alt_id, name: $name, description: $description, public: $public)"`
		}
	} `graphql:"provider(id: $providerId)"`
}

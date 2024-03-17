package provider

import "github.com/google/uuid"

type ScriptClaim struct {
	SubscriptionUuid uuid.UUID `json:"subscription_uuid"`
	ScriptGroupUuid  uuid.UUID `json:"script_group_uuid"`
	ScriptGroupAltID string    `json:"script_group_alt_id"`
	ScriptUuid       uuid.UUID `json:"script_uuid"`
	ScriptAltID      string    `json:"script_alt_id"`
}

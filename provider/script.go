package provider

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/myscribae/myscribae-sdk-go/gql"
	"github.com/myscribae/myscribae-sdk-go/utilities"
)

type Script struct {
	ScriptGroupID utilities.AltUUID
	AltID         utilities.AltUUID
	Uuid          *uuid.UUID
	Provider      *Provider
}

type CreateScriptInput struct {
	AltID            string               `json:"alt_id"`
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	Recurrence       utilities.Recurrence `json:"recurrence"`
	PriceInCents     utilities.MoneyValue `json:"price_in_cents"`
	SlaSec           utilities.UInt       `json:"sla_sec"`
	TokenLifetimeSec utilities.UInt       `json:"token_lifetime_sec"`
	Public           bool                 `json:"public"`
}

type UpdateScriptInput struct {
	Name             *string               `json:"name"`
	Description      *string               `json:"description"`
	PriceInCents     *utilities.MoneyValue `json:"price_in_cents"`
	SlaSec           *utilities.UInt       `json:"sla_sec"`
	TokenLifetimeSec *utilities.UInt       `json:"token_lifetime_sec"`
	Public           *bool                 `json:"public"`
}

func (s *Script) Create(ctx context.Context, input CreateScriptInput) (*uuid.UUID, error) {
	var mutation gql.CreateNewScript
	err := s.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"provider_id":        s.Provider.ID(),
		"script_group_id":    s.ScriptGroupID,
		"alt_id":             input.AltID,
		"name":               input.Name,
		"description":        input.Description,
		"recurrence":         input.Recurrence,
		"price_in_cents":     input.PriceInCents,
		"sla_sec":            input.SlaSec,
		"token_lifetime_sec": input.TokenLifetimeSec,
		"public":             input.Public,
	})

	if err != nil {
		return nil, err
	}

	s.Uuid = &mutation.Provider.ScriptGroup.Scripts.Create.Uuid
	return s.Uuid, nil
}

func (s *Script) Read(ctx context.Context) (*gql.ScriptProfile, error) {
	var query gql.GetScript
	if err := s.Provider.Client.Query(ctx, &query, map[string]interface{}{
		"provider_id":     s.Provider.ID(),
		"script_group_id": s.ScriptGroupID,
		"id":              s.AltID,
	}); err != nil {
		return nil, err
	}

	s.Uuid = &query.ProviderSelf.ScriptGroup.Script.Uuid
	return &gql.ScriptProfile{
		Uuid:             query.ProviderSelf.ScriptGroup.Script.Uuid,
		AltID:            query.ProviderSelf.ScriptGroup.Script.AltID,
		Name:             query.ProviderSelf.ScriptGroup.Script.Name,
		Description:      query.ProviderSelf.ScriptGroup.Script.Description,
		Recurrence:       query.ProviderSelf.ScriptGroup.Script.Recurrence,
		PriceInCents:     query.ProviderSelf.ScriptGroup.Script.PriceInCents,
		SlaSec:           uint(query.ProviderSelf.ScriptGroup.Script.SlaSec),
		TokenLifetimeSec: uint(query.ProviderSelf.ScriptGroup.Script.TokenLifetimeSec),
		Public:           query.ProviderSelf.ScriptGroup.Script.Public,
	}, nil
}

func (s *Script) Update(ctx context.Context, input UpdateScriptInput) (*uuid.UUID, error) {
	changes, err := input.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var mutation gql.EditScript
	err = s.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"provider_id":     s.Provider.ID(),
		"script_group_id": s.ScriptGroupID,
		"id":              s.AltID,
		"changes":         string(changes),
	})
	if err != nil {
		return nil, err
	}

	s.Uuid = &mutation.Provider.ScriptGroup.Script.Edit.Uuid
	return s.Uuid, nil
}

func (s *Script) Delete(ctx context.Context) error {
	var changes = struct {
		Public bool `json:"public"`
	} {
		Public: false,
	}
	changesBytes, err := json.Marshal(changes)
	if err != nil {
		return err
	}

	var mutation gql.EditScript
	err = s.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"provider_id":     s.Provider.ID(),
		"script_group_id": s.ScriptGroupID,
		"id":              s.AltID,
		"changes": 		  string(changesBytes),
	})
	if err != nil {
		return err
	}

	return nil
}

func (si *UpdateScriptInput) MarshalJSON() ([]byte, error) {
	// marshal as a map, not including optional fields if they are not set
	m := map[string]interface{}{}

	if si.Name != nil {
		m["name"] = *si.Name
	}

	if si.Description != nil {
		m["description"] = *si.Description
	}

	if si.PriceInCents != nil {
		m["price_in_cents"] = *si.PriceInCents
	}

	if si.SlaSec != nil {
		m["sla_secs"] = *si.SlaSec
	}

	if si.TokenLifetimeSec != nil {
		m["token_lifetime_secs"] = *si.TokenLifetimeSec
	}

	if si.Public != nil {
		m["public"] = *si.Public
	}

	return json.Marshal(m)
}

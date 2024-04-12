package provider

import (
	"context"

	"github.com/MyScribae/myscribae-sdk-go/gql"
	"github.com/google/uuid"
)

type Script struct {
	ScriptGroupUuid *uuid.UUID
	Uuid            *uuid.UUID
	AltID           string
	Provider        *Provider
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

func (s *Script) Create(ctx context.Context, input ScriptInput) (*uuid.UUID, error) {
	var mutation gql.CreateNewScript
	err := s.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"script_group_id":    s.ScriptGroupUuid,
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
		"script_group_id": s.ScriptGroupUuid,
		"id":              s.AltID,
	}); err != nil {
		return nil, err
	}

	s.Uuid = &query.Provider.ScriptGroup.Script.Uuid
	return &query.Provider.ScriptGroup.Script, nil
}

func (s *Script) Update(ctx context.Context, input ScriptInput) (*uuid.UUID, error) {
	var mutation gql.EditScript
	err := s.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"script_group_id":    s.ScriptGroupUuid,
		"id":                 s.AltID,
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
		panic("failed to update script: " + err.Error())
	}

	s.Uuid = &mutation.Provider.ScriptGroup.Script.Edit.Uuid
	return s.Uuid, nil
}

func (s *Script) Delete(ctx context.Context) error {
	var mutation gql.EditScript
	err := s.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"script_group_id": s.ScriptGroupUuid,
		"id":              s.AltID,
		"public":          false,
	})
	if err != nil {
		return err
	}

	return nil
}

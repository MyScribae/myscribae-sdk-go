package provider

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/myscribae/myscribae-sdk-go/gql"
	"github.com/myscribae/myscribae-sdk-go/utilities"
)

type ScriptGroup struct {
	AltID    utilities.AltUuid
	Uuid     *uuid.UUID
	Provider *Provider
}

type CreateScriptGroupInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type UpdateScriptGroupInput struct {
	AltID       *string `json:"alt_id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Public      *bool   `json:"public"`
}

func (sg *ScriptGroup) Update(ctx context.Context, profile UpdateScriptGroupInput) (*uuid.UUID, error) {
	var mutation gql.EditScriptGroup
	changes, err := profile.MarshalJSON()
	if err != nil {
		return nil, err
	}

	err = sg.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"provider_id": sg.Provider.ID(),
		"id":          sg.AltID,
		"changes":     string(changes),
	})
	if err != nil {
		panic("failed to update script group: " + err.Error())
	}

	sg.Uuid = &mutation.Provider.ScriptGroup.Edit.Uuid
	return sg.Uuid, nil
}

func (sg *ScriptGroup) Read(ctx context.Context) (*gql.ScriptGroupProfile, error) {
	var query gql.GetScriptGroup
	if err := sg.Provider.Client.Query(ctx, &query, map[string]interface{}{
		"id":          sg.AltID,
		"provider_id": sg.Provider.altId,
	}); err != nil {
		return nil, err
	}

	sg.Uuid = &query.ProviderSelf.ScriptGroup.Uuid
	return &query.ProviderSelf.ScriptGroup, nil
}

func (sg *ScriptGroup) Create(ctx context.Context, profile CreateScriptGroupInput) (*uuid.UUID, error) {
	var mutation gql.CreateNewScriptGroup
	err := sg.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
		"provider_id": sg.Provider.ID(),
		"alt_id":      sg.AltID.String(),
		"name":        profile.Name,
		"description": profile.Description,
		"public":      profile.Public,
	})

	if err != nil {
		return nil, err
	}

	sg.Uuid = &mutation.Provider.ScriptGroups.Create.Uuid

	return sg.Uuid, nil
}

func (sg *ScriptGroup) Delete(ctx context.Context) error {
	var public = false
	_, err := sg.Update(ctx, UpdateScriptGroupInput{
		Public: &public,
	})
	return err
}

func (sgi *UpdateScriptGroupInput) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	if sgi.AltID != nil {
		data["alt_id"] = *sgi.AltID
	}
	if sgi.Name != nil {
		data["name"] = *sgi.Name
	}
	if sgi.Description != nil {
		data["description"] = *sgi.Description
	}
	if sgi.Public != nil {
		data["public"] = *sgi.Public
	}

	return json.Marshal(data)
}

package provider

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Pritch009/myscribae-sdk-go/gql"
	"github.com/Pritch009/myscribae-sdk-go/utilities"
	"github.com/google/uuid"
)

type Script struct {
	ScriptGroup *ScriptGroup
	Uuid        *uuid.UUID
	Profile     ScriptInput
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

func (s *Script) Printf(format string, a ...interface{}) {
	log.Printf(fmt.Sprintf("[%s] %s", s.Profile.AltID, format), a...)
}

func (s *Script) Println(a ...interface{}) {
	log.Println(fmt.Sprintf("[%s]", s.Profile.AltID), a)
}

func (s *Script) Sync(ctx context.Context, remoteScript *gql.RemoteScript) {
	s.Println("Syncing script")

	if remoteScript != nil {
		hash := utilities.VersionHash(
			[]sql.NullString{
				utilities.NullString(&s.Profile.AltID),
				utilities.NullString(&s.Profile.Name),
				utilities.NullString(&s.Profile.Description),
				utilities.NotNullString(fmt.Sprintf("%d", s.Profile.PriceInCents)),
				utilities.NotNullString(s.Profile.Recurrence),
				utilities.NotNullString(fmt.Sprintf("%d", s.Profile.TokenLifetimeSec)),
				utilities.NotNullString(fmt.Sprintf("%d", s.Profile.SlaSec)),
				utilities.NotNullString(fmt.Sprintf("%t", s.Profile.Public)),
			},
		)

		if remoteScript.Version != hash {
			s.Println("updating script")

			var mutation gql.EditScript
			err := s.ScriptGroup.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
				"alt_id":             s.Profile.AltID,
				"name":               s.Profile.Name,
				"description":        s.Profile.Description,
				"recurrence":         s.Profile.Recurrence,
				"price_in_cents":     s.Profile.PriceInCents,
				"sla_sec":            s.Profile.SlaSec,
				"token_lifetime_sec": s.Profile.TokenLifetimeSec,
				"public":             s.Profile.Public,
			})
			if err != nil {
				panic("failed to update script: " + err.Error())
			}

			s.Uuid = &mutation.Provider.ScriptGroup.Script.Edit.Uuid
			s.Println("Script updated.")
		} else {
			s.Uuid = &remoteScript.Uuid
		}
	} else {
		log.Printf("Creating new script %s", s.Profile.Name)

		var mutation gql.CreateNewScript
		err := s.ScriptGroup.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
			"script_group_id":    s.ScriptGroup.Uuid,
			"alt_id":             s.Profile.AltID,
			"name":               s.Profile.Name,
			"description":        s.Profile.Description,
			"recurrence":         s.Profile.Recurrence,
			"price_in_cents":     s.Profile.PriceInCents,
			"sla_sec":            s.Profile.SlaSec,
			"token_lifetime_sec": s.Profile.TokenLifetimeSec,
			"public":             s.Profile.Public,
		})

		if err != nil {
			panic("failed to create script: " + err.Error())
		}

		s.Uuid = &mutation.Provider.ScriptGroup.Scripts.Create.Uuid
	}

	log.Println("Script up to date")
}

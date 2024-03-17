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

type ScriptGroup struct {
	Uuid     *uuid.UUID
	Provider *Provider
	Profile  ScriptGroupInput
	Scripts  []Script
}

func (sg *ScriptGroup) Printf(format string, a ...interface{}) {
	log.Printf(fmt.Sprintf("[%s] %s", sg.Profile.AltID, format), a...)
}

func (sg *ScriptGroup) Println(a ...interface{}) {
	log.Println(fmt.Sprintf("[%s]", sg.Profile.AltID), a)
}

func (sg *ScriptGroup) Sync(ctx context.Context, remoteScriptGroup *gql.RemoteScriptGroup) {
	sg.Printf("Syncing script group")

	if remoteScriptGroup != nil {
		// hash the local script group and compare to the remote
		hash := utilities.VersionHash(
			[]sql.NullString{
				utilities.NullString(&sg.Profile.AltID),
				utilities.NullString(&sg.Profile.Name),
				utilities.NullString(&sg.Profile.Description),
				utilities.NotNullString(fmt.Sprintf("%t", sg.Provider.Public)),
			},
		)

		if remoteScriptGroup.Version != hash {
			sg.Printf("Updating script group")

			var mutation gql.EditScriptGroup
			err := sg.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
				"id":          remoteScriptGroup.Uuid,
				"name":        sg.Profile.Name,
				"description": sg.Profile.Description,
				"public":      sg.Profile.Public,
			})
			if err != nil {
				panic("failed to update script group: " + err.Error())
			}

			sg.Uuid = &mutation.Provider.ScriptGroup.Edit.Uuid
			sg.Printf("Updated script group")
		} else {
			sg.Printf("Script group is up to date")
		}
	} else {
		sg.Printf("Creating script group %s", sg.Profile.Name)

		var mutation gql.CreateNewScriptGroup
		err := sg.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
			"provider_ sid": sg.Provider.Uuid,
			"alt_id":        sg.Profile.AltID,
			"name":          sg.Profile.Name,
			"description":   sg.Profile.Description,
			"public":        sg.Profile.Public,
		})

		if err != nil {
			panic("failed to create script group: " + err.Error())
		}

		sg.Uuid = &mutation.Provider.ScriptGroups.Create.Uuid
	}

	sg.Printf("Script group up to date")
	existingScripts := make(map[string]*gql.RemoteScript)
	for _, rs := range remoteScriptGroup.Scripts {
		existingScripts[rs.AltID] = &rs
	}

	scriptsSeen := make(map[string]bool)
	for _, s := range sg.Scripts {
		if _, ok := scriptsSeen[s.Profile.AltID]; ok {
			sg.Println("Script exists and has already been seen")
			continue
		} else {
			sg.Println("Syncing script")
		}

		s.ScriptGroup = sg
		scriptsSeen[s.Profile.AltID] = true
		rs := existingScripts[s.Profile.AltID]
		s.Sync(ctx, rs)
	}

	for _, rs := range remoteScriptGroup.Scripts {
		if _, ok := scriptsSeen[rs.AltID]; !ok {
			sg.Printf("Script %s has been abandoned", rs.AltID)

			var mutation gql.EditScript
			err := sg.Provider.Client.Mutate(ctx, &mutation, map[string]interface{}{
				"id":     rs.Uuid,
				"public": false,
			})

			if err != nil {
				panic("failed to abandon script: " + err.Error())
			}
		}
	}

	sg.Printf("Script group synced")
}

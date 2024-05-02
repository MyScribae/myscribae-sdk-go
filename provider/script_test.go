package provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/myscribae/myscribae-sdk-go/provider"
	"github.com/myscribae/myscribae-sdk-go/test_utils"
	"github.com/myscribae/myscribae-sdk-go/utilities"
)

func TestScriptCreateAndEdit(t *testing.T) {
	test_utils.SetupTest(t, func(ctx context.Context, prov *provider.Provider) error {
		// create a test script group
		sg, err := prov.ScriptGroup("test")
		if err != nil {
			return err
		}

		createOptions := provider.CreateScriptGroupInput{
			Name:        "Test Script Group",
			Description: "This is a test script group",
			Public:      true,
		}

		scriptGroupUuid, err := sg.Create(ctx, createOptions)
		if err != nil {
			return err
		}

		if scriptGroupUuid == nil {
			return fmt.Errorf("script group uuid is nil")
		}

		// create a test script
		script, err := prov.Script(sg.AltID, "test")
		if err != nil {
			return err
		}
		createScriptOptions := provider.CreateScriptInput{
			AltID:            "test",
			Name:             "Test Script",
			Description:      "This is a test script",
			Recurrence:       "monthly",
			PriceInCents:     100,
			SlaSec:           utilities.NewUInt(100),
			TokenLifetimeSec: utilities.NewUInt(100),
			Public:           true,
		}

		scriptUuid, err := script.Create(ctx, createScriptOptions)
		if err != nil {
			return err
		}

		if scriptUuid == nil {
			return fmt.Errorf("script uuid is nil")
		}

		// read the test script
		scriptData, err := script.Read(ctx)
		if err != nil {
			return err
		}

		if scriptData == nil {
			return fmt.Errorf("script data is nil")
		}

		if scriptData.Name != createScriptOptions.Name {
			return fmt.Errorf("script name is incorrect")
		}

		// edit the test script
		var (
			newName                  = "Test Script - Edited"
			newDescription           = "This is a test script - edited"
			newPriceInCents          = utilities.NewCentValue(200)
			newSlaSec           uint = 200
			newTokenLifetimeSec uint = 200
			newPublic                = false
		)

		_, err = script.Update(ctx, provider.UpdateScriptInput{
			Name:             &newName,
			Description:      &newDescription,
			PriceInCents:     &newPriceInCents,
			SlaSec:           utilities.NewUIntPointer(&newSlaSec),
			TokenLifetimeSec: utilities.NewUIntPointer(&newTokenLifetimeSec),
			Public:           &newPublic,
		})

		if err != nil {
			return err
		}

		scriptDataNew, err := script.Read(ctx)
		if err != nil {
			return err
		}

		if scriptDataNew == nil {
			return fmt.Errorf("script data is nil")
		}

		if scriptDataNew.Name != newName {
			return fmt.Errorf("script name is incorrect")
		}

		if scriptDataNew.Description != newDescription {
			return fmt.Errorf("script description is incorrect")
		}

		if scriptDataNew.PriceInCents != newPriceInCents.CentsValue() {
			return fmt.Errorf("script price in cents is incorrect")
		}

		if scriptDataNew.SlaSec != newSlaSec {
			return fmt.Errorf("script sla_sec is incorrect")
		}

		if scriptDataNew.TokenLifetimeSec != newTokenLifetimeSec {
			return fmt.Errorf("script token_lifetime_sec is incorrect")
		}

		if scriptDataNew.Public != newPublic {
			return fmt.Errorf("script public is incorrect")
		}

		// delete the test script
		err = script.Delete(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

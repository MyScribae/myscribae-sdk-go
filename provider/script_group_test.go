package provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/myscribae/myscribae-sdk-go/provider"
	"github.com/myscribae/myscribae-sdk-go/test_utils"
)

func TestCreateAndEditScriptGroup(t *testing.T) {
	test_utils.SetupTest(t, func(ctx context.Context, prov *provider.Provider) error {
		create_options := provider.CreateScriptGroupInput{
			Name:        "Test Script Group",
			Description: "This is a test script group",
			Public:      true,
		}

		sg, err := prov.ScriptGroup("test")
		if err != nil {
			return err
		}

		// create a test script group
		scriptGroupUuid, err := sg.Create(ctx, create_options)
		if err != nil {
			return err
		}

		if scriptGroupUuid == nil {
			return fmt.Errorf("script group uuid is nil")
		}

		// read the test script group
		groupData, err := sg.Read(ctx)
		if err != nil {
			return err
		}

		if groupData == nil {
			return fmt.Errorf("script group data is nil")
		}

		if groupData.Name != create_options.Name {
			return fmt.Errorf("script group name is incorrect")
		}

		// edit the test script group
		var (
			newName 	  = "Test Script Group - Edited"
			newDescription = "This is a test script group - edited"
			newPublic 	  = false
		)

		_, err = sg.Update(ctx, provider.UpdateScriptGroupInput{
			Name:        &newName,
			Description: &newDescription,
			Public:      &newPublic,
		})
		if err != nil {
			return err
		}

		groupDataNew, err := sg.Read(ctx)
		if err != nil {
			return err
		}
		if groupDataNew == nil {
			return fmt.Errorf("script group data is nil")
		}

		if groupDataNew.Name != newName {
			return fmt.Errorf("script group name is incorrect")
		}

		if groupDataNew.Description != newDescription {
			return fmt.Errorf("script group description is incorrect")
		}

		if groupDataNew.Public != newPublic {
			return fmt.Errorf("script group public is incorrect")
		}

		return nil
	})
}

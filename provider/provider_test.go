package provider_test

import (
	"context"
	"fmt"

	"testing"

	"github.com/myscribae/myscribae-sdk-go/provider"
	"github.com/myscribae/myscribae-sdk-go/test_utils"
)

func TestProviderEdit(t *testing.T) {
	test_utils.SetupTest(t, func(ctx context.Context, prov *provider.Provider) error {
		// attempt to create a get provider with an existing client
		provData, err := prov.Read(ctx)
		if err != nil {
			return fmt.Errorf("failed to read provider: %v", err)
		}

		if provData == nil {
			return fmt.Errorf("provider data is nil")
		}

		// attempt to change the provider name
		originalName := provData.Name
		newName := fmt.Sprintf("%s - %s", originalName, "Test")

		provUuid, err := prov.Update(ctx, provider.UpdateProviderProfileInput{
			Name: &newName,
		})

		if err != nil {
			return fmt.Errorf("failed to update provider: %v", err)
		}

		if provUuid == nil {
			return fmt.Errorf("provider uuid is nil")
		}

		if provUuid.String() != test_utils.TestProviderUuid {
			return fmt.Errorf("provider uuid is incorrect")
		}

		// attempt to read the provider again
		newProvData, err := prov.Read(ctx)
		if err != nil {
			return fmt.Errorf("failed to read provider: %v", err)
		}

		if newProvData == nil {
			return fmt.Errorf("provider data is nil")
		}

		if newProvData.Name != newName {
			return fmt.Errorf("provider name is incorrect")
		}

		// attempt to change the provider name back
		_, err = prov.Update(ctx, provider.UpdateProviderProfileInput{
			Name: &originalName,
		})
		if err != nil {
			return fmt.Errorf("failed to update provider: %v", err)
		}

		return nil
	})
}

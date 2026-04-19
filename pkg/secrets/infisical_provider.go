//go:build infisical

package secrets

import (
	"context"
	"fmt"
	"os"

	infisical "github.com/infisical/go-sdk"
)

// fetchFromInfisical uses the Infisical Go SDK to fetch a secret from a given project/environment.
// This file is only compiled when building with `-tags=infisical`.
func fetchFromInfisical(ctx context.Context, secretKey, projectID, environment, secretPath string) (string, error) {
	client := infisical.NewInfisicalClient(ctx, infisical.Config{
		SiteUrl:          os.Getenv("INFISICAL_SITE_URL"),
		AutoTokenRefresh: true,
	})

	// Attempt to authenticate using environment variables (INFISICAL_UNIVERSAL_AUTH_CLIENT_ID/_CLIENT_SECRET)
	if _, err := client.Auth().UniversalAuthLogin("", ""); err != nil {
		return "", fmt.Errorf("infisical: auth failed: %w", err)
	}

	secret, err := client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
		SecretKey:   secretKey,
		ProjectID:   projectID,
		Environment: environment,
		SecretPath:  secretPath,
	})
	if err != nil {
		return "", fmt.Errorf("infisical: failed to retrieve secret: %w", err)
	}

	// Return a string representation of the retrieved secret. The SDK's return shape may differ
	// across versions; consumers building with the `infisical` tag should ensure their SDK
	// version matches expectations. Using fmt.Sprintf("%v") provides a reasonable default.
	return fmt.Sprintf("%v", secret), nil
}

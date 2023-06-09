package key

import (
	"context"
	"github.com/lugondev/tx-builder/pkg/errors"
	authutils "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/utils"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/multitenancy"
)

// Key is a Checker for API Key authentication
type Key struct {
	key string
}

func New(key string) *Key {
	return &Key{
		key: key,
	}
}

// Check Parse and verify the validity of the Token (UUID or Access) and return a struct for a JWT (JSON Web Token)
func (c *Key) Check(ctx context.Context) (*multitenancy.UserInfo, error) {
	if c == nil || c.key == "" {
		return nil, nil
	}

	// Extract Key from context
	apiKeyCtx := authutils.APIKeyFromContext(ctx)
	if apiKeyCtx == "" {
		return nil, nil
	}

	if apiKeyCtx != c.key {
		return nil, errors.UnauthorizedError("invalid API key")
	}

	userInfo := multitenancy.NewAPIKeyUserInfo(apiKeyCtx)
	err := userInfo.ImpersonateTenant(authutils.TenantIDFromContext(ctx))
	if err != nil {
		return nil, err
	}

	// Impersonate username
	err = userInfo.ImpersonateUsername(authutils.UsernameFromContext(ctx))
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

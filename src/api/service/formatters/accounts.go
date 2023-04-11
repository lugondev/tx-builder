package formatters

import (
	"net/http"
	"strings"

	"github.com/lugondev/tx-builder/src/api/service/types"
	infra "github.com/lugondev/tx-builder/src/infra/api"

	"github.com/lugondev/tx-builder/src/entities"
)

func FormatCreateAccountRequest(req *types.CreateAccountRequest, defaultStoreID string) *entities.Account {
	acc := &entities.Account{
		Alias:      req.Alias,
		Attributes: req.Attributes,
		StoreID:    req.StoreID,
	}

	if acc.StoreID == "" {
		acc.StoreID = defaultStoreID
	}

	return acc
}

func FormatImportAccountRequest(req *types.ImportAccountRequest, defaultStoreID string) *entities.Account {
	acc := &entities.Account{
		Alias:      req.Alias,
		Attributes: req.Attributes,
		StoreID:    req.StoreID,
	}

	if acc.StoreID == "" {
		acc.StoreID = defaultStoreID
	}

	return acc
}

func FormatUpdateAccountRequest(req *types.UpdateAccountRequest) *entities.Account {
	return &entities.Account{
		Alias:      req.Alias,
		Attributes: req.Attributes,
		StoreID:    req.StoreID,
	}
}

func FormatAccountResponse(iden *entities.Account) *types.AccountResponse {
	return &types.AccountResponse{
		Alias:               iden.Alias,
		Attributes:          iden.Attributes,
		PublicKey:           iden.PublicKey.String(),
		CompressedPublicKey: iden.CompressedPublicKey.String(),
		TenantID:            iden.TenantID,
		OwnerID:             iden.OwnerID,
		StoreID:             iden.StoreID,
		CreatedAt:           iden.CreatedAt,
		UpdatedAt:           iden.UpdatedAt,
	}
}

func FormatAccountFilterRequest(req *http.Request) (*entities.AccountFilters, error) {
	filters := &entities.AccountFilters{}

	qAliases := req.URL.Query().Get("aliases")
	if qAliases != "" {
		filters.Aliases = strings.Split(qAliases, ",")
	}

	if err := infra.GetValidator().Struct(filters); err != nil {
		return nil, err
	}

	return filters, nil
}

package formatters

import (
	"github.com/lugondev/tx-builder/src/api/service/types"
	infra "github.com/lugondev/tx-builder/src/infra/api"
	"net/http"

	"github.com/lugondev/tx-builder/src/entities"
)

func FormatCreateAccountRequest(req *types.CreateAccountRequest, defaultStoreID string) *entities.Wallet {
	acc := &entities.Wallet{
		Attributes: req.Attributes,
		StoreID:    req.StoreID,
	}

	if acc.StoreID == "" {
		acc.StoreID = defaultStoreID
	}

	return acc
}

func FormatImportAccountRequest(req *types.ImportAccountRequest, defaultStoreID string) *entities.Wallet {
	acc := &entities.Wallet{
		Attributes: req.Attributes,
		StoreID:    req.StoreID,
	}

	if acc.StoreID == "" {
		acc.StoreID = defaultStoreID
	}

	return acc
}

func FormatUpdateAccountRequest(req *types.UpdateAccountRequest) *entities.Wallet {
	return &entities.Wallet{
		Attributes: req.Attributes,
		StoreID:    req.StoreID,
	}
}

func FormatAccountResponse(iden *entities.Wallet) *types.AccountResponse {
	return &types.AccountResponse{
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

	if err := infra.GetValidator().Struct(filters); err != nil {
		return nil, err
	}

	return filters, nil
}

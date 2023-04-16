package accounts

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lugondev/tx-builder/pkg/utils"

	qkm "github.com/lugondev/tx-builder/src/infra/signer-key-manager/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/pkg/errors"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/log"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/multitenancy"
	usecases "github.com/lugondev/tx-builder/src/api/business/use-cases"
	"github.com/lugondev/tx-builder/src/api/store"
	"github.com/lugondev/tx-builder/src/entities"
	"github.com/lugondev/wallet-signer-manager/pkg/client"
	qkmtypes "github.com/lugondev/wallet-signer-manager/src/stores/api/types"
)

const createAccountComponent = "use-cases.create-account"

type createAccountUseCase struct {
	db               store.DB
	searchUC         usecases.SearchAccountsUseCase
	keyManagerClient client.KeyManagerClient
	logger           *log.Logger
}

type Account struct {
	Address common.Address
	priv    *ecdsa.PrivateKey
}

func NewCreateAccountUseCase(
	db store.DB,
	searchUC usecases.SearchAccountsUseCase,
	keyManagerClient client.KeyManagerClient,
) usecases.CreateAccountUseCase {
	return &createAccountUseCase{
		db:               db,
		searchUC:         searchUC,
		keyManagerClient: keyManagerClient,
		logger:           log.NewLogger().SetComponent(createAccountComponent),
	}
}

func (uc *createAccountUseCase) Execute(ctx context.Context, acc *entities.Wallet, privateKey hexutil.Bytes, userInfo *multitenancy.UserInfo) (*entities.Wallet, error) {
	logger := uc.logger.WithContext(ctx)
	logger.Debug("creating new wallet")

	accounts, err := uc.searchUC.Execute(ctx,
		&entities.AccountFilters{ID: 9999, TenantID: userInfo.TenantID, OwnerID: userInfo.Username},
		userInfo)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	if len(accounts) > 0 {
		errMsg := "wallet already exists"
		logger.Error(errMsg)
		return nil, errors.AlreadyExistsError(errMsg).ExtendComponent(createAccountComponent)
	}

	var accountID = utils.GenerateKeyID()
	var resp *qkmtypes.WalletResponse
	if privateKey != nil {
		importedAccount, der := NewAccountFromPrivateKey(privateKey.String())
		if der != nil {
			logger.WithError(err).Error("invalid private key")
			return nil, errors.InvalidParameterError(der.Error()).ExtendComponent(createAccountComponent)
		}

		existingAcc, der := uc.db.Account().FindOneByPubkey(ctx, importedAccount.Address.Hex(), userInfo.AllowedTenants, userInfo.Username)
		if existingAcc != nil {
			errMsg := "account already exists"
			logger.Error(errMsg)
			return nil, errors.AlreadyExistsError(errMsg).ExtendComponent(createAccountComponent)
		}

		if der != nil && !errors.IsNotFoundError(der) {
			errMsg := "failed to get account"
			logger.WithError(der).Error(errMsg)
			return nil, errors.FromError(der).ExtendComponent(createAccountComponent)
		}

		resp, err = uc.keyManagerClient.ImportWallet(ctx, acc.StoreID, &qkmtypes.ImportWalletRequest{
			KeyID:      accountID,
			PrivateKey: privateKey,
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		})
	} else {
		resp, err = uc.keyManagerClient.CreateWallet(ctx, acc.StoreID, &qkmtypes.CreateWalletRequest{
			KeyID: accountID,
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		})
	}

	if err != nil {
		errMsg := "failed to import/create wallet"
		logger.WithError(err).Error(errMsg)
		return nil, errors.DependencyFailureError(errMsg).ExtendComponent(createAccountComponent)
	}

	acc.PublicKey = common.FromHex(resp.PublicKey)
	acc.CompressedPublicKey = common.FromHex(resp.CompressedPublicKey)
	acc.TenantID = userInfo.TenantID
	acc.OwnerID = userInfo.Username

	acc, err = uc.db.Account().Insert(ctx, acc)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	logger.WithField("address", resp.CompressedPublicKey).Info("wallet created successfully")
	return acc, nil
}

func NewAccountFromPrivateKey(priv string) (*Account, error) {
	prv, err := crypto.HexToECDSA(priv[2:])
	if err != nil {
		return nil, err
	}

	return &Account{
		priv:    prv,
		Address: crypto.PubkeyToAddress(prv.PublicKey),
	}, nil
}

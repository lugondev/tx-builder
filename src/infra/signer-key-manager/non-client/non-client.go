package nonclient

import (
	"context"
	"github.com/lugondev/tx-builder/pkg/errors"
	qkm "github.com/lugondev/wallet-signer-manager/pkg/client"
	"github.com/lugondev/wallet-signer-manager/pkg/jsonrpc"
	storestypes "github.com/lugondev/wallet-signer-manager/src/stores/api/types"
)

const errMessage = "Wallet Signer Manager is disabled"

type NonClient struct{}

var _ qkm.KeyManagerClient = &NonClient{}

func NewNonClient() *NonClient {
	return &NonClient{}
}

func (n NonClient) CreateWallet(ctx context.Context, storeName string, request *storestypes.CreateWalletRequest) (*storestypes.WalletResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) ImportWallet(ctx context.Context, storeName string, request *storestypes.ImportWalletRequest) (*storestypes.WalletResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) UpdateWallet(ctx context.Context, storeName, pubkey string, request *storestypes.UpdateWalletRequest) (*storestypes.WalletResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) GetWallet(ctx context.Context, storeName, pubkey string) (*storestypes.WalletResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) ListWallets(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) ListDeletedWallets(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) DeleteWallet(ctx context.Context, storeName, pubkey string) error {
	return errors.DependencyFailureError(errMessage)
}
func (n NonClient) DestroyWallet(ctx context.Context, storeName, pubkey string) error {
	return errors.DependencyFailureError(errMessage)
}
func (n NonClient) RestoreWallet(ctx context.Context, storeName, pubkey string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) Sign(ctx context.Context, storeName, account string, request *storestypes.SignWalletRequest) (string, error) {
	return "", errors.DependencyFailureError(errMessage)
}

func (n NonClient) Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

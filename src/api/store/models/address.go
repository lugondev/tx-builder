package models

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/lugondev/tx-builder/pkg/blockchain/bitcoin/chain"
	"github.com/lugondev/tx-builder/pkg/common"
	"github.com/samber/lo"
	"time"

	"github.com/lugondev/tx-builder/src/entities"
)

type Address struct {
	tableName struct{} `pg:"addresses"` // nolint:unused,structcheck // reason

	ID         int
	WalletID   int
	WalletType entities.AddressType
	Address    string

	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
}

func NewAddressesFromWallet(wallet *entities.Wallet) []*Address {
	return lo.Map[entities.AddressType, *Address](entities.AddressTypes, func(addressType entities.AddressType, _ int) *Address {
		return getAddress(wallet, addressType)
	})
}

func getAddress(wallet *entities.Wallet, addressType entities.AddressType) *Address {
	// generate switch case with addressTypes
	address := &Address{
		WalletType: addressType,
		WalletID:   wallet.ID,
		CreatedAt:  wallet.CreatedAt,
		UpdatedAt:  wallet.UpdatedAt,
	}
	pubkeyBtc, err := btcec.ParsePubKey(wallet.PublicKey)
	if err != nil {
		return nil
	}

	btcAddressesTestnet := chain.PubkeyToAddresses(pubkeyBtc, &chaincfg.TestNet3Params)
	btcAddressesMainnet := chain.PubkeyToAddresses(pubkeyBtc, &chaincfg.MainNetParams)
	switch addressType {
	case entities.AddressTypeBTCSegwit:
		address.Address = btcAddressesMainnet[common.Segwit]
	case entities.AddressTypeBTCLegacy:
		address.Address = btcAddressesMainnet[common.Legacy]
	case entities.AddressTypeBTCTaproot:
		address.Address = btcAddressesMainnet[common.Taproot]
	case entities.AddressTypeBTCNested:
		address.Address = btcAddressesMainnet[common.Nested]
	case entities.AddressTypeBTCSegwitTestnet:
		address.Address = btcAddressesTestnet[common.Segwit]
	case entities.AddressTypeBTCLegacyTestnet:
		address.Address = btcAddressesTestnet[common.Legacy]
	case entities.AddressTypeBTCTaprootTestnet:
		address.Address = btcAddressesTestnet[common.Taproot]
	case entities.AddressTypeBTCNestedTestnet:
		address.Address = btcAddressesTestnet[common.Nested]
	case entities.AddressTypeEVM:
		address.Address = ethcommon.BytesToAddress(pubkeyBtc.SerializeUncompressed()).String()
	default:
		return nil
	}

	return address
}

func (a *Address) ToEntity() *entities.Address {
	return &entities.Address{
		WalletType: a.WalletType,
		WalletID:   a.WalletID,
		Address:    a.Address,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

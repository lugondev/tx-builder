package author

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/txscript"
)

var _ txauthor.SecretsSource = (*MemorySecretStoreWif)(nil)

func NewMemorySecretStoreWif(keyMap map[string]*btcutil.WIF, params *chaincfg.Params) MemorySecretStoreWif {
	return MemorySecretStoreWif{
		keyMap: keyMap,
		params: params,
	}
}

type MemorySecretStoreWif struct {
	keyMap map[string]*btcutil.WIF
	params *chaincfg.Params
}

func (m MemorySecretStoreWif) GetKey(address btcutil.Address) (*btcec.PrivateKey, bool, error) {
	wif, found := m.keyMap[address.EncodeAddress()]
	if !found {
		return nil, false, fmt.Errorf("address not found")
	}
	return wif.PrivKey, wif.CompressPubKey, nil
}

func (m MemorySecretStoreWif) GetScript(address btcutil.Address) ([]byte, error) {
	return txscript.PayToAddrScript(address)
}

func (m MemorySecretStoreWif) ChainParams() *chaincfg.Params {
	return m.params
}

package author

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
)

var _ txauthor.SecretsSource = (*MemorySecretStore)(nil)

func NewMemorySecretStore(keyMap map[string]*btcec.PrivateKey, params *chaincfg.Params) MemorySecretStore {
	return MemorySecretStore{
		keyMap: keyMap,
		params: params,
	}
}

type MemorySecretStore struct {
	keyMap map[string]*btcec.PrivateKey
	params *chaincfg.Params
}

func (m MemorySecretStore) GetKey(address btcutil.Address) (*btcec.PrivateKey, bool, error) {
	privKey, found := m.keyMap[address.EncodeAddress()]
	if !found {
		return nil, false, fmt.Errorf("address not found")
	}
	return privKey, true, nil
}

func (m MemorySecretStore) GetScript(address btcutil.Address) ([]byte, error) {
	return txscript.PayToAddrScript(address)
}

func (m MemorySecretStore) ChainParams() *chaincfg.Params {
	return m.params
}

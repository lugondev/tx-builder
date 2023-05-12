package author

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	ecdsa2 "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lugondev/tx-builder/pkg/blockchain/bitcoin"
	"github.com/lugondev/tx-builder/pkg/blockchain/bitcoin/txscript"
)

var _ SecretsSource = (*MemorySecretStore)(nil)

func NewMemorySecretStore(addressMap map[string]*btcec.PrivateKey, pubkeyMap map[string][]byte, params *chaincfg.Params) MemorySecretStore {
	return MemorySecretStore{
		addressMap: addressMap,
		pubkeyMap:  pubkeyMap,
		params:     params,
	}
}

type MemorySecretStore struct {
	addressMap map[string]*btcec.PrivateKey
	pubkeyMap  map[string][]byte
	params     *chaincfg.Params
}

func (m MemorySecretStore) GetKey(address btcutil.Address) (*btcec.PrivateKey, bool, error) {
	privKey, found := m.addressMap[address.EncodeAddress()]
	if !found {
		return nil, false, fmt.Errorf("address not found")
	}
	return privKey, true, nil
}

func (m MemorySecretStore) GetPubkey(address btcutil.Address) ([]byte, bool, error) {
	pubkey, found := m.pubkeyMap[address.EncodeAddress()]
	if !found {
		return nil, false, fmt.Errorf("pubkey not found")
	}
	return pubkey, true, nil
}

func (m MemorySecretStore) Sign(pubkey []byte, data []byte) ([]byte, error) {
	privKey, found := m.addressMap[hexutil.Encode(pubkey)]
	if !found {
		return nil, fmt.Errorf("pubkey not found: %s", hexutil.Encode(pubkey))
	}
	sig := ecdsa2.Sign(privKey, data)
	//sig, err := bitcoin.MpcSign(data)
	//if err != nil {
	//	return nil, err
	//}
	////sig := sig1.Serialize()
	//fmt.Println("len sig", len(sig))
	//fmt.Println("hex sig", hex.EncodeToString(sig))
	return sig.Serialize(), nil
}

func (m MemorySecretStore) SignTaproot(pubkey []byte, data []byte) (*schnorr.Signature, error) {
	fmt.Println("sign taproot", hex.EncodeToString(data))
	fmt.Println("sign taproot pubkey", hexutil.Encode(pubkey))
	privKey, found := m.addressMap[hexutil.Encode(pubkey)]
	if !found {
		return nil, fmt.Errorf("pubkey not found: %s", hexutil.Encode(pubkey))
	}
	ecdsaPrivKey, err := crypto.HexToECDSA(hex.EncodeToString(privKey.Serialize()))
	if err != nil {
		return nil, err
	}
	sig, err := bitcoin.SignTaprootSignature(data, ecdsaPrivKey)
	if err != nil {
		return nil, err
	}
	return schnorr.ParseSignature(sig)
}

func (m MemorySecretStore) GetScript(address btcutil.Address) ([]byte, error) {
	return txscript.PayToAddrScript(address)
}

func (m MemorySecretStore) ChainParams() *chaincfg.Params {
	return m.params
}

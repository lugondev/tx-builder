// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package builder

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lugondev/tx-builder/pkg/blockchain/bitcoin/chain"
	"github.com/lugondev/tx-builder/pkg/common"
	"github.com/lugondev/tx-builder/pkg/xcrypto"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMpcEcdsa(t *testing.T) {
	hash := xcrypto.DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	// Party 1.
	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := xcrypto.PrvKeyFromBytes(p1.Bytes())
	pub1 := prv1.PubKey()
	party1 := xcrypto.NewEcdsaParty(prv1)
	defer party1.Close()

	// Party 2.
	p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	prv2 := xcrypto.PrvKeyFromBytes(p2.Bytes())
	pub2 := prv2.PubKey()
	party2 := xcrypto.NewEcdsaParty(prv2)
	defer party2.Close()

	// Phase 1.
	sharepub1 := party1.Phase1(pub2)
	sharepub2 := party2.Phase1(pub1)
	assert.Equal(t, sharepub1, sharepub2)
	t.Log("sharepub1:", hex.EncodeToString(sharepub1.Serialize()))
	pubkey, err := btcec.ParsePubKey(sharepub1.Serialize())
	if err != nil {
		t.Fatal(err)
	}
	addresses := chain.PubkeyToAddresses(pubkey, &chaincfg.TestNet3Params)
	t.Log("address legacy:", addresses[common.Legacy])

	// Phase 2.
	encpk1, encpub1, scalarR1 := party1.Phase2(hash)
	encpk2, encpub2, scalarR2 := party2.Phase2(hash)

	// Phase 3.
	shareR1 := party1.Phase3(scalarR2)
	shareR2 := party2.Phase3(scalarR1)
	assert.Equal(t, shareR1, shareR2)

	// Phase 4.
	sig1, err := party1.Phase4(encpk2, encpub2, shareR1)
	assert.Nil(t, err)
	sig2, err := party2.Phase4(encpk1, encpub1, shareR2)
	assert.Nil(t, err)

	// Phase 5.
	fs1, err := party1.Phase5(shareR1, sig2)
	assert.Nil(t, err)
	fs2, err := party2.Phase5(shareR2, sig1)
	assert.Nil(t, err)
	assert.Equal(t, fs1, fs2)

	// Verify.
	err = xcrypto.EcdsaVerify(sharepub1, hash, fs1)
	assert.Nil(t, err)
	t.Logf("\nKeys\n  x1: %x\n  x2: %x\n  Q:  %x\n\nSignatures\n  %x\nIs valid under Q?: %v",
		p1.Bytes(),
		p2.Bytes(),
		sharepub1.SerializeCompressed(),
		fs1,
		err == nil)
}

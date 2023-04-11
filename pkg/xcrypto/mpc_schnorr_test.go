// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcrypto

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMpcSchnorr(t *testing.T) {
	hash := DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	// Party 1.
	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := PrvKeyFromBytes(p1.Bytes())
	pub1 := prv1.PubKey()
	party1, err := NewSchnorrParty(prv1)
	assert.Nil(t, err)
	defer party1.Close()

	// Party 2.
	p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	prv2 := PrvKeyFromBytes(p2.Bytes())
	pub2 := prv2.PubKey()
	party2, err := NewSchnorrParty(prv2)
	assert.Nil(t, err)
	defer party2.Close()

	// Phase 1.
	sharepub1 := party1.Phase1(pub2)
	sharepub2 := party2.Phase1(pub1)
	assert.Equal(t, sharepub1, sharepub2)
	t.Log("sharepub1:", hex.EncodeToString(sharepub1.Serialize()))

	// Phase 2.
	r1 := party1.Phase2(hash)
	r2 := party2.Phase2(hash)

	// Phase 3.
	sharer1 := party1.Phase3(r2)
	sharer2 := party2.Phase3(r1)
	assert.Equal(t, sharer1, sharer2)

	// Phase 4.
	s1, err := party1.Phase4(sharepub1, sharer1)
	assert.Nil(t, err)
	s2, err := party2.Phase4(sharepub2, sharer2)
	assert.Nil(t, err)

	// Phase 5.
	fs1, err := party1.Phase5(sharer1, s1, s2)
	assert.Nil(t, err)
	fs2, err := party2.Phase5(sharer2, s1, s2)
	assert.Nil(t, err)
	assert.Equal(t, fs1, fs2)
	assert.Nil(t, err)

	assert.Nil(t, SchnorrVerify(sharepub1, fs1, hash))
	//signature, err := schnorr.ParseSignature(fs1)
	//assert.Nil(t, err)
	//parsePubKey, err := btcec.ParsePubKey(common2.FromHex("03201d82ca8f8ebe6542459911a5275e8650504050af538d469fda12568c49ed7b"))
	//assert.Nil(t, err)
	//fmt.Printf("parsePubKey: %x\n", parsePubKey.SerializeCompressed())
	//pubKeySchnorr, err := schnorr.ParsePubKey(schnorr.SerializePubKey(parsePubKey))
	//assert.Nil(t, err)
	//fmt.Printf("pubKeySchnorr: %x\n", pubKeySchnorr.SerializeCompressed())
	//fmt.Printf("sharepub1: %x\n", sharepub1.SerializeCompressed())
	//assert.True(t, signature.Verify(hash, pubKeySchnorr))
}

func BenchmarkMpcSchnorrKeyGen(b *testing.B) {
	// Party 1.
	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := PrvKeyFromBytes(p1.Bytes())
	pub1 := prv1.PubKey()
	party1, _ := NewSchnorrParty(prv1)
	defer party1.Close()

	// Party 2.
	p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	prv2 := PrvKeyFromBytes(p2.Bytes())
	party2, _ := NewSchnorrParty(prv2)
	defer party2.Close()

	for i := 0; i < b.N; i++ {
		party2.Phase1(pub1)
	}
}

func BenchmarkMpcSchnorrSigning(b *testing.B) {
	hash := DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	// Party 1.
	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := PrvKeyFromBytes(p1.Bytes())
	pub1 := prv1.PubKey()
	party1, _ := NewSchnorrParty(prv1)
	defer party1.Close()

	// Party 2.
	p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	prv2 := PrvKeyFromBytes(p2.Bytes())
	pub2 := prv2.PubKey()
	party2, _ := NewSchnorrParty(prv2)
	defer party2.Close()

	// Phase 1.
	sharepub1 := party1.Phase1(pub2)
	sharepub2 := party2.Phase1(pub1)

	for i := 0; i < b.N; i++ {
		// Phase 2.
		r1 := party1.Phase2(hash)
		r2 := party2.Phase2(hash)

		// Phase 3.
		sharer1 := party1.Phase3(r2)
		sharer2 := party2.Phase3(r1)

		// Phase 4.
		s1, _ := party1.Phase4(sharepub1, sharer1)
		s2, _ := party2.Phase4(sharepub2, sharer2)

		// Phase 5.
		if _, err := party1.Phase5(sharer1, s1, s2); err != nil {
			panic(err)
		}
	}
}

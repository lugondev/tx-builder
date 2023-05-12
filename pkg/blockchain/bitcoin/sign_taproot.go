package bitcoin

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func SignTaprootSignature(data []byte, ecdsaKey *ecdsa.PrivateKey) ([]byte, error) {
	key, _ := btcec.PrivKeyFromBytes(ecdsaKey.D.Bytes())

	// Before we sign the sighash, we'll need to apply the taptweak to the
	// private key based on the tapScriptRootHash.
	privKeyTweak := tweakTaprootPrivKey(*key, []byte{})
	//
	//// With the sighash constructed, we can sign it with the specified
	//// private key.
	signature, err := schnorr.Sign(privKeyTweak, data)
	//signature, err := kdb.Sign(tapScriptRootHash, sigHash)
	if err != nil {
		return nil, err
	}

	return signature.Serialize(), nil
}

func tweakTaprootPrivKey(privKey btcec.PrivateKey,
	scriptRoot []byte) *btcec.PrivateKey {
	// If the corresponding public key has an odd y coordinate, then we'll
	// negate the private key as specified in BIP 341.
	privKeyScalar := privKey.Key
	pubKeyBytes := privKey.PubKey().SerializeCompressed()
	if pubKeyBytes[0] == secp.PubKeyFormatCompressedOdd {
		privKeyScalar.Negate()
	}

	// Next, we'll compute the tap tweak hash that commits to the internal
	// key and the merkle script root. We'll snip off the extra parity byte
	// from the compressed serialization and use that directly.
	schnorrKeyBytes := pubKeyBytes[1:]
	tapTweakHash := chainhash.TaggedHash(
		chainhash.TagTapTweak, schnorrKeyBytes, scriptRoot,
	)

	// Map the private key to a ModNScalar which is needed to perform
	// operation mod the curve order.
	var tweakScalar btcec.ModNScalar
	tweakScalar.SetBytes((*[32]byte)(tapTweakHash))

	// Now that we have the private key in its may negated form, we'll add
	// the script root as a tweak. As we're using a ModNScalar all
	// operations are already normalized mod the curve order.
	privTweak := privKeyScalar.Add(&tweakScalar)

	return btcec.PrivKeyFromScalar(privTweak)
}

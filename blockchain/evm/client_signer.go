package evm

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func SignerFunc(chainID *big.Int, signFunc SignFunc) func(common.Address, *types.Transaction) (*types.Transaction, error) {
	return func(addr common.Address, txn *types.Transaction) (*types.Transaction, error) {
		signer := types.LatestSignerForChainID(chainID)
		txHash := signer.Hash(txn)
		sig, err := signFunc(txHash.Bytes()[:])
		if err != nil {
			fmt.Println("sign error", err)
			return nil, err
		}
		if len(sig) == 64 {
			sig = append(sig, 0x0)
		}
		if len(sig) != 65 {
			return nil, errors.New("wrong signature length")
		}

		for j := 0; j < 2; j++ {
			signedTxn, err := txn.WithSignature(signer, sig)

			if err != nil {
				fmt.Println("Error with signature txn", "detail", err)
				return nil, err
			}
			sender, err := types.Sender(signer, signedTxn)
			if sender.String() == addr.String() {
				return signedTxn, nil
			}

			fmt.Println("sender", "addr", sender, "require", addr)
			vPos := crypto.SignatureLength - 1
			sig[vPos] ^= 0x1
		}

		return nil, errors.New("wrong sender address")
	}
}

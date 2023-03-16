package evm

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// mumbaiClientRPCURL is the RPC URL used by default, to interact with the ethereum node.
	mumbaiClientRPCURL = "https://rpc.ankr.com/polygon_mumbai"
)

// Client holds the underlying RPC client instance.
type Client struct {
	EthClient *ethclient.Client
	ChainID   *big.Int
	Ctx       context.Context
}

// NewClientMumbai creates and returns a new JSON-RPC client to the Mumbai-Polygon node
func NewClientMumbai() (*Client, error) {
	return NewClient(mumbaiClientRPCURL, big.NewInt(80001))
}

// NewClient creates and returns a new JSON-RPC client to the EVM node
func NewClient(rpcURL string, chainID *big.Int) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dialing url: %v", rpcURL)
	}

	ctx := context.Background()
	clientChainID, err := client.ChainID(ctx)
	if err != nil || clientChainID.Cmp(chainID) != 0 {
		return nil, fmt.Errorf("mismatched chain id: expected %v, got %v", chainID, clientChainID)
	}
	return &Client{
		client,
		chainID,
		ctx,
	}, nil
}

// LatestBlock returns the block number at the current chain head.
func (client *Client) LatestBlock(ctx context.Context) (uint64, error) {
	header, err := client.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("fetching header: %v", err)
	}
	return header.Number.Uint64(), nil
}

type SignFunc func(txHash []byte) ([]byte, error)

func (client *Client) Transfer(txRequest *TxRequest, signFunc SignFunc) (*types.Transaction, error) {
	tx, err := txRequest.PrepareTransaction(client.Ctx, client)
	if err != nil {
		return nil, err
	}

	signer := types.LatestSignerForChainID(client.ChainID)
	txHash := signer.Hash(tx)
	sig, err := signFunc(txHash.Bytes()[:])
	if err != nil {
		return nil, err
	}
	signedTxn, err := tx.WithSignature(signer, sig)
	if err != nil {
		return nil, err
	}

	return signedTxn, nil
}

func (client *Client) TransactContract(txRequest *TxRequest, signFunc SignFunc) (*types.Transaction, error) {
	tx, err := txRequest.PrepareTransaction(client.Ctx, client)
	if err != nil {
		return nil, err
	}

	opt := &bind.TransactOpts{
		From: txRequest.From,
		Signer: func(addr common.Address, txn *types.Transaction) (*types.Transaction, error) {
			signer := types.LatestSignerForChainID(client.ChainID)
			txHash := signer.Hash(txn)
			sig, err := signFunc(txHash.Bytes()[:])
			if err != nil {
				fmt.Println("sign error", err)
				return nil, err
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
		},

		Nonce:    new(big.Int).SetUint64(tx.Nonce()),
		GasPrice: tx.GasPrice(),
		GasLimit: tx.Gas(),
		Value:    txRequest.Value,
		NoSend:   true,
	}
	bound := bind.NewBoundContract(*txRequest.To, abi.ABI{}, client.EthClient, client.EthClient, client.EthClient)
	transaction, err := bound.RawTransact(opt, txRequest.Data)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// SubmitTx to the underlying blockchain network.
func (client *Client) SubmitTx(tx *types.Transaction) error {
	return client.EthClient.SendTransaction(client.Ctx, tx)
}

// AccountNonce returns the current nonce of the account. This is the nonce to
// be used while building a new transaction.
func (client *Client) AccountNonce(targetAddr common.Address) (*big.Int, error) {
	nonce, err := client.EthClient.NonceAt(client.Ctx, targetAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce for '%v': %v", targetAddr, err)
	}

	return big.NewInt(int64(nonce)), nil
}

// AccountBalance returns the account balance for a given common.
func (client *Client) AccountBalance(targetAddr common.Address) (*big.Int, error) {
	balance, err := client.EthClient.BalanceAt(client.Ctx, targetAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance for '%v': %v", targetAddr, err)
	}

	return balance, nil
}

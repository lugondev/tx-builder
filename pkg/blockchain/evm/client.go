package evm

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ybbus/jsonrpc"
)

const (
	// mumbaiClientRPCURL is the RPC URL used by default, to interact with the ethereum node.
	mumbaiClientRPCURL = "https://rpc.ankr.com/polygon_mumbai"
)

// Client holds the underlying RPC client instance.
type Client struct {
	EthClient *ethclient.Client
	RpcClient jsonrpc.RPCClient
	ChainID   *big.Int
	Ctx       context.Context
}

// TxPoolInspect ethereum transaction pool datatype
type TxPoolInspect struct {
	Pending map[string]map[uint64]string `json:"pending"`
	Queued  map[string]map[uint64]string `json:"queued"`
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

	rpcClient := jsonrpc.NewClient(rpcURL)

	return &Client{
		client,
		rpcClient,
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

// Transfer sends native token to the given address.
func (client *Client) Transfer(txRequest *TxRequest, signFunc SignFunc) (*types.Transaction, error) {
	// checking: value is greater than zero
	if txRequest.Value.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("value must be greater than zero")
	}

	tx, err := txRequest.PrepareTransaction(client)
	if err != nil {
		return nil, err
	}

	return SignerFunc(client.ChainID, signFunc)(txRequest.From, tx)
}

// TransactContract sends data to the given contract address.
func (client *Client) TransactContract(txRequest *TxRequest, signFunc SignFunc) (*types.Transaction, error) {
	tx, err := txRequest.PrepareTransaction(client)
	if err != nil {
		return nil, err
	}

	opt := &bind.TransactOpts{
		From:   txRequest.From,
		Signer: SignerFunc(client.ChainID, signFunc),

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

// PendingAccountNonce returns the current pending nonce of the account
func (client *Client) PendingAccountNonce(targetAddr common.Address) (*big.Int, error) {
	pendingNonceAt, err := client.EthClient.PendingNonceAt(client.Ctx, targetAddr)
	if err != nil {
		return nil, err
	}

	response, err := client.RpcClient.Call("txpool_inspect")
	if err != nil {
		fmt.Println("txpool_inspect error", err)
		return big.NewInt(int64(pendingNonceAt)), err
	}
	if response.Error != nil {
		fmt.Println("txpool_inspect error", err)
		return big.NewInt(int64(pendingNonceAt)), response.Error
	}

	var (
		txPoolInspect  TxPoolInspect
		txPoolMaxCount uint64
	)
	if err = response.GetObject(&txPoolInspect); err != nil {
		return nil, err
	}
	pending := reflect.ValueOf(txPoolInspect.Pending)
	if pending.Kind() == reflect.Map {
		for _, key := range pending.MapKeys() {
			address := key.Interface().(string)
			tx := reflect.ValueOf(pending.MapIndex(key).Interface())
			if tx.Kind() == reflect.Map && strings.ToLower(targetAddr.String()) == strings.ToLower(address) {
				for _, key := range tx.MapKeys() {
					count := key.Interface().(uint64)
					if count > txPoolMaxCount {
						txPoolMaxCount = count
					}
				}
			}
		}
	}
	pendingNonce := pendingNonceAt
	if pendingNonceAt != 0 && txPoolMaxCount+1 > pendingNonceAt {
		pendingNonce = txPoolMaxCount + 1
	}

	return big.NewInt(int64(pendingNonce)), nil
}

// AccountBalance returns the account balance for a given common.
func (client *Client) AccountBalance(targetAddr common.Address) (*big.Int, error) {
	balance, err := client.EthClient.BalanceAt(client.Ctx, targetAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance for '%v': %v", targetAddr, err)
	}

	return balance, nil
}

package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/multitenancy"
	usecases "github.com/lugondev/tx-builder/src/api/business/use-cases"
	"github.com/lugondev/tx-builder/src/api/service/formatters"
	api "github.com/lugondev/tx-builder/src/api/service/types"
	infra "github.com/lugondev/tx-builder/src/infra/api"
	"github.com/lugondev/wallet-signer-manager/pkg/client"
	qkmstoretypes "github.com/lugondev/wallet-signer-manager/src/stores/api/types"
)

type AccountsController struct {
	ucs              usecases.AccountUseCases
	keyManagerClient client.KeyManagerClient
	storeName        string
}

func NewAccountsController(accountUCs usecases.AccountUseCases, keyManagerClient client.KeyManagerClient, qkmStoreID string) *AccountsController {
	return &AccountsController{
		accountUCs,
		keyManagerClient,
		qkmStoreID,
	}
}

// Append Add routes to router
func (c *AccountsController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/accounts").HandlerFunc(c.search)
	router.Methods(http.MethodPost).Path("/accounts").HandlerFunc(c.create)
	router.Methods(http.MethodPost).Path("/accounts/import").HandlerFunc(c.importKey)
	router.Methods(http.MethodGet).Path("/accounts/{address}").HandlerFunc(c.getOne)
	router.Methods(http.MethodPatch).Path("/accounts/{address}").HandlerFunc(c.update)
	router.Methods(http.MethodPost).Path("/accounts/{address}/sign-message").HandlerFunc(c.signMessage)
	//router.Methods(http.MethodPost).Path("/accounts/{address}/sign-typed-data").HandlerFunc(c.signTypedData)
	//router.Methods(http.MethodPost).Path("/accounts/verify-message").HandlerFunc(c.verifyMessageSignature)
	//router.Methods(http.MethodPost).Path("/accounts/verify-typed-data").HandlerFunc(c.verifyTypedDataSignature)
}

// @Summary      Creates a new Account
// @Description  Creates a new Account
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     JWTAuth
// @Param        request  body      api.CreateAccountRequest  true  "Account creation request"
// @Success      200      {object}  api.AccountResponse       "Account object"
// @Failure      400      {object}  infra.ErrorResponse    "Invalid request"
// @Failure      401      {object}  infra.ErrorResponse    "Unauthorized"
// @Failure      500      {object}  infra.ErrorResponse    "Internal server error"
// @Router       /accounts [post]
func (c *AccountsController) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req := &api.CreateAccountRequest{}
	err := infra.UnmarshalBody(request.Body, req)
	if err != nil {
		infra.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := c.ucs.Create().Execute(ctx, formatters.FormatCreateAccountRequest(req, c.storeName), nil, req.Chain,
		multitenancy.UserInfoValue(ctx))
	if err != nil {
		infra.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(acc))
}

// @Summary      Fetch a account by address
// @Description  Fetch a single account by address
// @Tags         Accounts
// @Produce      json
// @Security     ApiKeyAuth
// @Security     JWTAuth
// @Param        address  path      string                  true  "selected account address"
// @Success      200      {object}  api.AccountResponse     "Account found"
// @Failure      404      {object}  infra.ErrorResponse  "Account not found"
// @Failure      401      {object}  infra.ErrorResponse  "Unauthorized"
// @Failure      500      {object}  infra.ErrorResponse  "Internal server error"
// @Router       /accounts/{address} [get]
func (c *AccountsController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	pubkey := mux.Vars(request)["pubkey"] //utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	//if err != nil {
	//	infra.WriteError(rw, err.Error(), http.StatusBadRequest)
	//	return
	//}

	acc, err := c.ucs.Get().Execute(ctx, pubkey, multitenancy.UserInfoValue(ctx))
	if err != nil {
		infra.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(acc))
}

// @Summary      Search accounts by provided filters
// @Description  Get a list of filtered accounts
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     JWTAuth
// @Param        aliases  query     []string                false  "List of account aliases"  collectionFormat(csv)
// @Success      200      {array}   api.AccountResponse     "List of identities found"
// @Failure      400      {object}  infra.ErrorResponse  "Invalid filter in the request"
// @Failure      401      {object}  infra.ErrorResponse  "Unauthorized"
// @Failure      500      {object}  infra.ErrorResponse  "Internal server error"
// @Router       /accounts [get]
func (c *AccountsController) search(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	filters, err := formatters.FormatAccountFilterRequest(request)
	if err != nil {
		infra.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	accs, err := c.ucs.Search().Execute(ctx, filters, multitenancy.UserInfoValue(ctx))
	if err != nil {
		infra.WriteHTTPErrorResponse(rw, err)
		return
	}

	response := []*api.AccountResponse{}
	for _, acc := range accs {
		response = append(response, formatters.FormatAccountResponse(acc))
	}

	_ = json.NewEncoder(rw).Encode(response)
}

// @Summary      Creates a new Account by importing a private key
// @Description  Creates a new Account by importing a private key
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     JWTAuth
// @Param        request  body      api.ImportAccountRequest  true  "Account creation request"
// @Success      200      {object}  api.AccountResponse       "Account object"
// @Failure      400      {object}  infra.ErrorResponse    "Invalid request"
// @Failure      422      {object}  infra.ErrorResponse    "Unprocessable entity"
// @Failure      401      {object}  infra.ErrorResponse    "Unauthorized"
// @Failure      405      {object}  infra.ErrorResponse    "Not allowed"
// @Failure      500      {object}  infra.ErrorResponse    "Internal server error"
// @Router       /accounts/import [post]
func (c *AccountsController) importKey(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req := &api.ImportAccountRequest{}
	err := infra.UnmarshalBody(request.Body, req)
	if err != nil {
		infra.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := c.ucs.Create().Execute(ctx, formatters.FormatImportAccountRequest(req, c.storeName), req.PrivateKey, req.Chain,
		multitenancy.UserInfoValue(ctx))
	if err != nil {
		infra.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(acc))
}

// @Summary      Update account by Address
// @Description  Update a specific account by Address
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     JWTAuth
// @Param        request  body      api.UpdateAccountRequest  true  "Account update request"
// @Param        address  path      string                    true  "selected account address"
// @Success      200      {object}  api.AccountResponse       "Account found"
// @Failure      400      {object}  infra.ErrorResponse    "Invalid request"
// @Failure      401      {object}  infra.ErrorResponse    "Unauthorized"
// @Failure      404      {object}  infra.ErrorResponse    "Account not found"
// @Failure      500      {object}  infra.ErrorResponse    "Internal server error"
// @Router       /accounts/{address} [patch]
func (c *AccountsController) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	accRequest := &api.UpdateAccountRequest{}
	err := infra.UnmarshalBody(request.Body, accRequest)
	if err != nil {
		infra.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	acc := formatters.FormatUpdateAccountRequest(accRequest)
	pubkey := mux.Vars(request)["pubkey"] //utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	//if err != nil {
	//	infra.WriteError(rw, err.Error(), http.StatusBadRequest)
	//	return
	//}
	acc.PublicKey = common.FromHex(pubkey)

	accRes, err := c.ucs.Update().Execute(ctx, acc, multitenancy.UserInfoValue(ctx))

	if err != nil {
		infra.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(accRes))
}

// @Summary      Sign Message (EIP-191)
// @Description  Sign message, following EIP-191, data using selected account
// @Tags         Accounts
// @Accept       json
// @Produce      text/plain
// @Security     ApiKeyAuth
// @Security     JWTAuth
// @Param        request  body      api.SignMessageRequest  true  "Payload to sign"
// @Param        address  path      string                  true  "selected account address"
// @Success      200      {string}  string                  "Signed payload"
// @Failure      400      {object}  infra.ErrorResponse  "Invalid request"
// @Failure      401      {object}  infra.ErrorResponse  "Unauthorized"
// @Failure      404      {object}  infra.ErrorResponse  "Account not found"
// @Failure      500      {object}  infra.ErrorResponse  "Internal server error"
// @Router       /accounts/{address}/sign-message [post]
func (c *AccountsController) signMessage(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	payloadRequest := &api.SignMessageRequest{}
	err := infra.UnmarshalBody(request.Body, payloadRequest)
	if err != nil {
		infra.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	pubkey := mux.Vars(request)["pubkey"] //utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	//if err != nil {
	//	infra.WriteError(rw, err.Error(), http.StatusBadRequest)
	//	return
	//}

	_, err = c.ucs.Get().Execute(ctx, pubkey, multitenancy.UserInfoValue(ctx))
	if err != nil {
		infra.WriteError(rw, fmt.Sprintf("pubkey %s was not found", pubkey), http.StatusBadRequest)
		return
	}

	qkmStoreID := payloadRequest.StoreID
	if qkmStoreID == "" {
		qkmStoreID = c.storeName
	}

	signature, err := c.keyManagerClient.Sign(request.Context(), qkmStoreID, pubkey, &qkmstoretypes.SignWalletRequest{
		Data: payloadRequest.Data,
	})
	if err != nil {
		infra.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}

// @Summary      Signs typed data using an existing account following the EIP-712 standard
// @Description  Signs typed data using ECDSA and the private key of an existing account following the EIP-712 standard
// @Tags         Accounts
// @Accept       json
// @Produce      text/plain
// @Param        request  body      api.SignTypedDataRequest  true  "Typed data to sign"
// @Param        address  path      string                    true  "selected account address"
// @Success      200      {string}  string                    "Signed payload"
// @Failure      400      {object}  infra.ErrorResponse    "Invalid request"
// @Failure      401      {object}  infra.ErrorResponse    "Unauthorized"
// @Failure      404      {object}  infra.ErrorResponse    "Account not found"
// @Failure      422      {object}  infra.ErrorResponse    "Invalid parameters"
// @Failure      500      {object}  infra.ErrorResponse    "Internal server error"
// @Router       /accounts/{address}/sign-typed-data [post]
//func (c *AccountsController) signTypedData(rw http.ResponseWriter, request *http.Request) {
//	ctx := request.Context()
//	signRequest := &api.SignTypedDataRequest{}
//	err := infra.UnmarshalBody(request.Body, signRequest)
//	if err != nil {
//		infra.WriteError(rw, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
//	if err != nil {
//		infra.WriteError(rw, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	_, err = c.ucs.Get().Execute(ctx, *address, multitenancy.UserInfoValue(ctx))
//	if err != nil {
//		infra.WriteError(rw, fmt.Sprintf("account %s was not found", address), http.StatusBadRequest)
//		return
//	}
//
//	qkmStoreID := signRequest.StoreID
//	if qkmStoreID == "" {
//		qkmStoreID = c.storeName
//	}
//
//	signature, err := c.keyManagerClient.SignTypedData(ctx, qkmStoreID, address.Hex(), &qkmstoretypes.SignTypedDataRequest{
//		DomainSeparator: signRequest.DomainSeparator,
//		Types:           signRequest.Types,
//		Message:         signRequest.Message,
//		MessageType:     signRequest.MessageType,
//	})
//	if err != nil {
//		infra.WriteHTTPErrorResponse(rw, err)
//		return
//	}
//
//	_, _ = rw.Write([]byte(signature))
//}

// @Summary      Verifies the signature of a typed data message following the EIP-712 standard
// @Description  Verifies if a typed data message has been signed by the Ethereum account passed as argument following the EIP-712 standard
// @Tags         Accounts
// @Accept       json
// @Param        request  body  qkmutilstypes.VerifyTypedDataRequest  true  "Typed data to sign"
// @Success      204
// @Failure      400  {object}  infra.ErrorResponse  "Invalid request"
// @Failure      401  {object}  infra.ErrorResponse  "Unauthorized"
// @Failure      422  {object}  infra.ErrorResponse  "Invalid parameters"
// @Failure      500  {object}  infra.ErrorResponse  "Internal server error"
// @Router       /accounts/verify-typed-data [post]
//func (c *AccountsController) verifyTypedDataSignature(rw http.ResponseWriter, request *http.Request) {
//	verifyRequest := &qkmutilstypes.VerifyTypedDataRequest{}
//	err := infra.UnmarshalBody(request.Body, verifyRequest)
//	if err != nil {
//		infra.WriteError(rw, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	err = c.keyManagerClient.VerifyTypedData(request.Context(), verifyRequest)
//	if err != nil {
//		infra.WriteHTTPErrorResponse(rw, err)
//		return
//	}
//
//	rw.WriteHeader(http.StatusNoContent)
//}

// @Summary      Verifies the signature of a message (EIP-191)
// @Description  Verifies if a message has been signed by the Ethereum account passed as argument
// @Tags         Accounts
// @Accept       json
// @Param        request  body  qkmutilstypes.VerifyRequest  true  "signature and message to verify"
// @Success      204
// @Failure      400  {object}  infra.ErrorResponse  "Invalid request"
// @Failure      401  {object}  infra.ErrorResponse  "Unauthorized"
// @Failure      422  {object}  infra.ErrorResponse  "Invalid parameters"
// @Failure      500  {object}  infra.ErrorResponse  "Internal server error"
// @Router       /accounts/verify-message [post]
//func (c *AccountsController) verifyMessageSignature(rw http.ResponseWriter, request *http.Request) {
//	verifyRequest := &qkmutilstypes.VerifyRequest{}
//	err := infra.UnmarshalBody(request.Body, verifyRequest)
//	if err != nil {
//		infra.WriteError(rw, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	err = c.keyManagerClient.VerifyMessage(request.Context(), verifyRequest)
//	if err != nil {
//		infra.WriteHTTPErrorResponse(rw, err)
//		return
//	}
//
//	rw.WriteHeader(http.StatusNoContent)
//}

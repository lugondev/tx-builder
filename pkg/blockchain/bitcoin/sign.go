package bitcoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/pkg/xcrypto"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
)

func Sign(dataHash []byte, typeSign string) ([]byte, error) {

	url := "http://localhost:8200/v1/quorum/wallets/0x02f564c5d9f932acbb0c81438f0e4389509f87383e22d4f203e0bb09c33135e86a/sign"

	payload := strings.NewReader(fmt.Sprintf("{\"data\": \"%s\",\"type_sign\" :\"%s\"}", hexutil.Encode(dataHash), typeSign))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("X-Vault-Token", "s.rSvcGSWGZ3y3uqdGDzLMU4OA")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Vault-Namespace", "lugon-test")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
	var rData map[string]interface{}
	err := json.Unmarshal(body, &rData)
	if err != nil {
		return nil, err
	}
	data := rData["data"].(map[string]interface{})
	signature := data["signature"].(string)

	return common.FromHex(signature), nil
}

func MpcSign(dataHash []byte) ([]byte, error) {
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
	if bytes.Compare(sharepub1.Serialize(), sharepub2.Serialize()) != 0 {
		return nil, fmt.Errorf("sharepub1 != sharepub2")
	}

	// Phase 2.
	encpk1, encpub1, scalarR1 := party1.Phase2(dataHash)
	encpk2, encpub2, scalarR2 := party2.Phase2(dataHash)

	// Phase 3.
	shareR1 := party1.Phase3(scalarR2)
	shareR2 := party2.Phase3(scalarR1)
	if bytes.Compare(shareR1.X.Bytes(), shareR2.X.Bytes()) != 0 || bytes.Compare(shareR1.Y.Bytes(), shareR2.Y.Bytes()) != 0 {
		return nil, fmt.Errorf("shareR1 != shareR2")
	}

	// Phase 4.
	sig1, err := party1.Phase4(encpk2, encpub2, shareR1)
	if err != nil {
		return nil, err
	}
	sig2, err := party2.Phase4(encpk1, encpub1, shareR2)
	if err != nil {
		return nil, err
	}

	// Phase 5.
	fs1, err := party1.Phase5(shareR1, sig2)
	if err != nil {
		return nil, err
	}
	fs2, err := party2.Phase5(shareR2, sig1)
	if err != nil {
		return nil, err
	}
	if bytes.Compare(fs1, fs2) != 0 {
		return nil, fmt.Errorf("fs1 != fs2")
	}

	// Verify.
	err = xcrypto.EcdsaVerify(sharepub1, dataHash, fs1)

	return fs1, err
}

func MpcSchnorrSign(dataHash []byte) ([]byte, error) {
	//// Party 1.
	//p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	//prv1 := xcrypto.PrvKeyFromBytes(p1.Bytes())
	//pub1 := prv1.PubKey()
	//party1, err := xcrypto.NewSchnorrParty(prv1)
	//if err != nil {
	//	return nil, err
	//}
	//defer party1.Close()
	//
	//// Party 2.
	//p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	//prv2 := xcrypto.PrvKeyFromBytes(p2.Bytes())
	//pub2 := prv2.PubKey()
	//party2, err := xcrypto.NewSchnorrParty(prv2)
	//if err != nil {
	//	return nil, err
	//}
	//defer party2.Close()
	//
	//// Phase 1.
	//sharepub1 := party1.Phase1(pub2)
	//sharepub2 := party2.Phase1(pub1)
	//if bytes.Compare(sharepub1.Serialize(), sharepub2.Serialize()) != 0 {
	//	return nil, fmt.Errorf("sharepub1 != sharepub2")
	//}
	//
	//// Phase 2.
	//r1 := party1.Phase2(dataHash)
	//r2 := party2.Phase2(dataHash)
	//
	//// Phase 3.
	//sharer1 := party1.Phase3(r2)
	//sharer2 := party2.Phase3(r1)
	//if bytes.Compare(sharer1.X.Bytes(), sharer2.X.Bytes()) != 0 || bytes.Compare(sharer1.Y.Bytes(), sharer2.Y.Bytes()) != 0 {
	//	return nil, fmt.Errorf("shareR1 != shareR2")
	//}
	//
	//// Phase 4.
	//s1, err := party1.Phase4(sharepub1, sharer1)
	//if err != nil {
	//	return nil, err
	//}
	//s2, err := party2.Phase4(sharepub2, sharer2)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Phase 5.
	//fs1, err := party1.Phase5(sharer1, s1, s2)
	//if err != nil {
	//	return nil, err
	//}
	//fs2, err := party2.Phase5(sharer2, s1, s2)
	//if err != nil {
	//	return nil, err
	//}
	//if bytes.Compare(fs1, fs2) != 0 {
	//	return nil, fmt.Errorf("fs1 != fs2")
	//}
	//
	//return fs1, err

	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := xcrypto.PrvKeyFromBytes(p1.Bytes())

	pub1 := prv1.PubKey()
	fmt.Printf("pub: %x\n", pub1.Serialize())
	return xcrypto.SchnorrSign(prv1, dataHash)
}

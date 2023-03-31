package bitcoin

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"io/ioutil"
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

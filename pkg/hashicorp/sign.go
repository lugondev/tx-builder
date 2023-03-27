package hashicorp

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lugondev/tx-builder/pkg/hashicorp/client"
	"io/ioutil"
	"net/http"
	"strings"
)

func SignTest(data []byte) ([]byte, error) {
	cli, err := client.NewClient(&client.Config{
		MountPoint:    "quorum",
		Address:       "http://localhost:8200",
		CACert:        "",
		CAPath:        "",
		ClientCert:    "",
		ClientKey:     "",
		TLSServerName: "",
		Namespace:     "lugon-test",
		ClientTimeout: 0,
		RateLimit:     0,
		BurstLimit:    0,
		MaxRetries:    0,
		SkipVerify:    true,
	})
	if err != nil {
		return nil, err
	}
	cli.SetToken("s.LtsIPANb31NG3w85M0Nt5KcU")
	sign, err := cli.SignWallet("0x0304df1e31533b96e542cf4565846b5cfff68dd3b3eb5e3639f2f59c35c59dbe7a", data)
	if err != nil {
		return nil, err
	}
	sig := sign.Data["signature"]
	return common.FromHex(sig.(string)), nil
}

func SignByKeyManager(data []byte) ([]byte, error) {

	url := "http://0.0.0.0:8080/stores/wallet-signer/wallets/0304df1e31533b96e542cf4565846b5cfff68dd3b3eb5e3639f2f59c35c59dbe7a/sign"

	fmt.Println("data: ", common.Bytes2Hex(data))
	payload := strings.NewReader(fmt.Sprintf("{\"data\":\"0x%s\"}", common.Bytes2Hex(data)))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	return common.FromHex(string(body)), nil
}

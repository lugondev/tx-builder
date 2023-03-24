package hashicorp

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/lugondev/tx-builder/pkg/hashicorp/client"
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
	cli.SetToken("DevVaultToken")
	sign, err := cli.SignWallet("0x02555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f", data)
	if err != nil {
		return nil, err
	}
	sig := sign.Data["signature"]
	return common.FromHex(sig.(string)), nil
}

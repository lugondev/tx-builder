package hashicorp

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/lugondev/tx-builder/pkg/hashicorp/client"
	"testing"
)

func TestConnectVaultAndImport(t *testing.T) {
	cli, err := client.NewClient(&client.Config{
		MountPoint:    "quorum",
		Address:       "http://localhost:8200",
		CACert:        "",
		CAPath:        "",
		ClientCert:    "",
		ClientKey:     "",
		TLSServerName: "",
		Namespace:     "test-tx-builder",
		ClientTimeout: 0,
		RateLimit:     0,
		BurstLimit:    0,
		MaxRetries:    0,
		SkipVerify:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	cli.SetToken("s.MxMu6BHP1mReMI6xLPDu4BhL")
	importedKey, err := cli.ImportKey(map[string]interface{}{
		"private_key":       "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw=",
		"curve":             "secp256k1",
		"signing_algorithm": "ecdsa",
		"id":                "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(importedKey.Data)
}

func TestConnectVaultAndSign(t *testing.T) {
	cli, err := client.NewClient(&client.Config{
		MountPoint:    "quorum",
		Address:       "http://localhost:8200",
		CACert:        "",
		CAPath:        "",
		ClientCert:    "",
		ClientKey:     "",
		TLSServerName: "",
		Namespace:     "test-tx-builder",
		ClientTimeout: 0,
		RateLimit:     0,
		BurstLimit:    0,
		MaxRetries:    0,
		SkipVerify:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	cli.SetToken("s.MxMu6BHP1mReMI6xLPDu4BhL")
	importedKey, err := cli.Sign("1", common.FromHex("94918ae4ded1d00ee7ec34901f199f6e9d37c443b2bdd85f5b412309c9553c54"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(importedKey.Data)
}

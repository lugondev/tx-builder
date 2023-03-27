package builder

import (
	"fmt"
	"github.com/lugondev/tx-builder/pkg/common"
	"testing"
)

func TestAddressInfo(t *testing.T) {
	info := common.GetBTCAddressInfo("mkF4Rkh9bQoUujuk5zJnvcamXvTpUgSNss")
	fmt.Println("address:", info.Address)
	fmt.Println("type:", info.Type)
	fmt.Println("version:", info.Version)
}

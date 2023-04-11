package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestSign(t *testing.T) {

	signedTaproot, err := Sign(common.FromHex("0xdb337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"), "taproot")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("signed taproot: ", hex.EncodeToString(signedTaproot))
	signEcdsa, err := Sign(common.FromHex("0xdb337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"), "ecdsa")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("signed ecdsa: ", hex.EncodeToString(signEcdsa))

	signEth, err := Sign(common.FromHex("0xdb337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"), "eth_sign")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("signed eth: ", hex.EncodeToString(signEth))

}

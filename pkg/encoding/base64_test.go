package encoding

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"
)

const (
	encodedData = "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw="
	decodedHex  = "0xdb337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"
	//encodedData = "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI="
	//decodedHex  = "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2"
)

func TestEncodeToBase64(t *testing.T) {
	encoded := EncodeToBase64(common.FromHex(decodedHex))
	if encoded != encodedData {
		t.Errorf("Expected %s, got %s", encodedData, encoded)
	}
}

func TestDecodeBase64(t *testing.T) {
	decoded, err := DecodeBase64(encodedData)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !bytes.Equal(decoded, common.FromHex(decodedHex)) {
		t.Errorf("Expected %s, got %s", decodedHex, hexutil.Encode(decoded))
	}
}

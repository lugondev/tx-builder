package utils

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func IsHexString(s string) bool {
	if !strings.HasPrefix(s, "0x") {
		s = fmt.Sprintf("0x%s", s)
	}
	_, err := hexutil.Decode(s)
	return err == nil
}

func StringToHexBytes(v string) hexutil.Bytes {
	if v == "" {
		return nil
	}

	if vb, err := hexutil.Decode(v); err == nil {
		return vb
	}

	return nil
}

func HexBytesToString(hex hexutil.Bytes) string {
	return hexutil.Encode(hex)
}

func StringBigIntToHex(v string) *hexutil.Big {
	if v == "" {
		return nil
	}

	if bv, ok := new(big.Int).SetString(v, 10); ok {
		return (*hexutil.Big)(bv)
	}

	return nil
}

func HexToBigInt(v string) *hexutil.Big {
	if v == "" {
		return nil
	}

	if bv, err := hexutil.DecodeBig(v); err == nil {
		return (*hexutil.Big)(bv)
	}

	return nil
}

func HexToBigIntString(v *hexutil.Big) string {
	if v == nil {
		return ""
	}

	return v.ToInt().String()
}

func Uint64ToHex(v uint64) *hexutil.Uint64 {
	return ToPtr(hexutil.Uint64(v)).(*hexutil.Uint64)
}

func UintToHex(v uint) *hexutil.Uint {
	return ToPtr(hexutil.Uint(v)).(*hexutil.Uint)
}

// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package schnorr

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"testing"

	"crypto/ecdsa"

	"github.com/lugondev/tx-builder/pkg/xcrypto/secp256k1"
	"github.com/stretchr/testify/assert"
)

var (
	schnorrTests = []struct {
		d      string
		pk     string
		m      string
		sig    string
		err    error
		desc   string
		result bool
	}{
		{
			d:      "0000000000000000000000000000000000000000000000000000000000000001",
			pk:     "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
			m:      "0000000000000000000000000000000000000000000000000000000000000000",
			sig:    "787A848E71043D280C50470E8E1532B2DD5D20EE912A45DBDD2BD1DFBF187EF67031A98831859DC34DFFEEDDA86831842CCD0079E1F92AF177F7F22CC1DCED05",
			result: true,
		},

		{
			d:      "B7E151628AED2A6ABF7158809CF4F3C762E7160F38B4DA56A784D9045190CFEF",
			pk:     "02DFF1D77F2A671C5F36183726DB2341BE58FEAE1DA2DECED843240F7B502BA659",
			m:      "243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89",
			sig:    "2A298DACAE57395A15D0795DDBFD1DCB564DA82B0F269BC70A74F8220429BA1D1E51A22CCEC35599B8F266912281F8365FFC2D035A230434A1A64DC59F7013FD",
			result: true,
		},
		{
			d:      "C90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B14E5C7",
			pk:     "03FAC2114C2FBB091527EB7C64ECB11F8021CB45E8E7809D3C0938E4B8C0E5F84B",
			m:      "5E2D58D8B3BCDF1ABADEC7829054F90DDA9805AAB56C77333024B9D0A508B75C",
			sig:    "00DA9B08172A9B6F0466A2DEFD817F2D7AB437E0D253CB5395A963866B3574BE00880371D01766935B92D2AB4CD5C8A2A5837EC57FED7660773A05F0DE142380",
			result: true,
		},
		{
			d:      "6D6C66873739BC7BFB3526629670D0EA357E92CC4581490D62779AE15F6B787B",
			pk:     "026D7F1D87AB3BBC8BC01F95D9AECE1E659D6E33C880F8EFA65FACF83E698BBBF7",
			m:      "B2F0CD8ECB23C1710903F872C31B0FD37E15224AF457722A87C5E0C7F50FFFB3",
			sig:    "68CA1CC46F291A385E7C255562068357F964532300BEADFFB72DD93668C0C1CAC8D26132EB3200B86D66DE9C661A464C6B2293BB9A9F5B966E53CA736C7E504F",
			result: true,
		},
		{
			d:      "",
			pk:     "03DEFDEA4CDB677750A420FEE807EACF21EB9898AE79B9768766E4FAA04A2D4A34",
			m:      "4DF3C3F68FCC83B27E9D42C90431A72499F17875C81A599B566C9889B9696703",
			sig:    "00000000000000000000003B78CE563F89A0ED9414F5AA28AD0D96D6795F9C6302A8DC32E64E86A333F20EF56EAC9BA30B7246D6D25E22ADB8C6BE1AEB08D49D",
			result: true,
		},
		{
			d:      "",
			pk:     "031B84C5567B126440995D3ED5AABA0565D71E1834604819FF9C17F5E9D5DD078F",
			m:      "0000000000000000000000000000000000000000000000000000000000000000",
			sig:    "52818579ACA59767E3291D91B76B637BEF062083284992F2D95F564CA6CB4E3530B1DA849C8E8304ADC0CFE870660334B3CFC18E825EF1DB34CFAE3DFC5D8187",
			result: true,
			desc:   "test fails if jacobi symbol of x(R) instead of y(R) is used",
		},
		{
			d:      "",
			pk:     "03FAC2114C2FBB091527EB7C64ECB11F8021CB45E8E7809D3C0938E4B8C0E5F84B",
			m:      "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			sig:    "570DD4CA83D4E6317B8EE6BAE83467A1BF419D0767122DE409394414B05080DCE9EE5F237CBD108EABAE1E37759AE47F8E4203DA3532EB28DB860F33D62D49BD",
			result: true,
			desc:   "test fails if msg is reduced",
		},
		{
			d:      "",
			pk:     "02DFF1D77F2A671C5F36183726DB2341BE58FEAE1DA2DECED843240F7B502BA659",
			m:      "243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89",
			sig:    "2A298DACAE57395A15D0795DDBFD1DCB564DA82B0F269BC70A74F8220429BA1DFA16AEE06609280A19B67A24E1977E4697712B5FD2943914ECD5F730901B4AB7",
			result: false,
			err:    errors.New("signature verification failed"),
			desc:   "incorrect R residuosity",
		},
		{
			d:      "",
			pk:     "03FAC2114C2FBB091527EB7C64ECB11F8021CB45E8E7809D3C0938E4B8C0E5F84B",
			m:      "5E2D58D8B3BCDF1ABADEC7829054F90DDA9805AAB56C77333024B9D0A508B75C",
			sig:    "00DA9B08172A9B6F0466A2DEFD817F2D7AB437E0D253CB5395A963866B3574BED092F9D860F1776A1F7412AD8A1EB50DACCC222BC8C0E26B2056DF2F273EFDEC",
			result: false,
			err:    errors.New("signature verification failed"),
			desc:   "negated message hash",
		},
		{
			d:      "",
			pk:     "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
			m:      "0000000000000000000000000000000000000000000000000000000000000000",
			sig:    "787A848E71043D280C50470E8E1532B2DD5D20EE912A45DBDD2BD1DFBF187EF68FCE5677CE7A623CB20011225797CE7A8DE1DC6CCD4F754A47DA6C600E59543C",
			result: false,
			err:    errors.New("signature verification failed"),
			desc:   "negated s value",
		},
		{
			d:      "",
			pk:     "03DFF1D77F2A671C5F36183726DB2341BE58FEAE1DA2DECED843240F7B502BA659",
			m:      "243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89",
			sig:    "2A298DACAE57395A15D0795DDBFD1DCB564DA82B0F269BC70A74F8220429BA1D1E51A22CCEC35599B8F266912281F8365FFC2D035A230434A1A64DC59F7013FD",
			result: false,
			err:    errors.New("signature verification failed"),
			desc:   "negated public key",
		},
		{
			d:      "",
			pk:     "02DFF1D77F2A671C5F36183726DB2341BE58FEAE1DA2DECED843240F7B502BA659",
			m:      "243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89",
			sig:    "00000000000000000000000000000000000000000000000000000000000000009E9D01AF988B5CEDCE47221BFA9B222721F3FA408915444A4B489021DB55775F",
			result: false,
			err:    errors.New("signature verification failed"),
			desc:   "sG - eP is infinite. Test fails in single verification if jacobi(y(inf)) is defined as 1 and x(inf) as 0",
		},
		{
			d:      "",
			pk:     "02DFF1D77F2A671C5F36183726DB2341BE58FEAE1DA2DECED843240F7B502BA659",
			m:      "243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89",
			sig:    "0000000000000000000000000000000000000000000000000000000000000001D37DDF0254351836D84B1BD6A795FD5D523048F298C4214D187FE4892947F728",
			result: false,
			err:    errors.New("signature verification failed"),
			desc:   "sG - eP is infinite. Test fails in single verification if jacobi(y(inf)) is defined as 1 and x(inf) as 1",
		},
		{
			d:      "",
			pk:     "02DFF1D77F2A671C5F36183726DB2341BE58FEAE1DA2DECED843240F7B502BA659",
			m:      "243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89",
			sig:    "4A298DACAE57395A15D0795DDBFD1DCB564DA82B0F269BC70A74F8220429BA1D1E51A22CCEC35599B8F266912281F8365FFC2D035A230434A1A64DC59F7013FD",
			result: false,
			err:    errors.New("signature verification failed"),
			desc:   "sig[0:32] is not an X coordinate on the curve",
		},
		{
			d:      "",
			pk:     "02DFF1D77F2A671C5F36183726DB2341BE58FEAE1DA2DECED843240F7B502BA659",
			m:      "243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89",
			sig:    "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC2F1E51A22CCEC35599B8F266912281F8365FFC2D035A230434A1A64DC59F7013FD",
			result: false,
			err:    errors.New("r is larger than or equal to field size"),
			desc:   "sig[0:32] is equal to field size",
		},
	}
)

func TestSchnorrSign(t *testing.T) {
	curve := secp256k1.SECP256K1()
	for _, test := range schnorrTests {
		if test.d == "" {
			continue
		}
		d, _ := new(big.Int).SetString(test.d, 16)
		prv := &ecdsa.PrivateKey{
			D: d,
		}
		prv.Curve = curve

		mbytes, _ := hex.DecodeString(test.m)
		r, s, err := Sign(prv, mbytes)
		assert.Nil(t, err)

		var sig [64]byte
		copy(sig[:32], IntToByte(r))
		copy(sig[32:], IntToByte(s))
		got := hex.EncodeToString(sig[:])
		want := strings.ToLower(test.sig)
		assert.Equal(t, want, got)
	}
}

func TestSchnorrVerify(t *testing.T) {
	for _, test := range schnorrTests {
		curve := secp256k1.SECP256K1()
		pkbytes, _ := hex.DecodeString(test.pk)
		x, y := secp256k1.SecUnmarshal(curve, pkbytes)
		pub := &ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		}

		mbytes, _ := hex.DecodeString(test.m)
		sig, _ := hex.DecodeString(test.sig)
		r := new(big.Int).SetBytes(sig[:32])
		s := new(big.Int).SetBytes(sig[32:])
		got := Verify(pub, mbytes, r, s)
		want := test.result
		assert.Equal(t, want, got)
	}
}

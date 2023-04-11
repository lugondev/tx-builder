// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcrypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivKeys(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{
			name: "check curve",
			key: []byte{
				0xea, 0xf0, 0x2c, 0xa3, 0x48, 0xc5, 0x24, 0xe6,
				0x39, 0x26, 0x55, 0xba, 0x4d, 0x29, 0x60, 0x3c,
				0xd1, 0xa7, 0x34, 0x7d, 0x9d, 0x65, 0xcf, 0xe9,
				0x3c, 0xe1, 0xeb, 0xff, 0xdc, 0xa2, 0x26, 0x94,
			},
		},
	}

	for _, test := range tests {
		priv := PrvKeyFromBytes(test.key)
		pub := priv.PubKey()

		_, err := PubKeyFromBytes(pub.SerializeUncompressed())
		if err != nil {
			t.Errorf("%s privkey: %v", test.name, err)
			continue
		}

		serializedKey := priv.Serialize()
		if !bytes.Equal(serializedKey, test.key) {
			t.Errorf("%s unexpected serialized bytes - got: %x, "+
				"want: %x", test.name, serializedKey, test.key)
		}
	}
}

func TestPrivKeyAdd(t *testing.T) {
	hex1 := []byte{0x01}
	prvkey1 := PrvKeyFromBytes(hex1)

	hex2 := []byte{0x02}
	prvkey2 := prvkey1.Add(hex2)

	hex3 := []byte{0x03}
	prvkey3 := PrvKeyFromBytes(hex3)
	assert.Equal(t, prvkey2, prvkey3)
}

package utils

import (
	"bytes"
	"github.com/btcsuite/btcd/btcec/v2"
	"testing"
)

func TestParseHexToPubkey(t *testing.T) {
	pubkeyCompressed := "0x02e47e16da1880463b011a77c3fa2b7bfec5a78cf426a4eb927623925f3710aaf2"
	pubkeyUnCompressed := "0x04e47e16da1880463b011a77c3fa2b7bfec5a78cf426a4eb927623925f3710aaf20b70b0ef1ec12434ff068618cea8e6f5b27bc86515b9c61435565c145eeb6ff6"

	testPubkey(t, pubkeyCompressed, pubkeyUnCompressed)
}

func testPubkey(t *testing.T, pubkeyCompressed, pubkeyUnCompressed string) {
	if !IsHexString(pubkeyCompressed) || !IsHexString(pubkeyUnCompressed) {
		t.Fatal("expected hex string")
	}
	bytesPubkeyCompressed := StringToHexBytes(pubkeyCompressed)
	bytesPubkeyUnCompressed := StringToHexBytes(pubkeyUnCompressed)

	publicKeyCompressed, err := btcec.ParsePubKey(bytesPubkeyCompressed)
	if err != nil {
		t.Fatal(err)
	}

	publicKeyUnCompressed, err := btcec.ParsePubKey(bytesPubkeyUnCompressed)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(publicKeyUnCompressed.SerializeCompressed(), publicKeyCompressed.SerializeCompressed()) != 0 {
		t.Fatal("expected equal")
	}

	if bytes.Compare(publicKeyCompressed.SerializeUncompressed(), publicKeyUnCompressed.SerializeUncompressed()) != 0 {
		t.Fatal("expected equal")
	}

	if bytes.Compare(publicKeyCompressed.SerializeUncompressed(), bytesPubkeyUnCompressed) != 0 {
		t.Fatal("expected equal")
	}

	if bytes.Compare(publicKeyCompressed.SerializeCompressed(), bytesPubkeyCompressed) != 0 {
		t.Fatal("expected equal")
	}

	if bytes.Compare(publicKeyUnCompressed.SerializeUncompressed(), bytesPubkeyUnCompressed) != 0 {
		t.Fatal("expected equal")
	}

	if bytes.Compare(publicKeyUnCompressed.SerializeCompressed(), bytesPubkeyCompressed) != 0 {
		t.Fatal("expected equal")
	}
}

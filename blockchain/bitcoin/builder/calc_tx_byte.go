package builder

import (
	"github.com/lugondev/tx-builder/pkg/common"
	"strings"
)

const P2pkhInSize = 148
const P2pkhOutSize = 34
const P2shOutSize = 32
const P2shP2wpkhOutSize = 32
const P2shP2wshOutSize = 32
const P2shP2wpkhInSize = 90.75
const P2wpkhInSize = 67.75
const P2wpkhOutSize = 31
const P2wshOutSize = 43
const P2trOutSize = 43
const P2trInSize = 57.25
const PubkeySize = 33
const SignatureSize = 72

func getSizeOfScriptLengthElement(length float64) float64 {
	if length < 75 {
		return 1
	} else if length <= 255 {
		return 2
	} else if length <= 65535 {
		return 3
	} else if length <= 4294967295 {
		return 5
	} else {
		panic("Invalid script length element")
	}
}

func getSizeOfVarInt(length float64) float64 {
	if length < 253 {
		return 1
	} else if length < 65535 {
		return 3
	} else if length < 4294967295 {
		return 5
	} else if length < 18446744073709551615 {
		return 9
	} else {
		panic("Invalid varint length")
	}
}

func getTxOverheadVBytes(inputScript string, inputCount, outputCount float64) float64 {
	witnessVbytes := float64(0)
	if inputScript == "P2PKH" || inputScript == "P2SH" {
		witnessVbytes = float64(0)
	} else { // Transactions with segwit inputs have extra overhead
		witnessVbytes = 0.25 + 0.25 + float64(inputCount)/4 // witness element count per input
	}

	return 4 + getSizeOfVarInt(inputCount) + getSizeOfVarInt(outputCount) + 4 + witnessVbytes
}

func getTxOverheadExtraRawBytes(inputScript string, inputCount float64) float64 {
	witnessBytes := float64(0)
	// Returns the remaining 3/4 bytes per witness bytes
	if inputScript == "P2PKH" || inputScript == "P2SH" {
		witnessBytes = float64(0)
	} else { // Transactions with segwit inputs have extra overhead
		witnessBytes = 0.25 + 0.25 + float64(inputCount)/4 // witness element count per input
	}

	return witnessBytes * 3
}

func CalculateTxBytes(fromAddress string, inputCount float64, addressOutputs []string) float64 {
	// Validate transaction input attributes
	inputScript := strings.ToUpper(common.GetBTCAddressInfo(fromAddress).Version)
	signaturesPerInput := float64(1)
	pubkeyPerInput := float64(1)

	// Validate transaction output attributes
	p2pkhOutputCount := float64(0)
	p2shOutputCount := float64(0)
	p2shP2wpkhOutputCount := float64(0)
	p2shP2wshOutputCount := float64(0)
	p2wpkhOutputCount := float64(0)
	p2wshOutputCount := float64(0)
	p2trOutputCount := float64(0)

	for i := range addressOutputs {
		if common.GetBTCAddressInfo(addressOutputs[i]).GetVersion() == "P2PKH" {
			p2pkhOutputCount++
		}
		if common.GetBTCAddressInfo(addressOutputs[i]).GetVersion() == "P2SH" {
			p2shOutputCount++
		}
		if common.GetBTCAddressInfo(addressOutputs[i]).GetVersion() == "P2SH-P2WPKH" {
			p2shP2wpkhOutputCount++
		}
		if common.GetBTCAddressInfo(addressOutputs[i]).GetVersion() == "P2SH-P2WSH" {
			p2shP2wshOutputCount++
		}
		if common.GetBTCAddressInfo(addressOutputs[i]).GetVersion() == "P2WPKH" {
			p2wpkhOutputCount++
		}
		if common.GetBTCAddressInfo(addressOutputs[i]).GetVersion() == "P2WSH" {
			p2wshOutputCount++
		}
		if common.GetBTCAddressInfo(addressOutputs[i]).GetVersion() == "P2TR" {
			p2trOutputCount++
		}
	}
	return calc(inputScript, inputCount, signaturesPerInput, pubkeyPerInput, p2pkhOutputCount, p2shOutputCount, p2shP2wpkhOutputCount, p2shP2wshOutputCount, p2wpkhOutputCount, p2wshOutputCount, p2trOutputCount)
}

func calc(
	inputScript string,
	inputCount, inputSignatureCount, inputPubkeyCount,
	p2pkhOutputCount,
	p2shOutputCount,
	p2shP2wpkhOutputCount,
	p2shP2wshOutputCount,
	p2wpkhOutputCount,
	p2wshOutputCount,
	p2trOutputCount float64,
) float64 {
	outputCount := p2pkhOutputCount + p2shOutputCount + p2shP2wpkhOutputCount + p2shP2wshOutputCount + p2wpkhOutputCount + p2wshOutputCount + p2trOutputCount

	// In most cases the input size is predictable. For multisig inputs we need to perform a detailed calculation
	inputSize := float64(0) // in virtual bytes
	inputWitnessSize := float64(0)
	switch inputScript {
	case "P2PKH":
		inputSize = P2pkhInSize
		break
	case "P2SH-P2WPKH":
		inputSize = P2shP2wpkhInSize
		inputWitnessSize = 107 // size(signature) + signature + size(pubkey) + pubkey
		break
	case "P2WPKH":
		inputSize = P2wpkhInSize
		inputWitnessSize = 107 // size(signature) + signature + size(pubkey) + pubkey
		break
	case "P2TR": // Only consider the cooperative taproot signing path; assume multisig is done via aggregate signatures
		inputSize = P2trInSize
		inputWitnessSize = 65 // getSizeOfVarInt(schnorrSignature) + schnorrSignature;
		break
	case "P2SH":
		redeemScriptSize := 1 + // OP_M
			inputPubkeyCount*(1+PubkeySize) + // OP_PUSH33 <pubkey>
			1 + // OP_N
			1 // OP_CHECKMULTISIG
		scriptSigSize := 1 + float64(inputSignatureCount*(1+SignatureSize)) + getSizeOfScriptLengthElement(redeemScriptSize) + redeemScriptSize
		inputSize = 32 + 4 + getSizeOfVarInt(scriptSigSize) + scriptSigSize + 4
		break
	case "P2SH-P2WSH":
	case "P2WSH":
		redeemScriptSize := 1 + // OP_M
			inputPubkeyCount*(1+PubkeySize) + // OP_PUSH33 <pubkey>
			1 + // OP_N
			1 // OP_CHECKMULTISIG
		inputWitnessSize = 1 + float64(inputSignatureCount*(1+SignatureSize)) + getSizeOfScriptLengthElement(redeemScriptSize) + redeemScriptSize
		inputSize = 36 + // outpoint (spent UTXO ID)
			inputWitnessSize/4 + // witness program
			4 // nSequence
		if inputScript == "P2SH-P2WSH" {
			inputSize += 32 + 3 // P2SH wrapper (redeemscript hash) + overhead?
		}
	}

	txVBytes := getTxOverheadVBytes(inputScript, inputCount, outputCount) +
		inputSize*inputCount +
		P2pkhOutSize*p2pkhOutputCount +
		P2shOutSize*p2shOutputCount +
		P2shP2wpkhOutSize*p2shP2wpkhOutputCount +
		P2shP2wshOutSize*p2shP2wshOutputCount +
		P2wpkhOutSize*p2wpkhOutputCount +
		P2wshOutSize*p2wshOutputCount +
		P2trOutSize*p2trOutputCount

	return getTxOverheadExtraRawBytes(inputScript, inputCount) + txVBytes + (inputWitnessSize*inputCount)*3/4
}

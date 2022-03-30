package deepcoverclient

// Go packages
import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Bytes2HexString returns a hexadecimal string representation
// of an input byte array
func Bytes2HexString(indata []byte) string {
	return hex.EncodeToString(indata)
}

// HexString2Bytes returns a byte array representation
// of an input hexadecimal string
func HexString2Bytes(indata string) []byte {
	src, err := hex.DecodeString(indata)
	if err != nil {
		fmt.Println("Conversion error: ", err.Error())
	}
	// return the bytes decoded before the error
	return src
}

// VerifyDeepCoverSignature virifies the signature using public keys and calculated digest
func VerifyDeepCoverSignature(pub *ecdsa.PublicKey, digest []byte, signatureR string, signatureS string) bool {
	var ecdsaSig ecdsaSignature

	r := new(big.Int)
	s := new(big.Int)

	r.SetString(signatureR, 16)
	s.SetString(signatureS, 16)

	ecdsaSig.R = r
	ecdsaSig.S = s

	fmt.Printf("Converted signature R and S\n -->\tR: %d,\n -->\tS: %d\n\n", ecdsaSig.R, ecdsaSig.S)
	return ecdsa.Verify(pub, digest, ecdsaSig.R, ecdsaSig.S)
}

// CalcucateMessageDigest calculates the message digest
func CalcucateMessageDigest(message []byte, romID []byte) []byte {
	/*
		this function requires folowing packages:

		"crypto/sha256"
		"encoding/hex"

	*/

	var totalBuffer []byte

	if len(romID) != 8 {
		// define default romID of length 8
		// this will, of course, result in the verification failure
		// TO DO: apply different error handling mechanism
		romID = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}

	// Reconstruct DeepCover data fields for message digest (SHA256) calculation
	/*
		-------------------------------------+----------
		ROMID=romID []byte    				 |	 8 bytes
		PAGE0=FFFF........................FF | 	32 bytes
		BUFFER=xxx........................xx | 	32 bytes <-hash of transaction bytes in blockchain
		PAGE=00 							 |	 1 byte
		MANID=0000 							 |	 2 bytes
		-------------------------------------+ total 75bytes
	*/

	// hardcoded values from DeepCover eeprom memory
	page0HexString := "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
	pageHexString := "00"
	mainIDHexString := "0000"

	page0, err := hex.DecodeString(page0HexString)
	printError(err)
	page, err := hex.DecodeString(pageHexString)
	printError(err)
	mainID, err := hex.DecodeString(mainIDHexString)
	printError(err)

	// totalBuffer = romID + page0 + message + page + mainID
	totalBuffer = append(totalBuffer, romID...)
	totalBuffer = append(totalBuffer, page0...)
	totalBuffer = append(totalBuffer, message...)
	totalBuffer = append(totalBuffer, page...)
	totalBuffer = append(totalBuffer, mainID...)

	h := sha256.New()
	h.Write(totalBuffer)
	var sha256cs = h.Sum(nil)

	return sha256cs
}

func printError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
	} else {
		//fmt.Println("No error")
	}
}

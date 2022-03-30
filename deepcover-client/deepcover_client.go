package deepcoverclient

/*
	to build project with this package that uses imported C functions, following flags need to be set when calling go build:
	CC=/path/to/C/cross/compiler/for/RPi/arm-linux-gnueabihf-gcc
	CGO_ENABLED=1
	GOARCH=arm
	GOARM=7  (6 -if compatibility with RPi_1 is required)
	GOOS=linux

	full command example
	CC=/opt/beyond/rpi-toolchain/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/bin/arm-linux-gnueabihf-gcc CGO_ENABLED=1 GOARCH=arm GOARM=7 GOOS=linux go build -v -x main.go
*/

/*
#cgo CFLAGS: -I/opt/beyond/rpi-toolchain/arm-bcm2708/arm-linux-gnueabihf/arm-linux-gnueabihf/sysroot/usr/include/
#cgo LDFLAGS: -L/home/developer/go/src/github.com/vincepg13/bp-sdk/beyond/vendor/github.com/vincepg13/bp-sdk/deepcover-client/dcdriver/ -ldcdriver

#include "dcdriver/dcdriver.h"
#include <stdlib.h>
*/
import "C"

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"unsafe"
)

// GoGetDcID gets the DeepCover ID (romID)
func GoGetDcID(pg int) []byte {
	src := C.GoBytes(unsafe.Pointer(C.getDeepCoverID(C.int(pg))), 8)
	fmt.Println(" DeepCoverID: ", Bytes2HexString(src))
	return src
}

// GetDcID gets the DeepCover ID (romID) with the index 0
func GetDcID() []byte {
	return GoGetDcID(0)
}

// GoGetPageData gets the specified by index page data
func GoGetPageData(pg int, skiphdr int) []byte {
	src := C.GoBytes(unsafe.Pointer(C.getPageData(C.int(pg), C.int(skiphdr))), 32)
	return src
}

// GetPubKeyAX gets the PubKeyAX from the DeepCover
func GetPubKeyAX() []byte {
	return GoGetPageData(16, 1)
}

// GetPubKeyAY gets the PubKeyAX from the DeepCover
func GetPubKeyAY() []byte {
	return GoGetPageData(17, 1)
}

// GetPubKeyA gets the PubKeyA as the combination of PubKeyAX&PubKeyAY
func GetPubKeyA() []byte {
	var totalPubKey []byte
	totalPubKey = append(totalPubKey, GetPubKeyAX()...)
	totalPubKey = append(totalPubKey, GetPubKeyAY()...)

	return totalPubKey
}

// GoCalcSignature calculates the signature of the input data using DeepCover chip
func GoCalcSignature(indata []byte) []byte {
	cdata := C.CBytes(indata)
	cptrInData := (*C.uchar)(cdata)
	//C.computeReadPageAuthentication(cptrInData, C.int(1))
	out := C.GoBytes(unsafe.Pointer(C.computeReadPageAuthentication(cptrInData, C.int(1))), 64)
	C.free(unsafe.Pointer(cdata))
	return out
}

// SignData calculates the signature of the input data
func SignData(indata []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(indata)
	sha256cs := h.Sum(nil)
	fmt.Println(" Compressed input data to SHA256: ", Bytes2HexString(sha256cs))
	retVal := GoCalcSignature(sha256cs)
	fmt.Println(" DeepCover signature: ", Bytes2HexString(retVal))
	return retVal, nil
}

// GetPubKey returns publick key
func GetPubKey() *ecdsa.PublicKey {
	pubKey := byteToPublicKey(GetPubKeyAX(), GetPubKeyAY())

	return pubKey
}

func hexToPrivateKey(privKeyHex string) *ecdsa.PrivateKey {
	bytes, err := hex.DecodeString(privKeyHex)
	print(err)

	return byteToPrivateKey(bytes)
}

func byteToPrivateKey(privKey []byte) *ecdsa.PrivateKey {
	k := new(big.Int)
	k.SetBytes(privKey)

	priv := new(ecdsa.PrivateKey)
	curve := elliptic.P256()
	priv.PublicKey.Curve = curve
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())
	fmt.Printf("Calculated PubKey from PrivKey-->X: %d, Y: %d", priv.PublicKey.X, priv.PublicKey.Y)
	fmt.Printf("\n\n")

	return priv
}

func hexToPublicKey(pubKeyXpart string, pubKeyYpart string) *ecdsa.PublicKey {
	xBytes, _ := hex.DecodeString(pubKeyXpart)
	yBytes, _ := hex.DecodeString(pubKeyYpart)

	return byteToPublicKey(xBytes, yBytes)
}

func byteToPublicKey(pubKeyXpart []byte, pubKeyYpart []byte) *ecdsa.PublicKey {
	x := new(big.Int)
	x.SetBytes(pubKeyXpart)

	y := new(big.Int)
	y.SetBytes(pubKeyYpart)

	pub := new(ecdsa.PublicKey)
	pub.X = x
	pub.Y = y

	pub.Curve = elliptic.P256()

	return pub
}

type ecdsaSignature struct {
	R, S *big.Int
}

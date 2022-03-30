/* Copyright (C) beyond protocol inc. - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms and conditions defined in the file 'LICENSE.txt',
 * which is part of this source code package.
 */

package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/types"
)

func main() {
	if len(os.Args) != 2 {
		file := filepath.Base(os.Args[0])
		fmt.Println("Bech32 byte encoder")
		fmt.Println("Usage:")
		fmt.Println("\t" + file + " <bytes>")
		fmt.Println("\t\twhere <bytes> is a hex-string-encoded byte array")
		fmt.Println("Example:")
		fmt.Println("\t" + file + " 48656c6c6f20476f7068657221")
		fmt.Println("Remarks: resulting bech32 string will be prefixed with the following string: '" + types.Bech32PrefixAccAddr + "'.")
		return
	}
	argsWithoutProg := os.Args[1:]
	src := []byte(argsWithoutProg[0])
	accaddressbytes := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(accaddressbytes, src)
	if err != nil {
		log.Fatal(err)
	}
	accaddress := types.AccAddress(accaddressbytes)
	mykey := accaddress.String()
	fmt.Println("Bech32 string: " + mykey)
}

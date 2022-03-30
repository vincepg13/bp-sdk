/* Copyright (C) beyond protocol inc. - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms and conditions defined in the file 'LICENSE.txt',
 * which is part of this source code package.
 */

package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/types"
)

func main() {
	if len(os.Args) != 2 {
		file := filepath.Base(os.Args[0])
		fmt.Println("Bech32 string decoder")
		fmt.Println("Usage:")
		fmt.Println("\t" + file + " <bech32 encoded string>")
		fmt.Println("Example:")
		fmt.Println("\t" + file + " cosmosaccaddr1qgps2m09t2")
		fmt.Println("Remarks: it is assumed that the bech32 string uses the following prefix: '" + types.Bech32PrefixAccAddr + "'.")
		return
	}
	argsWithoutProg := os.Args[1:][0]
	//  src := []byte(argsWithoutProg[0])
	//  accaddressbytes := make([]byte, hex.DecodedLen(len(src)))
	//  _, err := hex.Decode(accaddressbytes, src)
	//  if err != nil {
	// 	 log.Fatal(err)
	//  }
	addr, _ := types.AccAddressFromBech32(argsWithoutProg)
	accountHex := "6163636F756E743A" // hex of the "account:" prefix used in the backend key-value store
	mykey := accountHex + hex.EncodeToString(addr)
	fmt.Println("Account bytes (hex): " + mykey)
}

/* Copyright (C) beyond protocol inc. - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms and conditions defined in the file 'LICENSE.txt',
 * which is part of this source code package.
 */

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type JsonBlock struct {
	Id     string `json:"id,omitempty"`
	Result Result `json:"result,omitempty"`
}

type Result struct {
	Block_meta BlockMeta `json:"block_meta,omitempty"`
	Block      Block     `json:"block,omitempty"`
}

type BlockMeta struct {
	Block_id BlockId `json:"block_id,omitempty"`
	Header   Header  `json:"header,omitempty"`
}

type BlockId struct {
	Hash string `json:"hash,omitempty"`
}

type Header struct {
	Chain_id         string    `json:"chain_id,omitempty"`
	Height           string    `json:"height,omitempty"`
	Time             time.Time `json:"time,omitempty"`
	Num_txs          string    `json:"num_txs,omitempty"`
	Total_txs        string    `json:"total_txs,omitempty"`
	App_hash         string    `json:"app_hash,omitempty"`
	Validators_hash  string    `json:"validators_hash,omitempty"`
	Last_commit_hash string    `json:"last_commit_hash,omitempty"`
	Consensus_hash   string    `json:"consensus_hash,omitempty"`
}

type Block struct {
	Data Data `json:"data,omitempty"`
}

type Data struct {
}

// our main function
func main() {
	router := mux.NewRouter()
	// API routes
	router.HandleFunc("/block/{id}", GetBlock).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// get a single block by ID
func GetBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id != "" {
		var netClient = &http.Client{
			Timeout: time.Second * 10,
		}

		response, err := netClient.Get("http://ec2-3-16-60-36.us-east-2.compute.amazonaws.com:26657/block?height=" + id)
		if err != nil {
			log.Fatal(err)
			json.NewEncoder(w).Encode(&Block{})
		} else {
			defer response.Body.Close()
			{
				blockData := &JsonBlock{}
				json.NewDecoder(response.Body).Decode(blockData)
				json.NewEncoder(w).Encode(blockData)
			}
		}
	}
}

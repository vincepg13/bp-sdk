# Introduction

The goal of beyond protocol is to develop a blockchain centric protocol to enable devices to communicate and exchange value in a trusted network through a universal protocol on top of TCP/IP.

The devices on the edge or the IoT devices must be able to:
- identify themselves
- communicate peer to peer in a secure way
- exchange value - incorporate a cryptocurrency wallet (if resources permit) 
- approve transactions through a consensus

Protocol features:
- identification
- credentials
- trust and authentication
- provisioning and configuration
- crypto / digital wallet capabilities
- consistent hashing
- dynamic ledger
  - TTL
  - Replication
  - Network latency
  - PUF randomness

# Install the SDK

On the host where you intend to run a Beyond protocol node, install *go* by following the [official docs](https://golang.org/doc/install). Remember to set your $GOPATH, $GOBIN, and $PATH environment variables.

Next, let's install the testnet's version of the Beyond protocol SDK. Here we'll use the develop branch, which contains the latest  release.

```
git clone https://github.com/vincepg13/bp-sdk
cd demo-sdk && git checkout develop
cd beyond
make get_tools && make get_vendor_deps && make build
```

That will install the beyondd and beyondcli binaries. Verify that everything is OK:

```
$ beyondcli version
$ beyondd version
```

The host is now ready to be configured as a Master node.

# Setting up a Master node

A Master node is running a full node, maintaining the entire blockchain.
First, initialize a full node by running the following command:

```
$ beyondd init --name <your_custom_node_name>
```

If the operation completed successfully, a node ID will be produced. Mark this information as you will need it later to complete the configuration.
This operation created a folder named '$HOME/.beyondd' with two subfolders: 'config' and 'data'. Inside the 'config' sub-folder there is a file named 'priv-validator.json' containing the validator private and public key. Also, there is the auto-generated genesis.json, containing validator information and initial account setup.

## Joining Beyond protocol testnet

If you want to join the master node to the Beyond protocol public testnet, replace the genesis.json file with the one from https://raw.githubusercontent.com/beyondprotocol/testnets/master/latest/genesis.json.
In order for your node to be able to find peers on the testnet, you should adjust the '$HOME/.beyondd/config/config.toml' file. Template for this configuration file can be found at https://raw.githubusercontent.com/beyondprotocol/testnets/master/latest/config.toml.
At the moment, the easiest way to find peers is to populate the 'persistent_peers' setting in the configuration file using the value provided in the template.

## Setting up your own testnet

If you are establishing a new network, you need to replace the genesis.json on each master node with your own template. First, you need to collect all the validator configurations (the 'priv-validator.json' files) from all Master nodes that you have initialized and insert their public keys into the genesis.json template file.
For example, if you have initialized four nodes, all of them being validators, the resulting section in the genesis.json file might look like the following snippet:

```
"validators": [
  {
    "pub_key": {
      "type": "bynd/PubKeyEd25519",
      "value": "ELtz5SZV6erWDQJdzfxzRDSSvYOfNVf3v3mq+P6sakU="
    },
    "power": "10",
    "name": ""
  },
  {
    "pub_key": {
      "type": "bynd/PubKeyEd25519",
      "value": "+DcGYN+mMIfCzT+wcwfLN6k4ToS3/6zZkr96uDJPGk4="
    },
    "power": "10",
    "name": ""
  },
  {
    "pub_key": {
      "type": "bynd/PubKeyEd25519",
      "value": "EYfOhuprfy6iuuTPQfzrj/O8bltrzDDPDW2brZ1u/hc="
    },
    "power": "10",
    "name": ""
  },
  {
    "pub_key": {
      "type": "bynd/PubKeyEd25519",
      "value": "gLctQ102Y6Hfo+WtIkhDogHjaTjebDfvvJ4AFDJ0LcM="
    },
    "power": "10",
    "name": ""
  }
]
```

Once you have the genesis.json file ready on your development machine, it is time to upload it to all validator nodes.

## Configuring Master node

Whether you are joining an existing testnet or building up a new one, additional configuration needs to be taken care of on each new Master node. You can specify the node configuration in the .beyondd/config/config.toml file that has been automatically generated when you ran the 'beyondd init' command.
The following settings need to be adjusted:
- Specify a unique validator node name in the 'moniker' property.
- Each node needs to be able to find other peers on the network. Specify a comma-separated list of peer validator nodes in the 'persistent_peers' property. Each node is identified by its ID, returned by the 'beyondd init' command, and its IP address, e.g. 846e22c4641d07258daa99c7fb455d28588fc29b@xxx.xxx.xxx.xxx:26656
- Set the value of the 'addr_book_strict' property to false if the peer validator nodes are running in private network using private IP addresses.

# Starting a Master node

Once the configuration is complete, you can start the full node with the following command:

```
$ beyondd start
```

After at least three validators are running you should see blocks being periodically created and committed.

# Setting up accounts

Accounts with their private keys are stored in "$HOME/.beyondcli". Private keys are stored armored (encrypted using a user-provided password) in database and are unarmored by beyondcli each time we sign a transaction.

Besides generic accounts we support mobility accounts used in scenarios related to wireless charging of vehicles. We use special DeepCover Secure Authenticator chip, which stores encrypted keys in EEPROM. Once the key is generated and stored, the EEPROM address is made write protected. When using this type of account, the transaction payload is signed by the chip. Transaction is verified on MasterNode with ECDSA algorithm using correct parameters such as Curve equation, X and Y part of public key, message digest and R,S parts of the signature.

The following command creates a new account on the client node:

```
$ beyondcli keys add fastcar
Enter a passphrase for your key:
Repeat the passphrase:
NAME:   TYPE:   ADDRESS:                                            PUBKEY:
fastcar local   byndaddr19vvpsyqrfx8mha8kh229d3g23k4f6mcyfffr5e     byndpub1addwnpepqtduxppjugzdgwslw9en929pc8eqxz0yffs4qlqawc257f0gs0y3uxmvn5p
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.
 
choice captain minor case grunt fragile blanket creek crane act maid sorry kiss brass glance rely silly lesson lab picnic close goddess sick lawsuit
```

## Listing accounts

Use the following command of beyondcli to list existing accounts, registered locally:

```
$ beyondcli keys list
NAME:   TYPE:   ADDRESS:                                          PUBKEY:
car     local   byndaddr1c4r7cgdyzpkdnnet8h8lgwf5n7zq3ctddggfx3   byndpub1addwnpepqfqjy2rh27cmvxx0m0h84j4z6kzxvr6m4ad40hrltqu6wl8cp4frxz4mjhy
fastcar local   byndaddr19vvpsyqrfx8mha8kh229d3g23k4f6mcyfffr5e   byndpub1addwnpepqtduxppjugzdgwslw9en929pc8eqxz0yffs4qlqawc257f0gs0y3uxmvn5p
station local   byndaddr1mv2e82j0sxrg5c7dr7t5nc8cc6zvhgyqs9phrd   byndpub1addwnpepqdw295n2flsgzeldl4hggpkeunlpaqhwv03mwpxz62uelhc0k6ny7au7g56
tesla   local   byndaddr12jltzqm663qt0rz5y354d06fgtv5mtndsq3t23   byndpub1addwnpepqwze7u7ef0jfs3ehxfl6ghs6tzvyqlsez8ahg2dxkalquxn7ha345euhnsq
```

By default the configuration is taken from "$HOME/.beyondcli". We can override this by specifying --home parameter.

## Query account balance

Here we query a charging station account, which shows among other fields:
- balance of coins and their denomination,
- sequence number to prevent replay attacks and
- current price in kWh at which it sells energy. 

```
$ beyondcli account byndaddr1mv2e82j0sxrg5c7dr7t5nc8cc6zvhgyqs9phrd --node=beyond.link:26657
```

Every client communicates with a designated Master node that we choose using the --node parameter. The parameter specifies the RCP endpoint of the Master node, which is by default exposed on port 26657.

If the communication succeeds, JSON output of the above command will be:

```
{
  "type": "beyond/Account",
  "value": {
    "BaseAccount": {
      "address": "byndaddr1mv2e82j0sxrg5c7dr7t5nc8cc6zvhgyqs9phrd",
      "coins": [
        {
          "denom": "byndcoin",
          "amount": "500000"
        }
      ],
      "public_key": {
        "type": "tendermint/PubKeySecp256k1",
        "value": "byndpub1addwnpepqdw295n2flsgzeldl4hggpkeunlpaqhwv03mwpxz62uelhc0k6ny7au7g56"
      },
      "account_number": "1",
      "sequence": "100"
    },
    "name": "station",
    "macAddress": "00-05-9A-3C-7A-00",
    "price": "10"
  }
}
```

## InitOrder command

InitOrder initializes new order and deterministically increments orderNumber for buyer of energy. OrderNumber is kept in application state on Master nodes for each beyond account.

```
beyondcli initOrder --from=car --amount=2 --to=byndaddr1j6j0f6kvs92zazh3tmjvvu8ga8x4h9hf5jxum9 --sequence=0 --chain-id=beyond-chain --node=beyond.link:26657
Password to sign with 'car':
Committed at block 43 (tx hash: 378EA8577E67DF99AB4E159A35E9E439CEBC9023)
```

## FinalizeOrder command

FinalizeOrder performs actual transfer of coins which is calculated by provided charge in kWh multiplied by current price of the seller specified in "--to" address. FinalizeOrder checks correct InitOrder and links it via OrderNumber. Transaction details (inputs, outputs) and orderNumber are appended to transaction tags. Those can be queried later for analytics.

```
beyondcli finalizeOrder --from car --charge=2 --to=byndaddr1svf6jrfvfed33avtza9y9h8ckmcylpqsseeysn --sequence=0 --chain-id=beyond-chain --node=beyond.link:26657
Password to sign with 'station':
Committed at block 83 (tx hash: 5D2219CE78657A2A6D462B3DCA79E47601FFFE80)
```

# Querying the blockchain

Blockchain data can be queried via the light client (beyondcli) using its CLI interface or the REST API.
For instance, to query a committed transaction using an attached tag, the following CLI command could be used:

```
$ beyondcli tendermint txs --tag orderNumber=15453 --node=beyond.link:26657
```

The same query could be made using the exposed REST API. For this, the client should first be initialized with a REST server running. Here, the ID of the blockchain must also be supplied.

```
$ beyondcli rest-server --node=beyond.link:26657 --laddr "tcp://localhost:26650 --chain-id beyond-chain"
```

Now we can issue a request against the endpoint.

```
$ curl "http://localhost:26650/txs?tag=orderNumber=15453"
```

If we wanted to list all transactions committed in a block at a given height, this is the command we would issue:

```
$ curl "http://localhost:26650/txs?tag=tx.height=175572"
```

By default, the light client verifies the query results obtained from the blockchain (using the validator public keys, latest block hash and the Merkle proofs provided together with the query results). Also, the REST server establishes a secure SSL channel with the connecting client using a self-signed certificate. 
Query response verification can be disabled when starting a REST server:

```
$ beyondcli rest-server --node=beyond.link:26657 --laddr "tcp://localhost:26650 --trust-node"
```

The SSL/TLS layer can also be skipped during initialization of the REST server:

```
$ beyondcli rest-server --node=beyond.link:26657 --laddr "tcp://localhost:26650 --chain-id beyond-chain --insecure"
```






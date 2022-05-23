# btcio
Tool to perform Bitcoin transactions

## disclaimer
> The use of this tool does not guarantee security or usability for any
> particular purpose. Please review the code and use at your own risk.

> Don't trust, verify

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system.

```bash
go install github.com/kubetrail/btcio@latest
```
Add autocompletion for `bash` to your `.bashrc`
```bash
source <(btcio completion bash)
```

## runtime prerequisite
`bitcoind`, needs to be serving RPC that this tool connects to. Furthermore,
make sure that _all_ transactions are indexed using config option
`txindex=1` shown below.

Download [Bitcoin core](https://bitcoin.org/en/download)

Create a configuration file:
```bash
cat ~/.bitcoin/bitcoin.conf
```
```text
rpcuser=yourrpcuser
rpcpassword=yourrpcpass
server=1
testnet=0
txindex=1
```

Set `testnet=1`, when running on testnet.

`rpcuser` and `rpcpassword` will be deprecated in the future, however, they work for now.
`regtest` implies regression test setup, which we switch off implying connecting
to the `mainnet`

`txindex` indexes all transactions, not just those for the associated wallet.

Run daemon:
```bash
./bitcoind
```
```text
2022-05-20T00:29:26Z Bitcoin Core version v23.0.0 (release build)
2022-05-20T00:29:26Z Assuming ancestors of block 000000000000000000052d314a259755ca65944e68df6b12a067ea8f1f5a7091 have valid signatures.
2022-05-20T00:29:26Z Setting nMinimumChainWork=00000000000000000000000000000000000000002927cdceccbd5209e81e80db
2022-05-20T00:29:26Z Using the 'sse4(1way),sse41(4way),avx2(8way)' SHA256 implementation
2022-05-20T00:29:26Z Using RdSeed as additional entropy source
2022-05-20T00:29:26Z Using RdRand as an additional entropy source
2022-05-20T00:29:26Z Default data directory /home/sdeoras/.bitcoin
2022-05-20T00:29:26Z Using data directory /home/sdeoras/.bitcoin
2022-05-20T00:29:26Z Config file: /home/sdeoras/.bitcoin/bitcoin.conf
2022-05-20T00:29:26Z Config file arg: regtest="0"
2022-05-20T00:29:26Z Config file arg: rpcpassword=****
2022-05-20T00:29:26Z Config file arg: rpcuser=****
2022-05-20T00:29:26Z Config file arg: server="1"
2022-05-20T00:29:26Z Config file arg: txindex="1"
2022-05-20T00:29:26Z Using at most 125 automatic connections (1024 file descriptors available)
2022-05-20T00:29:26Z Using 16 MiB out of 32/2 requested for signature cache, able to store 524288 elements
2022-05-20T00:29:26Z Using 16 MiB out of 32/2 requested for script execution cache, able to store 524288 elements
2022-05-20T00:29:26Z Script verification uses 11 additional threads
2022-05-20T00:29:26Z scheduler thread start
2022-05-20T00:29:26Z HTTP: creating work queue of depth 16
2022-05-20T00:29:26Z Config options rpcuser and rpcpassword will soon be deprecated. Locally-run instances may remove rpcuser to use cookie-based auth, or may be replaced with rpcauth. Please see share/rpcauth for rpcauth auth generation.
2022-05-20T00:29:26Z HTTP: starting 4 worker threads
2022-05-20T00:29:26Z Using wallet directory /home/sdeoras/.bitcoin/wallets
2022-05-20T00:29:26Z init message: Verifying wallet(s)…
2022-05-20T00:29:26Z Using SQLite Version 3.32.1
2022-05-20T00:29:26Z Using wallet /home/sdeoras/.bitcoin/wallets
2022-05-20T00:29:26Z Using SQLite Version 3.32.1
2022-05-20T00:29:26Z Using wallet /home/sdeoras/.bitcoin/wallets
2022-05-20T00:29:26Z Using /16 prefix for IP bucketing
2022-05-20T00:29:26Z init message: Loading P2P addresses…
2022-05-20T00:29:26Z Loaded 22930 addresses from peers.dat  68ms
2022-05-20T00:29:26Z init message: Loading banlist…
2022-05-20T00:29:26Z SetNetworkActive: true
... more such lines
2022-05-20T00:28:38Z UpdateTip: new best=00000000000000000018341c712847a87f6f498a3f6ff7148e3759758c082556 height=488965 version=0x20000000 log2_work=87.247470 tx=260369222 date='2017-10-09T00:25:59Z' progress=0.354673 cache=566.7MiB(4264665txo)
2022-05-20T00:28:38Z UpdateTip: new best=00000000000000000089ddba3d14c6a099f5e2f3be8e53a6cfeba96aefe11c79 height=488966 version=0x20000000 log2_work=87.247508 tx=260369425 date='2017-10-09T00:26:49Z' progress=0.354673 cache=566.5MiB(4263566txo)
2022-05-20T00:28:38Z UpdateTip: new best=000000000000000000c17af840a28377518361f842aec03cf963a129e381c96c height=488967 version=0x20000000 log2_work=87.247546 tx=260371459 date='2017-10-09T00:38:59Z' progress=0.354676 cache=566.8MiB(4265775txo)
2022-05-20T00:28:38Z UpdateTip: new best=00000000000000000046bde1140c5e4bbd2f66b5b3744e68a5472285f2825c31 height=488968 version=0x20000000 log2_work=87.247583 tx=260371460 date='2017-10-09T00:39:07Z' progress=0.354676 cache=566.8MiB(4265776txo)
2022-05-20T00:28:38Z UpdateTip: new best=000000000000000000206ced5efeee63f05eb55a5c3429caa76139b668bb4a7a height=488969 version=0x20000000 log2_work=87.247621 tx=260373156 date='2017-10-09T00:50:13Z' progress=0.354678 cache=567.0MiB(4267286txo)
2022-05-20T00:28:38Z UpdateTip: new best=000000000000000000b8b20882f2eb4a7bae281fbab4a535b0130325e8969363 height=488970 version=0x20000000 log2_work=87.247659 tx=260375346 date='2017-10-09T01:36:30Z' progress=0.354681 cache=567.3MiB(4269581txo)
```

It can take several hours to a few days for `bitcond` to sync up.

### make rpc calls via curl
Get hash for the first block!
```bash
curl --user yourrpcuser:yourrpcpass --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getblockhash", "params": [0]}' -H 'content-type: text/plain;' http://127.0.0.1:8332/
```
```json
{"result":"000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f","error":null,"id":"curltest"}
```

Using this hash, get the block details:
```bash
curl --user yourrpcuser:yourrpcpass \
  --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getblock", "params": ["000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"]}' -H 'content-type: text/plain;' \
  http://127.0.0.1:8332/ \
  | jq '.'
```
```json
{
  "result": {
    "hash": "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
    "confirmations": 399646,
    "height": 0,
    "version": 1,
    "versionHex": "00000001",
    "merkleroot": "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
    "time": 1231006505,
    "mediantime": 1231006505,
    "nonce": 2083236893,
    "bits": "1d00ffff",
    "difficulty": 1,
    "chainwork": "0000000000000000000000000000000000000000000000000000000100010001",
    "nTx": 1,
    "nextblockhash": "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048",
    "strippedsize": 285,
    "size": 285,
    "weight": 1140,
    "tx": [
      "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"
    ]
  },
  "error": null,
  "id": "curltest"
}
```

It has just one transaction, which can be further inspected:
```bash
curl \
  --user yourrpcuser:yourrpcpass \
  --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getrawtransaction", "params": ["4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"]}' -H 'content-type: text/plain;' \
  http://127.0.0.1:8332/ \
  | jq '.'
```
```json
{
  "result": null,
  "error": {
    "code": -5,
    "message": "The genesis block coinbase is not considered an ordinary transaction and cannot be retrieved"
  },
  "id": "curltest"
}
```

As you see, that transaction was the genesis transaction and has some restrictions.
Let's query next block:
```bash
curl \
  --user yourrpcuser:yourrpcpass \
  --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getblockhash", "params": [1]}' -H 'content-type: text/plain;' \
  http://127.0.0.1:8332/
```
```json
{"result":"00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048","error":null,"id":"curltest"}
```

```bash
curl --user yourrpcuser:yourrpcpass \
  --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getblock", "params": ["00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048"]}' -H 'content-type: text/plain;' \
  http://127.0.0.1:8332/ \
  | jq '.'
```
```json
{
  "result": {
    "hash": "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048",
    "confirmations": 511344,
    "height": 1,
    "version": 1,
    "versionHex": "00000001",
    "merkleroot": "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
    "time": 1231469665,
    "mediantime": 1231469665,
    "nonce": 2573394689,
    "bits": "1d00ffff",
    "difficulty": 1,
    "chainwork": "0000000000000000000000000000000000000000000000000000000200020002",
    "nTx": 1,
    "previousblockhash": "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
    "nextblockhash": "000000006a625f06636b8bb6ac7b960a8d03705d1ace08b1a19da3fdcc99ddbd",
    "strippedsize": 215,
    "size": 215,
    "weight": 860,
    "tx": [
      "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098"
    ]
  },
  "error": null,
  "id": "curltest"
}
```

Investigate the transaction
```bash
curl \
  --user yourrpcuser:yourrpcpass \
  --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getrawtransaction", "params": ["0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098"]}' -H 'content-type: text/plain;' \
  http://127.0.0.1:8332/ \
  | jq '.'
```
```json
{
  "result": "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704ffff001d0104ffffffff0100f2052a0100000043410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac00000000",
  "error": null,
  "id": "curltest"
}
```

Decode the transaction:
```bash
curl \
  --user yourrpcuser:yourrpcpass \
  --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "decoderawtransaction", "params": ["01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704ffff001d0104ffffffff0100f2052a0100000043410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac00000000"]}' -H 'content-type: text/plain;' \
  http://127.0.0.1:8332/ \
  | jq '.'
```
```json
{
  "result": {
    "txid": "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
    "hash": "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
    "version": 1,
    "size": 134,
    "vsize": 134,
    "weight": 536,
    "locktime": 0,
    "vin": [
      {
        "coinbase": "04ffff001d0104",
        "sequence": 4294967295
      }
    ],
    "vout": [
      {
        "value": 50.00000000,
        "n": 0,
        "scriptPubKey": {
          "asm": "0496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858ee OP_CHECKSIG",
          "desc": "pk(0496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858ee)#qnv32gt7",
          "hex": "410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac",
          "type": "pubkey"
        }
      }
    ]
  },
  "error": null,
  "id": "curltest"
}
```

## query transactions
Assuming steps above are working, we can now perform bitcoin transactions. First, it is required
to obtain UTXO for the sender, which is done by inspecting the last transaction

You can inspect the transaction [here](https://www.blockchain.com/btc/tx/0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098)

Examples below are shown for the BTC testnet, so make sure `bitcoind` is running as a testnet node.
```bash
btcio --rpc-endpoint=127.0.0.1:18332 \
  utxo 72a49ff05829f6c31a089a9c7413498cb18190ffc839208e67a27cc15933a298 \
  | jq '.'
```
```json
{
  "hex": "",
  "txid": "72a49ff05829f6c31a089a9c7413498cb18190ffc839208e67a27cc15933a298",
  "hash": "72a49ff05829f6c31a089a9c7413498cb18190ffc839208e67a27cc15933a298",
  "size": 109,
  "vsize": 109,
  "weight": 436,
  "version": 1,
  "locktime": 0,
  "vin": [
    {
      "coinbase": "04e10e4a4d0169062f503253482f",
      "sequence": 4294967295
    }
  ],
  "vout": [
    {
      "value": 50,
      "n": 0,
      "scriptPubKey": {
        "asm": "02a741071164b40b01c4ad28913c4aa2a1015cc5b064f0c802272552f17ae08750 OP_CHECKSIG",
        "hex": "2102a741071164b40b01c4ad28913c4aa2a1015cc5b064f0c802272552f17ae08750ac",
        "type": "pubkey"
      }
    }
  ]
}
```

Based on this output the script can be obtained `2102a741071164b40b01c4ad28913c4aa2a1015cc5b064f0c802272552f17ae08750ac`
and we also get to know the unspent amount in `value`

## send transaction
In order to send transaction you will need to know following:
* previous transaction hash
* script hex string
* sender's private key in WIF format
* receiver's address in compressed format
* amount to send in sats
```bash
btcio send --amount=5000
```

The output is the transaction ID that can be used to verify it

## references
* [Bitcoin API reference](https://developer.bitcoin.org/reference/rpc/index.html)
* https://stackoverflow.com/questions/20481478/how-to-check-bitcoin-address-balance-from-my-application
* https://bitcoin.stackexchange.com/questions/59179/determining-if-input-is-a-block-id-transaction-id-or-address
* https://www.reddit.com/r/crypto/comments/guctw4/finding_sha256_partial_collisions_via_the_bitcoin/
* https://bitcoin.stackexchange.com/questions/1781/nonce-size-will-it-always-be-big-enough
* https://levelup.gitconnected.com/bitcoin-proof-of-work-the-only-article-you-will-ever-have-to-read-4a1fcd76a294
* https://medium.com/swlh/create-raw-bitcoin-transaction-and-sign-it-with-golang-96b5e10c30aa
* https://en.bitcoin.it/wiki/Script
* https://bitcoin.stackexchange.com/questions/96865/why-does-vout-sometimes-not-have-address
* 

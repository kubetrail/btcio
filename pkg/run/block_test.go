package run

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/kubetrail/btcio/pkg/flags"
)

func TestBlock(t *testing.T) {
	connCfg := &rpcclient.ConnConfig{
		Host:         flags.DefaultRPCEndpoint,
		User:         flags.DefaultRPCUser,
		Pass:         flags.DefaultRPCPassword,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
		Certificates: nil,
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to get a new rpc client: %w", err))
	}
	defer client.Shutdown()

	blockHash, err := client.GetBlockHash(0)
	if err != nil {
		t.Fatal(err)
	}

	block, err := client.GetBlock(blockHash)
	if err != nil {
		t.Fatal(err)
	}

	outBlock, err := Parse(block)
	if err != nil {
		t.Fatal(err)
	}

	outBlock.Header.Hash = hex.EncodeToString((*blockHash)[:])
	count := 0
	h := *blockHash
	for i := len(h) - 1; i >= 0; i-- {
		if h[i] != 0 {
			break
		}
		count++
	}

	outBlock.Header.NumContiguousZerosInHash = count

	jb, err := json.MarshalIndent(outBlock, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(jb))
}

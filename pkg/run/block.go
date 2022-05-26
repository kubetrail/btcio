package run

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/kubetrail/bip32/pkg/keys"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/kubetrail/btcio/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Block defines a custom view of actual block in the blockchain
type Block struct {
	Header       *BlockHeader `json:"header,omitempty" yaml:"header,omitempty"`
	Transactions []string     `json:"transactions,omitempty" yaml:"transactions,omitempty"`
}

// BlockHeader defines information about a block
type BlockHeader struct {
	// Version of the block.  This is not the same as the protocol version.
	Version int32 `json:"version,omitempty" yaml:"version,omitempty"`

	// Hash of the block
	Hash string `json:"hash,omitempty" yaml:"hash,omitempty"`

	// NumContiguousZerosInHash is the number of trailing zeros in the binary version of the hash
	NumContiguousZerosInHash int `json:"numContiguousZerosInHash,omitempty" yaml:"numContiguousZerosInHash,omitempty"`

	// PrevBlock is the hash of the previous block header in the blockchain
	PrevBlock string `json:"prevBlock,omitempty" yaml:"prevBlock,omitempty"`

	// Merkle tree reference to hash of all transactions for the block.
	MerkleRoot string `json:"merkleRoot,omitempty" yaml:"merkleRoot,omitempty"`

	// Time the block was created.  This is, unfortunately, encoded as a
	// uint32 on the wire and therefore is limited to 2106.
	Timestamp time.Time `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`

	// Difficulty target for the block.
	Bits uint32 `json:"bits,omitempty" yaml:"bits,omitempty"`

	// Nonce used to generate the block.
	Nonce uint32 `json:"nonce,omitempty" yaml:"nonce,omitempty"`
}

func reverseString(input string) string {
	b := []byte(input)
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-1-i] = b[len(b)-1-i], b[i]
	}

	return string(b)
}

func Parse(block *wire.MsgBlock, chainHash *chainhash.Hash) (*Block, error) {
	b := &Block{
		Header: &BlockHeader{
			Version:    block.Header.Version,
			PrevBlock:  reverseString(block.Header.PrevBlock.String()),
			MerkleRoot: reverseString(block.Header.MerkleRoot.String()),
			Timestamp:  block.Header.Timestamp,
			Bits:       block.Header.Bits,
			Nonce:      block.Header.Nonce,
		},
		Transactions: make([]string, len(block.Transactions)),
	}

	for i := range block.Transactions {
		b.Transactions[i] = block.Transactions[i].TxHash().String()
	}

	b.Header.Hash = (*chainHash).String()

	count := 0

	bb := &bytes.Buffer{}
	bw := bufio.NewWriter(bb)

	for _, hashByte := range *chainHash {
		if _, err := fmt.Fprintf(bw, "%08b", hashByte); err != nil {
			return nil, fmt.Errorf("failed to write to output: %w", err)
		}
	}

	if err := bw.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush output: %w", err)
	}

	hashBytesAsString := bb.Bytes()
	for i := len(hashBytesAsString) - 1; i >= 0; i-- {
		if hashBytesAsString[i] != '0' {
			break
		}
		count++
	}

	b.Header.NumContiguousZerosInHash = count

	return b, nil
}

func GetBlock(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.BlockHash, cmd.Flag(flags.BlockHash))
	blockHash := viper.GetString(flags.BlockHash)

	prompt, err := prompts.Status()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	if len(blockHash) == 0 {
		if len(args) == 0 {
			if prompt {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter block hash or number: "); err != nil {
					return fmt.Errorf("failed to write to output: %w", err)
				}
			}
			blockHash, err = keys.Read(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read pub addr from input: %w", err)
			}
		} else {
			blockHash = args[0]
		}
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         persistentFlags.RPCEndpoint,
		User:         persistentFlags.RPCUser,
		Pass:         persistentFlags.RPCPassword,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
		Certificates: nil,
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return fmt.Errorf("failed to get a new rpc client: %w", err)
	}
	defer client.Shutdown()

	var chainHash *chainhash.Hash
	if len(blockHash) < 10 {
		x, err := strconv.ParseInt(blockHash, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse block number as decimal integer: %w", err)
		}

		chainHash, err = client.GetBlockHash(x)
		if err != nil {
			return fmt.Errorf("failed to get block hash: %w", err)
		}

		blockHash = chainHash.String()
	} else {
		chainHash, err = chainhash.NewHashFromStr(blockHash)
		if err != nil {
			return fmt.Errorf("failed to get new chainhash from block hash string: %w", err)
		}
	}

	block, err := client.GetBlock(chainHash)
	if err != nil {
		return fmt.Errorf("failed to get block for given hash: %w", err)
	}

	outBlock, err := Parse(block, chainHash)
	if err != nil {
		return fmt.Errorf("failed to parse block: %w", err)
	}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatYaml:
		b, err := yaml.Marshal(outBlock)
		if err != nil {
			return fmt.Errorf("failed to serialize output to yaml: %w", err)
		}

		if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
			return fmt.Errorf("failed to write yaml to output: %w", err)
		}
	case flags.OutputFormatJson:
		b, err := json.Marshal(outBlock)
		if err != nil {
			return fmt.Errorf("failed to serialize output to json: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(b)); err != nil {
			return fmt.Errorf("failed to write json to output: %w", err)
		}
	case flags.OutputFormatNative:
		b, err := json.MarshalIndent(outBlock, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to serialize output to json: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(b)); err != nil {
			return fmt.Errorf("failed to write json to output: %w", err)
		}
	default:
		return fmt.Errorf("invalid output format")
	}

	return nil
}

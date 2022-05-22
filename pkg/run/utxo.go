package run

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/kubetrail/btcio/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Utxo(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.TransactionHash, cmd.Flag(flags.TransactionHash))
	txHash := viper.GetString(flags.TransactionHash)

	if len(txHash) == 0 && len(args) > 0 {
		txHash = args[0]
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

	tx, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return fmt.Errorf("failed to generate new chainhash: %w", err)
	}

	txRaw, err := client.GetRawTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	bb := &bytes.Buffer{}
	bw := bufio.NewWriter(bb)
	if err := txRaw.MsgTx().Serialize(bw); err != nil {
		return fmt.Errorf("failed to serialize tx msg: %w", err)
	}

	if err := bw.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	txResults, err := client.DecodeRawTransaction(bb.Bytes())
	if err != nil {
		return fmt.Errorf("failed to decode raw transaction: %w", err)
	}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatNative, flags.OutputFormatYaml:
		b, err := yaml.Marshal(txResults)
		if err != nil {
			return fmt.Errorf("failed to serialize output to yaml: %w", err)
		}

		if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
			return fmt.Errorf("failed to write yaml to output: %w", err)
		}
	case flags.OutputFormatJson:
		b, err := json.Marshal(txResults)
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

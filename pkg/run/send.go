package run

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/kubetrail/btcio/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Send(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.Network, cmd.Flag(flags.Network))
	_ = viper.BindPFlag(flags.TransactionHash, cmd.Flag(flags.TransactionHash))
	_ = viper.BindPFlag(flags.PkScript, cmd.Flag(flags.PkScript))
	_ = viper.BindPFlag(flags.Key, cmd.Flag(flags.Key))
	_ = viper.BindPFlag(flags.Addr, cmd.Flag(flags.Addr))
	_ = viper.BindPFlag(flags.Amount, cmd.Flag(flags.Amount))

	network := viper.GetString(flags.Network)
	txHash := viper.GetString(flags.TransactionHash)
	pkScript := viper.GetString(flags.PkScript)
	privKey := viper.GetString(flags.Key)
	addr := viper.GetString(flags.Addr)
	amount := viper.GetInt64(flags.Amount)

	var params chaincfg.Params
	switch strings.ToLower(network) {
	case flags.DefaultNetworkMainnet:
		params = chaincfg.MainNetParams
	case flags.DefaultNetworkTestnet:
		params = chaincfg.TestNet3Params
	default:
		return fmt.Errorf("invalid network, only mainnet or testnet are allowed")
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

	newTx, err := CreateTx(
		&txConfig{
			txId:     txHash,
			pkScript: pkScript,
			privKey:  privKey,
			addr:     addr,
			amount:   amount,
			params:   &params,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create a new transaction: %w", err)
	}

	chainHash, err := client.SendRawTransaction(newTx, false)
	if err != nil {
		return fmt.Errorf("failed to send raw transaction: %w", err)
	}

	type output struct {
		TxHash string `json:"txHash,omitempty" yaml:"txHash,omitempty"`
	}

	out := &output{TxHash: chainHash.String()}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatNative:
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), out.TxHash); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	case flags.OutputFormatYaml:
		b, err := yaml.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to serialize output to yaml: %w", err)
		}

		if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
			return fmt.Errorf("failed to write yaml to output: %w", err)
		}
	case flags.OutputFormatJson:
		b, err := json.Marshal(out)
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

type txConfig struct {
	txId     string
	pkScript string
	privKey  string
	addr     string
	amount   int64
	params   *chaincfg.Params
}

func CreateTx(config *txConfig) (*wire.MsgTx, error) {
	// 1 or unit-amount in Bitcoin is equal to 1 satoshi and 1 Bitcoin = 100000000 satoshi

	// extracting destination address as []byte from function argument (destination string)
	destinationAddr, err := btcutil.DecodeAddress(config.addr, config.params)
	if err != nil {
		return nil, fmt.Errorf("failed to decode addr: %w", err)
	}

	destinationAddrByte, err := txscript.PayToAddrScript(destinationAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination addr byte: %w", err)
	}

	// creating a new bitcoin transaction, different sections of the tx, including
	// input list (contain UTXOs) and output-list (contain destination address and usually our address)
	// in next steps, sections will be field and pass to sign
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	utxoHash, err := chainhash.NewHashFromStr(config.txId)
	if err != nil {
		return nil, fmt.Errorf("faied to generate new chainhash: %w", err)
	}

	// the second argument is vout or Tx-index, which is the index
	// of spending UTXO in the transaction that Txid referred to
	// in this case is 0, but can vary different numbers
	outPoint := wire.NewOutPoint(utxoHash, 0)

	// making the input, and adding it to transaction
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	// adding the destination address and the amount to
	// the transaction as output
	redeemTxOut := wire.NewTxOut(config.amount, destinationAddrByte)
	redeemTx.AddTxOut(redeemTxOut)

	// now sign the transaction
	if err := SignTx(config.privKey, config.pkScript, redeemTx); err != nil {
		return nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	return redeemTx, nil
}

func SignTx(privKey string, pkScript string, redeemTx *wire.MsgTx) error {

	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return err
	}

	sourcePKScript, err := hex.DecodeString(pkScript)
	if err != nil {
		return nil
	}

	// since there is only one input in our transaction
	// we use 0 as second argument, if the transaction
	// has more args, should pass related index
	signature, err := txscript.SignatureScript(redeemTx, 0, sourcePKScript, txscript.SigHashAll, wif.PrivKey, true)
	if err != nil {
		return nil
	}

	// since there is only one input, and want to add
	// signature to it use 0 as index
	redeemTx.TxIn[0].SignatureScript = signature

	return nil
}

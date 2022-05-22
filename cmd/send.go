/*
Copyright © 2022 kubetrail.io authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/kubetrail/btcio/pkg/flags"
	"github.com/kubetrail/btcio/pkg/run"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send Bitcoins",
	Long: `Sending BTC transaction requires knowledge of
previous transaction id for the sender address and
the associated public key script. Assuming that the
previous transaction id is known, the public key script
can be obtained by parsing the transaction details.
`,
	RunE: run.Send,
}

func init() {
	rootCmd.AddCommand(sendCmd)
	f := sendCmd.Flags()

	f.String(flags.Network, flags.DefaultNetworkMainnet, "BTC network: mainnet or testnet")
	f.String(flags.TransactionHash, "", "Previous transaction hash")
	f.String(flags.PkScript, "", "PK Script")
	f.String(flags.Key, "", "Private WIF key")
	f.String(flags.Addr, "", "Receiver addr")
	f.Int64(flags.Amount, 0, "Amount to send (Sats)")
}

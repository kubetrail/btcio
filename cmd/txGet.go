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

// utxoCmd represents the txGet command
var utxoCmd = &cobra.Command{
	Use:   "utxo",
	Short: "Get transaction UTXO, script etc.",
	Long:  ``,
	RunE:  run.Utxo,
	Args:  cobra.MaximumNArgs(1),
}

func init() {
	rootCmd.AddCommand(utxoCmd)
	f := utxoCmd.Flags()

	f.String(flags.TransactionHash, "", "Transaction hash")
}

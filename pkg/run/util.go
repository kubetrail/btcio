package run

import (
	"github.com/kubetrail/btcio/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type persistentFlagValues struct {
	OutputFormat string `json:"outputFormat,omitempty"`
	RPCEndpoint  string `json:"rpcEndPoint,omitempty"`
	RPCUser      string `json:"rpcUser,omitempty"`
	RPCPassword  string `json:"rpcPassword,omitempty"`
}

func getPersistentFlags(cmd *cobra.Command) persistentFlagValues {
	rootCmd := cmd.Root().PersistentFlags()

	_ = viper.BindPFlag(flags.OutputFormat, rootCmd.Lookup(flags.OutputFormat))
	_ = viper.BindPFlag(flags.RPCEndpoint, rootCmd.Lookup(flags.RPCEndpoint))
	_ = viper.BindPFlag(flags.RPCUser, rootCmd.Lookup(flags.RPCUser))
	_ = viper.BindPFlag(flags.RPCPassword, rootCmd.Lookup(flags.RPCPassword))

	_ = viper.BindEnv(flags.RPCEndpoint, flags.DefaultRPCEndpointEnvVarKey)

	outputFormat := viper.GetString(flags.OutputFormat)
	rpcEndpoint := viper.GetString(flags.RPCEndpoint)
	rpcUser := viper.GetString(flags.RPCUser)
	rpcPassword := viper.GetString(flags.RPCPassword)

	return persistentFlagValues{
		OutputFormat: outputFormat,
		RPCEndpoint:  rpcEndpoint,
		RPCUser:      rpcUser,
		RPCPassword:  rpcPassword,
	}
}

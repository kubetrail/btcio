package flags

const (
	RPCEndpoint         = "rpc-endpoint"
	RPCUser             = "rpc-user"
	RPCPassword         = "rpc-password"
	OutputFormat        = "output-format"
	Network             = "network"
	TransactionHash     = "tx-hash"
	BlockHash           = "block-hash"
	PkScriptNumber      = "pk-script-numer"
	PkScriptNumberShort = "n"
	PkScript            = "pk-script"
	Amount              = "amount"
	Key                 = "key"
	Addr                = "addr"
)

const (
	DefaultRPCEndpoint = "127.0.0.1:8332"
	DefaultRPCUser     = "yourrpcuser"
	DefaultRPCPassword = "yourrpcpass"

	DefaultRPCEndpointEnvVarKey = "RPC_ENDPOINT"

	DefaultNetworkMainnet = "mainnet"
	DefaultNetworkTestnet = "testnet"

	OutputFormatNative = "native"
	OutputFormatJson   = "json"
	OutputFormatYaml   = "yaml"
)

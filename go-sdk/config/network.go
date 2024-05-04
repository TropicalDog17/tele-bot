package config

import (
	"github.com/InjectiveLabs/sdk-go/client/common"
)

type Config struct {
	NetworkType          string
	ChainId              string
	ExchangeGrpcEndpoint string
	ExplorerGrpcEndpoint string
	TmEndpoint           string
	ChainGrpcEndpoint    string
	LcdEndpoint          string
	Fee_denom            string
	ChronosEndpoint      string
}

func DefaultConfig() *Config {
	return &Config{
		NetworkType:          "local",
		ChainId:              "injective-1",
		ExchangeGrpcEndpoint: "https://localhost:4444",
		ExplorerGrpcEndpoint: "https://localhost:9091",
		TmEndpoint:           "tcp://localhost:26657",
		ChainGrpcEndpoint:    "http://localhost:9900",
		LcdEndpoint:          "https://localhost:10337",
		Fee_denom:            "inj",
		ChronosEndpoint:      "https://localhost:4445",
	}
}

func DefaultNetwork() common.Network {
	network := common.NewNetwork()
	network.ChainId = "injective-1"
	network.ExchangeGrpcEndpoint = "tcp://localhost:9910"
	network.ExplorerGrpcEndpoint = "tcp://localhost:9091"
	network.TmEndpoint = "tcp://localhost:26657"
	network.ChainGrpcEndpoint = "tcp://localhost:9900"
	network.LcdEndpoint = "https://localhost:10337"
	network.Fee_denom = "inj"
	return network
}

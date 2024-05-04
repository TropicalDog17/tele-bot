package chain

import (
	"fmt"
	"os"

	chainclient "github.com/InjectiveLabs/sdk-go/client/chain"
	"github.com/InjectiveLabs/sdk-go/client/common"
	"github.com/TropicalDog17/orderbook-go-sdk/config"
	cosmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
)

const DefaultLocalGasPrice = "100000000000000inj"

type ChainClient struct {
	chainClient   chainclient.ChainClient
	SenderAddress cosmtypes.AccAddress
	CosmosKeyring keyring.Keyring
	ClientCtx     cosmclient.Context
}

// Get chain client with signing key prepared.
func NewChainClient(keyName string) ChainClient {
	network := config.DefaultNetwork()
	senderAddress, cosmosKeyring, err := chainclient.InitCosmosKeyring(
		os.Getenv("HOME")+"/.injectived",
		"injectived",
		"test",
		keyName,
		"12345678",
		"", // keyring will be used if pk not provided
		false,
	)
	if err != nil {
		panic(err)
	}
	clientCtx, err := chainclient.NewClientContext(
		"injective-1", // TODO: refactor hard code
		senderAddress.String(),
		cosmosKeyring,
	)
	fmt.Println("senderAddress: ", senderAddress.String())
	if err != nil {
		panic(err)
	}
	tmClient, err := rpchttp.New("http://localhost:26657", "/websocket")
	if err != nil {
		panic(err)
	}

	clientCtx = clientCtx.WithNodeURI("http://localhost:26657").WithClient(tmClient)
	chainClient, err := chainclient.NewChainClient(
		clientCtx,
		network,
		common.OptionGasPrices(DefaultLocalGasPrice),
	)
	if err != nil {
		panic(err)
	}
	return ChainClient{
		chainClient:   chainClient,
		SenderAddress: senderAddress,
		CosmosKeyring: cosmosKeyring,
		ClientCtx:     clientCtx,
	}
}

func (c *ChainClient) GetInjectiveChainClient() chainclient.ChainClient {
	return c.chainClient
}

func (c *ChainClient) AdjustKeyring(keyName string) {
	network := config.DefaultNetwork()
	senderAddress, cosmosKeyring, err := chainclient.InitCosmosKeyring(
		os.Getenv("HOME")+"/.injectived",
		"injectived",
		"test",
		keyName,
		"12345678",
		"", // keyring will be used if pk not provided
		false,
	)
	if err != nil {
		panic(err)
	}
	c.SenderAddress = senderAddress
	c.CosmosKeyring = cosmosKeyring
	clientCtx, err := chainclient.NewClientContext(
		"injective-1", // TODO: refactor hard code
		senderAddress.String(),
		cosmosKeyring,
	)
	if err != nil {
		panic(err)
	}
	tmClient, err := rpchttp.New("http://localhost:26657", "/websocket")
	if err != nil {
		panic(err)
	}

	clientCtx = clientCtx.WithNodeURI("http://localhost:26657").WithClient(tmClient)

	chainClient, err := chainclient.NewChainClient(
		clientCtx,
		network,
		common.OptionGasPrices(DefaultLocalGasPrice),
	)
	if err != nil {
		panic(err)
	}
	c.chainClient = chainClient
}

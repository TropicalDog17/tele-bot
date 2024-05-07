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

type ChainClient interface {
	GetInjectiveChainClient() chainclient.ChainClient
	AdjustKeyring(keyName string)
	AdjustKeyringFromPrivateKey(privateKey string)
	TransferToken(toAddress string, amount float64, denom string) (string, error)
	GetSenderAddress() cosmtypes.AccAddress
	GetBalance(address string, denom string) (float64, error)
}

type ChainClientStruct struct {
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
	return &ChainClientStruct{
		chainClient:   chainClient,
		SenderAddress: senderAddress,
		CosmosKeyring: cosmosKeyring,
		ClientCtx:     clientCtx,
	}
}

func (c *ChainClientStruct) GetInjectiveChainClient() chainclient.ChainClient {
	return c.chainClient
}

func NewChainClientFromPrivateKey(privateKey string) ChainClient {
	network := config.DefaultNetwork()
	senderAddress, cosmosKeyring, err := chainclient.InitCosmosKeyring(
		"",
		"",
		"",
		"",
		"",
		privateKey,
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
	return &ChainClientStruct{
		chainClient:   chainClient,
		SenderAddress: senderAddress,
		CosmosKeyring: cosmosKeyring,
		ClientCtx:     clientCtx,
	}
}

func (c *ChainClientStruct) AdjustKeyring(keyName string) {
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

func (c *ChainClientStruct) AdjustKeyringFromPrivateKey(privateKey string) {
	network := config.DefaultNetwork()
	senderAddress, cosmosKeyring, err := chainclient.InitCosmosKeyring(
		"",
		"",
		"",
		"",
		"",
		privateKey, // keyring will be used if pk not provided
		false,
	)
	if err != nil {
		panic(err)
	}
	// fund the account with some dust tokens
	_, err = c.TransferToken(senderAddress.String(), 0.00001, "inj")
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

func (c *ChainClientStruct) GetSenderAddress() cosmtypes.AccAddress {
	return c.SenderAddress
}

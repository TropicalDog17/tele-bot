package internal

import (
	"context"
	"time"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	spotExchangePB "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/chain"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/exchange"
	customsdktypes "github.com/TropicalDog17/orderbook-go-sdk/pkg/types"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/redis/go-redis/v9"
	"gopkg.in/telebot.v3"
)

type CoinGecko interface {
	FetchUsdPriceMap(denoms ...string) (map[string]float64, error)
	GetPriceInUsd(denoms ...string) (map[string]map[string]float64, error)
	GetAPIKey() string
}

type Bot interface {
	Delete(msg telebot.Editable) error
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
	Handle(endpoint interface{}, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc)
	ProcessUpdate(u telebot.Update)
}

type BotClient interface {
	GetPrice(ticker string) (float64, bool)
	SetPrice(ticker string, price float64)
	GetBalances(address string, denoms []string) (map[string]float64, error)
	TransferToken(to string, amount float64, denom string) (string, error)
	GetAddress() string
	GetDecimal(denom string) int32
	PlaceSpotOrder(denomIn, denomOut string, amount, price float64) (string, error)
	GetActiveOrders(marketId string) ([]types.LimitOrderInfo, error)
	CancelOrder(marketID, orderHash string) (string, error)
	ToMessage(order types.LimitOrderInfo, showDetail bool) string
	GetRedisInstance() RedisClient
	GetExchangeClient() *exchange.MbClient
	GetActiveMarkets() (map[string]string, error)
}

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
	SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	HKeys(ctx context.Context, key string) *redis.StringSliceCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
}

type ExchangeClient interface {
	CancelOrder(ctx context.Context, marketID string, orderID string) (string, error)
	GetChainClient() chain.ChainClient
	GetDecimals(ctx context.Context, marketId string) (baseDecimal int32, quoteDecimal int32)
	GetMarketSummary(marketId string) (customsdktypes.MarketSummary, error)
	GetMarketSummaryFromTicker(ticker string) (customsdktypes.MarketSummary, error)
	GetPrice(ticker string) (float64, error)
	GetSpotMarket(marketId string) (*exchangetypes.SpotMarket, error)
	GetSpotMarketFromTicker(ticker string) (*exchangetypes.SpotMarket, error)
	NewSpotOrder(orderType exchangetypes.OrderType, marketId string, price float64, quantity float64) customsdktypes.SpotOrder
	PlaceSpotOrder(order customsdktypes.SpotOrder) (string, error)
	GetActiveMarkets(ctx context.Context, req *spotExchangePB.MarketsRequest) ([]*spotExchangePB.SpotMarketInfo, error)
}

type TeleContext interface {
	// Bot returns the bot instance.
	Bot() *telebot.Bot

	// Update returns the original update.
	Update() telebot.Update

	// Message returns stored message if such presented.
	Message() *telebot.Message

	// Callback returns stored callback if such presented.
	Callback() *telebot.Callback

	// Query returns stored query if such presented.
	Query() *telebot.Query

	// InlineResult returns stored inline result if such presented.
	InlineResult() *telebot.InlineResult

	// ShippingQuery returns stored shipping query if such presented.
	ShippingQuery() *telebot.ShippingQuery

	// PreCheckoutQuery returns stored pre checkout query if such presented.
	PreCheckoutQuery() *telebot.PreCheckoutQuery

	// Poll returns stored poll if such presented.
	Poll() *telebot.Poll

	// PollAnswer returns stored poll answer if such presented.
	PollAnswer() *telebot.PollAnswer

	// ChatMember returns chat member changes.
	ChatMember() *telebot.ChatMemberUpdate

	// ChatJoinRequest returns the chat join request.
	ChatJoinRequest() *telebot.ChatJoinRequest

	// Migration returns both migration from and to chat IDs.
	Migration() (int64, int64)

	// Topic returns the topic changes.
	Topic() *telebot.Topic

	// Sender returns the current recipient, depending on the context type.
	// Returns nil if user is not presented.
	Sender() *telebot.User

	// Chat returns the current chat, depending on the context type.
	// Returns nil if chat is not presented.
	Chat() *telebot.Chat

	// Recipient combines both Sender and Chat functions. If there is no user
	// the chat will be returned. The native context cannot be without sender,
	// but it is useful in the case when the context created intentionally
	// by the NewContext constructor and have only Chat field inside.
	Recipient() telebot.Recipient

	// Text returns the message text, depending on the context type.
	// In the case when no related data presented, returns an empty string.
	Text() string

	// Entities returns the message entities, whether it's media caption's or the text's.
	// In the case when no entities presented, returns a nil.
	Entities() telebot.Entities

	// Data returns the current data, depending on the context type.
	// If the context contains command, returns its arguments string.
	// If the context contains payment, returns its payload.
	// In the case when no related data presented, returns an empty string.
	Data() string

	// Args returns a raw slice of command or callback arguments as strings.
	// The message arguments split by space, while the callback's ones by a "|" symbol.
	Args() []string

	// Send sends a message to the current recipient.
	// See Send from bot.go.
	Send(what interface{}, opts ...interface{}) error

	// SendAlbum sends an album to the current recipient.
	// See SendAlbum from bot.go.
	SendAlbum(a telebot.Album, opts ...interface{}) error

	// Reply replies to the current message.
	// See Reply from bot.go.
	Reply(what interface{}, opts ...interface{}) error

	// Forward forwards the given message to the current recipient.
	// See Forward from bot.go.
	Forward(msg telebot.Editable, opts ...interface{}) error

	// ForwardTo forwards the current message to the given recipient.
	// See Forward from bot.go
	ForwardTo(to telebot.Recipient, opts ...interface{}) error

	// Edit edits the current message.
	// See Edit from bot.go.
	Edit(what interface{}, opts ...interface{}) error

	// EditCaption edits the caption of the current message.
	// See EditCaption from bot.go.
	EditCaption(caption string, opts ...interface{}) error

	// EditOrSend edits the current message if the update is callback,
	// otherwise the content is sent to the chat as a separate message.
	EditOrSend(what interface{}, opts ...interface{}) error

	// EditOrReply edits the current message if the update is callback,
	// otherwise the content is replied as a separate message.
	EditOrReply(what interface{}, opts ...interface{}) error

	// Delete removes the current message.
	// See Delete from bot.go.
	Delete() error

	// DeleteAfter waits for the duration to elapse and then removes the
	// message. It handles an error automatically using b.OnError callback.
	// It returns a Timer that can be used to cancel the call using its Stop method.
	DeleteAfter(d time.Duration) *time.Timer

	// Notify updates the chat action for the current recipient.
	// See Notify from bot.go.
	Notify(action telebot.ChatAction) error

	// Ship replies to the current shipping query.
	// See Ship from bot.go.
	Ship(what ...interface{}) error

	// Accept finalizes the current deal.
	// See Accept from bot.go.
	Accept(errorMessage ...string) error

	// Answer sends a response to the current inline query.
	// See Answer from bot.go.
	Answer(resp *telebot.QueryResponse) error

	// Respond sends a response for the current callback query.
	// See Respond from bot.go.
	Respond(resp ...*telebot.CallbackResponse) error

	// Get retrieves data from the context.
	Get(key string) interface{}

	// Set saves data in the context.
	Set(key string, val interface{})
}

// Code generated by MockGen. DO NOT EDIT.
// Source: internal/interface.go
//
// Generated by this command:
//
//	mockgen -source=internal/interface.go -destination tests/mocks/interface.go
//

// Package mock_internal is a generated GoMock package.
package mock_internal

import (
	context "context"
	reflect "reflect"
	time "time"

	types "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	injective_spot_exchange_rpcpb "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	chain "github.com/TropicalDog17/orderbook-go-sdk/pkg/chain"
	types0 "github.com/TropicalDog17/orderbook-go-sdk/pkg/types"
	internal "github.com/TropicalDog17/tele-bot/internal"
	types1 "github.com/TropicalDog17/tele-bot/internal/types"
	redis "github.com/redis/go-redis/v9"
	gomock "go.uber.org/mock/gomock"
	telebot "gopkg.in/telebot.v3"
)

// MockCoinGecko is a mock of CoinGecko interface.
type MockCoinGecko struct {
	ctrl     *gomock.Controller
	recorder *MockCoinGeckoMockRecorder
}

// MockCoinGeckoMockRecorder is the mock recorder for MockCoinGecko.
type MockCoinGeckoMockRecorder struct {
	mock *MockCoinGecko
}

// NewMockCoinGecko creates a new mock instance.
func NewMockCoinGecko(ctrl *gomock.Controller) *MockCoinGecko {
	mock := &MockCoinGecko{ctrl: ctrl}
	mock.recorder = &MockCoinGeckoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoinGecko) EXPECT() *MockCoinGeckoMockRecorder {
	return m.recorder
}

// FetchUsdPriceMap mocks base method.
func (m *MockCoinGecko) FetchUsdPriceMap(denoms ...string) (map[string]float64, error) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range denoms {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FetchUsdPriceMap", varargs...)
	ret0, _ := ret[0].(map[string]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchUsdPriceMap indicates an expected call of FetchUsdPriceMap.
func (mr *MockCoinGeckoMockRecorder) FetchUsdPriceMap(denoms ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchUsdPriceMap", reflect.TypeOf((*MockCoinGecko)(nil).FetchUsdPriceMap), denoms...)
}

// GetAPIKey mocks base method.
func (m *MockCoinGecko) GetAPIKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAPIKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAPIKey indicates an expected call of GetAPIKey.
func (mr *MockCoinGeckoMockRecorder) GetAPIKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAPIKey", reflect.TypeOf((*MockCoinGecko)(nil).GetAPIKey))
}

// GetPriceInUsd mocks base method.
func (m *MockCoinGecko) GetPriceInUsd(denoms ...string) (map[string]map[string]float64, error) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range denoms {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPriceInUsd", varargs...)
	ret0, _ := ret[0].(map[string]map[string]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPriceInUsd indicates an expected call of GetPriceInUsd.
func (mr *MockCoinGeckoMockRecorder) GetPriceInUsd(denoms ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriceInUsd", reflect.TypeOf((*MockCoinGecko)(nil).GetPriceInUsd), denoms...)
}

// MockBot is a mock of Bot interface.
type MockBot struct {
	ctrl     *gomock.Controller
	recorder *MockBotMockRecorder
}

// MockBotMockRecorder is the mock recorder for MockBot.
type MockBotMockRecorder struct {
	mock *MockBot
}

// NewMockBot creates a new mock instance.
func NewMockBot(ctrl *gomock.Controller) *MockBot {
	mock := &MockBot{ctrl: ctrl}
	mock.recorder = &MockBotMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBot) EXPECT() *MockBotMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockBot) Delete(msg telebot.Editable) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockBotMockRecorder) Delete(msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBot)(nil).Delete), msg)
}

// Handle mocks base method.
func (m_2 *MockBot) Handle(endpoint any, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc) {
	m_2.ctrl.T.Helper()
	varargs := []any{endpoint, h}
	for _, a := range m {
		varargs = append(varargs, a)
	}
	m_2.ctrl.Call(m_2, "Handle", varargs...)
}

// Handle indicates an expected call of Handle.
func (mr *MockBotMockRecorder) Handle(endpoint, h any, m ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{endpoint, h}, m...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockBot)(nil).Handle), varargs...)
}

// ProcessUpdate mocks base method.
func (m *MockBot) ProcessUpdate(u telebot.Update) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ProcessUpdate", u)
}

// ProcessUpdate indicates an expected call of ProcessUpdate.
func (mr *MockBotMockRecorder) ProcessUpdate(u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessUpdate", reflect.TypeOf((*MockBot)(nil).ProcessUpdate), u)
}

// Send mocks base method.
func (m *MockBot) Send(to telebot.Recipient, what any, opts ...any) (*telebot.Message, error) {
	m.ctrl.T.Helper()
	varargs := []any{to, what}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Send", varargs...)
	ret0, _ := ret[0].(*telebot.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send.
func (mr *MockBotMockRecorder) Send(to, what any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{to, what}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockBot)(nil).Send), varargs...)
}

// MockBotClient is a mock of BotClient interface.
type MockBotClient struct {
	ctrl     *gomock.Controller
	recorder *MockBotClientMockRecorder
}

// MockBotClientMockRecorder is the mock recorder for MockBotClient.
type MockBotClientMockRecorder struct {
	mock *MockBotClient
}

// NewMockBotClient creates a new mock instance.
func NewMockBotClient(ctrl *gomock.Controller) *MockBotClient {
	mock := &MockBotClient{ctrl: ctrl}
	mock.recorder = &MockBotClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBotClient) EXPECT() *MockBotClientMockRecorder {
	return m.recorder
}

// CancelOrder mocks base method.
func (m *MockBotClient) CancelOrder(marketID, orderHash string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelOrder", marketID, orderHash)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CancelOrder indicates an expected call of CancelOrder.
func (mr *MockBotClientMockRecorder) CancelOrder(marketID, orderHash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelOrder", reflect.TypeOf((*MockBotClient)(nil).CancelOrder), marketID, orderHash)
}

// GetActiveMarkets mocks base method.
func (m *MockBotClient) GetActiveMarkets() (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveMarkets")
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveMarkets indicates an expected call of GetActiveMarkets.
func (mr *MockBotClientMockRecorder) GetActiveMarkets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveMarkets", reflect.TypeOf((*MockBotClient)(nil).GetActiveMarkets))
}

// GetActiveOrders mocks base method.
func (m *MockBotClient) GetActiveOrders(marketId string) ([]types1.LimitOrderInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveOrders", marketId)
	ret0, _ := ret[0].([]types1.LimitOrderInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveOrders indicates an expected call of GetActiveOrders.
func (mr *MockBotClientMockRecorder) GetActiveOrders(marketId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveOrders", reflect.TypeOf((*MockBotClient)(nil).GetActiveOrders), marketId)
}

// GetAddress mocks base method.
func (m *MockBotClient) GetAddress() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddress")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAddress indicates an expected call of GetAddress.
func (mr *MockBotClientMockRecorder) GetAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddress", reflect.TypeOf((*MockBotClient)(nil).GetAddress))
}

// GetBalances mocks base method.
func (m *MockBotClient) GetBalances(address string, denoms []string) (map[string]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalances", address, denoms)
	ret0, _ := ret[0].(map[string]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalances indicates an expected call of GetBalances.
func (mr *MockBotClientMockRecorder) GetBalances(address, denoms any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalances", reflect.TypeOf((*MockBotClient)(nil).GetBalances), address, denoms)
}

// GetDecimal mocks base method.
func (m *MockBotClient) GetDecimal(denom string) int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDecimal", denom)
	ret0, _ := ret[0].(int32)
	return ret0
}

// GetDecimal indicates an expected call of GetDecimal.
func (mr *MockBotClientMockRecorder) GetDecimal(denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDecimal", reflect.TypeOf((*MockBotClient)(nil).GetDecimal), denom)
}

// GetPrice mocks base method.
func (m *MockBotClient) GetPrice(ticker string) (float64, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrice", ticker)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetPrice indicates an expected call of GetPrice.
func (mr *MockBotClientMockRecorder) GetPrice(ticker any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrice", reflect.TypeOf((*MockBotClient)(nil).GetPrice), ticker)
}

// GetRedisInstance mocks base method.
func (m *MockBotClient) GetRedisInstance() internal.RedisClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRedisInstance")
	ret0, _ := ret[0].(internal.RedisClient)
	return ret0
}

// GetRedisInstance indicates an expected call of GetRedisInstance.
func (mr *MockBotClientMockRecorder) GetRedisInstance() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRedisInstance", reflect.TypeOf((*MockBotClient)(nil).GetRedisInstance))
}

// PlaceSpotOrder mocks base method.
func (m *MockBotClient) PlaceSpotOrder(denomIn, denomOut string, amount, price float64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlaceSpotOrder", denomIn, denomOut, amount, price)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlaceSpotOrder indicates an expected call of PlaceSpotOrder.
func (mr *MockBotClientMockRecorder) PlaceSpotOrder(denomIn, denomOut, amount, price any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlaceSpotOrder", reflect.TypeOf((*MockBotClient)(nil).PlaceSpotOrder), denomIn, denomOut, amount, price)
}

// SetPrice mocks base method.
func (m *MockBotClient) SetPrice(ticker string, price float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPrice", ticker, price)
}

// SetPrice indicates an expected call of SetPrice.
func (mr *MockBotClientMockRecorder) SetPrice(ticker, price any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPrice", reflect.TypeOf((*MockBotClient)(nil).SetPrice), ticker, price)
}

// ToMessage mocks base method.
func (m *MockBotClient) ToMessage(order types1.LimitOrderInfo, showDetail bool) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToMessage", order, showDetail)
	ret0, _ := ret[0].(string)
	return ret0
}

// ToMessage indicates an expected call of ToMessage.
func (mr *MockBotClientMockRecorder) ToMessage(order, showDetail any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToMessage", reflect.TypeOf((*MockBotClient)(nil).ToMessage), order, showDetail)
}

// TransferToken mocks base method.
func (m *MockBotClient) TransferToken(to string, amount float64, denom string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferToken", to, amount, denom)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransferToken indicates an expected call of TransferToken.
func (mr *MockBotClientMockRecorder) TransferToken(to, amount, denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferToken", reflect.TypeOf((*MockBotClient)(nil).TransferToken), to, amount, denom)
}

// MockRedisClient is a mock of RedisClient interface.
type MockRedisClient struct {
	ctrl     *gomock.Controller
	recorder *MockRedisClientMockRecorder
}

// MockRedisClientMockRecorder is the mock recorder for MockRedisClient.
type MockRedisClientMockRecorder struct {
	mock *MockRedisClient
}

// NewMockRedisClient creates a new mock instance.
func NewMockRedisClient(ctrl *gomock.Controller) *MockRedisClient {
	mock := &MockRedisClient{ctrl: ctrl}
	mock.recorder = &MockRedisClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedisClient) EXPECT() *MockRedisClientMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(*redis.StringCmd)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockRedisClientMockRecorder) Get(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedisClient)(nil).Get), ctx, key)
}

// HDel mocks base method.
func (m *MockRedisClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	m.ctrl.T.Helper()
	varargs := []any{ctx, key}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "HDel", varargs...)
	ret0, _ := ret[0].(*redis.IntCmd)
	return ret0
}

// HDel indicates an expected call of HDel.
func (mr *MockRedisClientMockRecorder) HDel(ctx, key any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, key}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HDel", reflect.TypeOf((*MockRedisClient)(nil).HDel), varargs...)
}

// HGet mocks base method.
func (m *MockRedisClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HGet", ctx, key, field)
	ret0, _ := ret[0].(*redis.StringCmd)
	return ret0
}

// HGet indicates an expected call of HGet.
func (mr *MockRedisClientMockRecorder) HGet(ctx, key, field any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HGet", reflect.TypeOf((*MockRedisClient)(nil).HGet), ctx, key, field)
}

// HGetAll mocks base method.
func (m *MockRedisClient) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HGetAll", ctx, key)
	ret0, _ := ret[0].(*redis.MapStringStringCmd)
	return ret0
}

// HGetAll indicates an expected call of HGetAll.
func (mr *MockRedisClientMockRecorder) HGetAll(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HGetAll", reflect.TypeOf((*MockRedisClient)(nil).HGetAll), ctx, key)
}

// HKeys mocks base method.
func (m *MockRedisClient) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HKeys", ctx, key)
	ret0, _ := ret[0].(*redis.StringSliceCmd)
	return ret0
}

// HKeys indicates an expected call of HKeys.
func (mr *MockRedisClientMockRecorder) HKeys(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HKeys", reflect.TypeOf((*MockRedisClient)(nil).HKeys), ctx, key)
}

// HSet mocks base method.
func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...any) *redis.IntCmd {
	m.ctrl.T.Helper()
	varargs := []any{ctx, key}
	for _, a := range values {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "HSet", varargs...)
	ret0, _ := ret[0].(*redis.IntCmd)
	return ret0
}

// HSet indicates an expected call of HSet.
func (mr *MockRedisClientMockRecorder) HSet(ctx, key any, values ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, key}, values...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HSet", reflect.TypeOf((*MockRedisClient)(nil).HSet), varargs...)
}

// SAdd mocks base method.
func (m *MockRedisClient) SAdd(ctx context.Context, key string, members ...any) *redis.IntCmd {
	m.ctrl.T.Helper()
	varargs := []any{ctx, key}
	for _, a := range members {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SAdd", varargs...)
	ret0, _ := ret[0].(*redis.IntCmd)
	return ret0
}

// SAdd indicates an expected call of SAdd.
func (mr *MockRedisClientMockRecorder) SAdd(ctx, key any, members ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, key}, members...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SAdd", reflect.TypeOf((*MockRedisClient)(nil).SAdd), varargs...)
}

// SMembers mocks base method.
func (m *MockRedisClient) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SMembers", ctx, key)
	ret0, _ := ret[0].(*redis.StringSliceCmd)
	return ret0
}

// SMembers indicates an expected call of SMembers.
func (mr *MockRedisClientMockRecorder) SMembers(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SMembers", reflect.TypeOf((*MockRedisClient)(nil).SMembers), ctx, key)
}

// SRem mocks base method.
func (m *MockRedisClient) SRem(ctx context.Context, key string, members ...any) *redis.IntCmd {
	m.ctrl.T.Helper()
	varargs := []any{ctx, key}
	for _, a := range members {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SRem", varargs...)
	ret0, _ := ret[0].(*redis.IntCmd)
	return ret0
}

// SRem indicates an expected call of SRem.
func (mr *MockRedisClientMockRecorder) SRem(ctx, key any, members ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, key}, members...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SRem", reflect.TypeOf((*MockRedisClient)(nil).SRem), varargs...)
}

// Set mocks base method.
func (m *MockRedisClient) Set(arg0 context.Context, arg1 string, arg2 any, arg3 time.Duration) *redis.StatusCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*redis.StatusCmd)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockRedisClientMockRecorder) Set(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRedisClient)(nil).Set), arg0, arg1, arg2, arg3)
}

// MockExchangeClient is a mock of ExchangeClient interface.
type MockExchangeClient struct {
	ctrl     *gomock.Controller
	recorder *MockExchangeClientMockRecorder
}

// MockExchangeClientMockRecorder is the mock recorder for MockExchangeClient.
type MockExchangeClientMockRecorder struct {
	mock *MockExchangeClient
}

// NewMockExchangeClient creates a new mock instance.
func NewMockExchangeClient(ctrl *gomock.Controller) *MockExchangeClient {
	mock := &MockExchangeClient{ctrl: ctrl}
	mock.recorder = &MockExchangeClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExchangeClient) EXPECT() *MockExchangeClientMockRecorder {
	return m.recorder
}

// CancelOrder mocks base method.
func (m *MockExchangeClient) CancelOrder(ctx context.Context, marketID, orderID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelOrder", ctx, marketID, orderID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CancelOrder indicates an expected call of CancelOrder.
func (mr *MockExchangeClientMockRecorder) CancelOrder(ctx, marketID, orderID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelOrder", reflect.TypeOf((*MockExchangeClient)(nil).CancelOrder), ctx, marketID, orderID)
}

// GetActiveMarkets mocks base method.
func (m *MockExchangeClient) GetActiveMarkets(ctx context.Context, req *injective_spot_exchange_rpcpb.MarketsRequest) ([]*injective_spot_exchange_rpcpb.SpotMarketInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveMarkets", ctx, req)
	ret0, _ := ret[0].([]*injective_spot_exchange_rpcpb.SpotMarketInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveMarkets indicates an expected call of GetActiveMarkets.
func (mr *MockExchangeClientMockRecorder) GetActiveMarkets(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveMarkets", reflect.TypeOf((*MockExchangeClient)(nil).GetActiveMarkets), ctx, req)
}

// GetChainClient mocks base method.
func (m *MockExchangeClient) GetChainClient() *chain.ChainClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChainClient")
	ret0, _ := ret[0].(*chain.ChainClient)
	return ret0
}

// GetChainClient indicates an expected call of GetChainClient.
func (mr *MockExchangeClientMockRecorder) GetChainClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChainClient", reflect.TypeOf((*MockExchangeClient)(nil).GetChainClient))
}

// GetDecimals mocks base method.
func (m *MockExchangeClient) GetDecimals(ctx context.Context, marketId string) (int32, int32) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDecimals", ctx, marketId)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(int32)
	return ret0, ret1
}

// GetDecimals indicates an expected call of GetDecimals.
func (mr *MockExchangeClientMockRecorder) GetDecimals(ctx, marketId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDecimals", reflect.TypeOf((*MockExchangeClient)(nil).GetDecimals), ctx, marketId)
}

// GetMarketSummary mocks base method.
func (m *MockExchangeClient) GetMarketSummary(marketId string) (types0.MarketSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMarketSummary", marketId)
	ret0, _ := ret[0].(types0.MarketSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMarketSummary indicates an expected call of GetMarketSummary.
func (mr *MockExchangeClientMockRecorder) GetMarketSummary(marketId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMarketSummary", reflect.TypeOf((*MockExchangeClient)(nil).GetMarketSummary), marketId)
}

// GetMarketSummaryFromTicker mocks base method.
func (m *MockExchangeClient) GetMarketSummaryFromTicker(ticker string) (types0.MarketSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMarketSummaryFromTicker", ticker)
	ret0, _ := ret[0].(types0.MarketSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMarketSummaryFromTicker indicates an expected call of GetMarketSummaryFromTicker.
func (mr *MockExchangeClientMockRecorder) GetMarketSummaryFromTicker(ticker any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMarketSummaryFromTicker", reflect.TypeOf((*MockExchangeClient)(nil).GetMarketSummaryFromTicker), ticker)
}

// GetPrice mocks base method.
func (m *MockExchangeClient) GetPrice(ticker string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrice", ticker)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPrice indicates an expected call of GetPrice.
func (mr *MockExchangeClientMockRecorder) GetPrice(ticker any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrice", reflect.TypeOf((*MockExchangeClient)(nil).GetPrice), ticker)
}

// GetSpotMarket mocks base method.
func (m *MockExchangeClient) GetSpotMarket(marketId string) (*types.SpotMarket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpotMarket", marketId)
	ret0, _ := ret[0].(*types.SpotMarket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpotMarket indicates an expected call of GetSpotMarket.
func (mr *MockExchangeClientMockRecorder) GetSpotMarket(marketId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpotMarket", reflect.TypeOf((*MockExchangeClient)(nil).GetSpotMarket), marketId)
}

// GetSpotMarketFromTicker mocks base method.
func (m *MockExchangeClient) GetSpotMarketFromTicker(ticker string) (*types.SpotMarket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpotMarketFromTicker", ticker)
	ret0, _ := ret[0].(*types.SpotMarket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpotMarketFromTicker indicates an expected call of GetSpotMarketFromTicker.
func (mr *MockExchangeClientMockRecorder) GetSpotMarketFromTicker(ticker any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpotMarketFromTicker", reflect.TypeOf((*MockExchangeClient)(nil).GetSpotMarketFromTicker), ticker)
}

// NewSpotOrder mocks base method.
func (m *MockExchangeClient) NewSpotOrder(orderType types.OrderType, marketId string, price, quantity float64) types0.SpotOrder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSpotOrder", orderType, marketId, price, quantity)
	ret0, _ := ret[0].(types0.SpotOrder)
	return ret0
}

// NewSpotOrder indicates an expected call of NewSpotOrder.
func (mr *MockExchangeClientMockRecorder) NewSpotOrder(orderType, marketId, price, quantity any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSpotOrder", reflect.TypeOf((*MockExchangeClient)(nil).NewSpotOrder), orderType, marketId, price, quantity)
}

// PlaceSpotOrder mocks base method.
func (m *MockExchangeClient) PlaceSpotOrder(order types0.SpotOrder) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlaceSpotOrder", order)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlaceSpotOrder indicates an expected call of PlaceSpotOrder.
func (mr *MockExchangeClientMockRecorder) PlaceSpotOrder(order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlaceSpotOrder", reflect.TypeOf((*MockExchangeClient)(nil).PlaceSpotOrder), order)
}

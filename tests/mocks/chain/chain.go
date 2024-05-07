// Code generated by MockGen. DO NOT EDIT.
// Source: go-sdk/pkg/chain/types.go
//
// Generated by this command:
//
//	mockgen -source=go-sdk/pkg/chain/types.go -destination tests/mocks/chain/chain.go
//

// Package mock_chain is a generated GoMock package.
package mock_chain

import (
	reflect "reflect"

	chain "github.com/InjectiveLabs/sdk-go/client/chain"
	types "github.com/cosmos/cosmos-sdk/types"
	gomock "go.uber.org/mock/gomock"
)

// MockChainClient is a mock of ChainClient interface.
type MockChainClient struct {
	ctrl     *gomock.Controller
	recorder *MockChainClientMockRecorder
}

// MockChainClientMockRecorder is the mock recorder for MockChainClient.
type MockChainClientMockRecorder struct {
	mock *MockChainClient
}

// NewMockChainClient creates a new mock instance.
func NewMockChainClient(ctrl *gomock.Controller) *MockChainClient {
	mock := &MockChainClient{ctrl: ctrl}
	mock.recorder = &MockChainClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChainClient) EXPECT() *MockChainClientMockRecorder {
	return m.recorder
}

// AdjustKeyring mocks base method.
func (m *MockChainClient) AdjustKeyring(keyName string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AdjustKeyring", keyName)
}

// AdjustKeyring indicates an expected call of AdjustKeyring.
func (mr *MockChainClientMockRecorder) AdjustKeyring(keyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdjustKeyring", reflect.TypeOf((*MockChainClient)(nil).AdjustKeyring), keyName)
}

// AdjustKeyringFromPrivateKey mocks base method.
func (m *MockChainClient) AdjustKeyringFromPrivateKey(privateKey string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AdjustKeyringFromPrivateKey", privateKey)
}

// AdjustKeyringFromPrivateKey indicates an expected call of AdjustKeyringFromPrivateKey.
func (mr *MockChainClientMockRecorder) AdjustKeyringFromPrivateKey(privateKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdjustKeyringFromPrivateKey", reflect.TypeOf((*MockChainClient)(nil).AdjustKeyringFromPrivateKey), privateKey)
}

// GetBalance mocks base method.
func (m *MockChainClient) GetBalance(address, denom string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", address, denom)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockChainClientMockRecorder) GetBalance(address, denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockChainClient)(nil).GetBalance), address, denom)
}

// GetInjectiveChainClient mocks base method.
func (m *MockChainClient) GetInjectiveChainClient() chain.ChainClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInjectiveChainClient")
	ret0, _ := ret[0].(chain.ChainClient)
	return ret0
}

// GetInjectiveChainClient indicates an expected call of GetInjectiveChainClient.
func (mr *MockChainClientMockRecorder) GetInjectiveChainClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInjectiveChainClient", reflect.TypeOf((*MockChainClient)(nil).GetInjectiveChainClient))
}

// GetSenderAddress mocks base method.
func (m *MockChainClient) GetSenderAddress() types.AccAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSenderAddress")
	ret0, _ := ret[0].(types.AccAddress)
	return ret0
}

// GetSenderAddress indicates an expected call of GetSenderAddress.
func (mr *MockChainClientMockRecorder) GetSenderAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSenderAddress", reflect.TypeOf((*MockChainClient)(nil).GetSenderAddress))
}

// TransferToken mocks base method.
func (m *MockChainClient) TransferToken(toAddress string, amount float64, denom string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferToken", toAddress, amount, denom)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransferToken indicates an expected call of TransferToken.
func (mr *MockChainClientMockRecorder) TransferToken(toAddress, amount, denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferToken", reflect.TypeOf((*MockChainClient)(nil).TransferToken), toAddress, amount, denom)
}

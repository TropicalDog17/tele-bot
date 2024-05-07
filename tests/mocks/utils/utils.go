// Code generated by MockGen. DO NOT EDIT.
// Source: internal/utils/utils.go
//
// Generated by this command:
//
//	mockgen -source=internal/utils/utils.go -destination tests/mocks/utils/utils.go
//

// Package mock_utils is a generated GoMock package.
package mock_utils

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUtilsInterface is a mock of UtilsInterface interface.
type MockUtilsInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUtilsInterfaceMockRecorder
}

// MockUtilsInterfaceMockRecorder is the mock recorder for MockUtilsInterface.
type MockUtilsInterfaceMockRecorder struct {
	mock *MockUtilsInterface
}

// NewMockUtilsInterface creates a new mock instance.
func NewMockUtilsInterface(ctrl *gomock.Controller) *MockUtilsInterface {
	mock := &MockUtilsInterface{ctrl: ctrl}
	mock.recorder = &MockUtilsInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUtilsInterface) EXPECT() *MockUtilsInterfaceMockRecorder {
	return m.recorder
}

// GenerateMnemonic mocks base method.
func (m *MockUtilsInterface) GenerateMnemonic() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateMnemonic")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateMnemonic indicates an expected call of GenerateMnemonic.
func (mr *MockUtilsInterfaceMockRecorder) GenerateMnemonic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateMnemonic", reflect.TypeOf((*MockUtilsInterface)(nil).GenerateMnemonic))
}

// GetEncryptedMnemonic mocks base method.
func (m *MockUtilsInterface) GetEncryptedMnemonic(mnemonic, password string) (string, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEncryptedMnemonic", mnemonic, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetEncryptedMnemonic indicates an expected call of GetEncryptedMnemonic.
func (mr *MockUtilsInterfaceMockRecorder) GetEncryptedMnemonic(mnemonic, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEncryptedMnemonic", reflect.TypeOf((*MockUtilsInterface)(nil).GetEncryptedMnemonic), mnemonic, password)
}

// MnemonicChallenge mocks base method.
func (m *MockUtilsInterface) MnemonicChallenge(mnemonic string, indexes [3]int, providedWords [3]string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MnemonicChallenge", mnemonic, indexes, providedWords)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MnemonicChallenge indicates an expected call of MnemonicChallenge.
func (mr *MockUtilsInterfaceMockRecorder) MnemonicChallenge(mnemonic, indexes, providedWords any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MnemonicChallenge", reflect.TypeOf((*MockUtilsInterface)(nil).MnemonicChallenge), mnemonic, indexes, providedWords)
}

// SplitMnemonic mocks base method.
func (m *MockUtilsInterface) SplitMnemonic(mnemonic string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SplitMnemonic", mnemonic)
	ret0, _ := ret[0].([]string)
	return ret0
}

// SplitMnemonic indicates an expected call of SplitMnemonic.
func (mr *MockUtilsInterfaceMockRecorder) SplitMnemonic(mnemonic any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SplitMnemonic", reflect.TypeOf((*MockUtilsInterface)(nil).SplitMnemonic), mnemonic)
}

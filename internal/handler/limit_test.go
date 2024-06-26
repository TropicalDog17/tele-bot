package handler_test

import (
	"strconv"
	"testing"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/handler"
	"github.com/TropicalDog17/tele-bot/internal/types"
	mock_internal "github.com/TropicalDog17/tele-bot/tests/mocks"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v3"
)

func TestHandleLimitOrder_BtnLimitOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockBot := mock_internal.NewMockBot(ctrl)
	mockBotClient := mock_internal.NewMockBotClient(ctrl)
	mockClients := make(map[string]internal.BotClient)
	mockClients["sender"] = mockBotClient

	mockStoredMessages := GetStoredMessagesForTest(10)
	mockLimitOrderInfo := GetMockLimitOrderInfo()
	mockCurrentStep := "limitAmount"

	// Create the handler instance
	mockBot.EXPECT().Handle(gomock.Any(), gomock.Any()).Times(10)
	handler.HandleLimitOrder(mockBot, &tele.Group{}, mockClients, &mockStoredMessages[0], &mockStoredMessages[1], &mockLimitOrderInfo, &mockCurrentStep, GetMockReplyMarkup(), GetMockReplyMarkup(), GetMockReplyMarkup(), GetMockReplyMarkup())

	// Call the function

}

func GetStoredMessagesForTest(amount int) []tele.StoredMessage {
	storedMessages := make([]tele.StoredMessage, amount)
	for i := 0; i < amount; i++ {
		storedMessages[i] = tele.StoredMessage{
			ChatID:    0,
			MessageID: strconv.Itoa(i),
		}
	}
	return storedMessages
}

func GetMockLimitOrderInfo() types.LimitOrderInfo {
	return *types.NewLimitOrderInfo()
}

func GetMockReplyMarkup() *tele.ReplyMarkup {
	return &tele.ReplyMarkup{}
}

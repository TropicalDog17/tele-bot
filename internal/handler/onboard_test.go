package handler

import (
	"testing"

	"github.com/TropicalDog17/tele-bot/internal"
	mock_internal "github.com/TropicalDog17/tele-bot/tests/mocks"
	mock_utils "github.com/TropicalDog17/tele-bot/tests/mocks/utils"
	"github.com/awnumar/memguard"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v3"
)

func TestHandleOnboard(t *testing.T) {
	// Create mock bot and client
	ctrl := gomock.NewController(t)
	b := mock_internal.NewMockBot(ctrl)
	mockSender := &telebot.User{
		ID:       1,
		Username: "sender",
	}
	client := mock_internal.NewMockBotClient(ctrl)
	ctxMock := mock_internal.NewMockTeleContext(ctrl)
	currentStep := "start"
	clients := make(map[string]internal.BotClient)
	clients[mockSender.Username] = client
	ctxMock.EXPECT().Reply(gomock.Any(), gomock.Any()).MinTimes(1).Return(nil)
	// ctxMock.EXPECT().Message().MinTimes(1).Return(&telebot.Message{Text: "/start"})
	b.EXPECT().Handle("/start", gomock.Any()).MinTimes(1)

	HandleOnboard(b, clients, &currentStep)
	err := HandleStart(ctxMock, client, &currentStep)
	require.NoError(t, err)

	require.Equal(t, "addPassword", currentStep)
}

func TestHandleStorePrivateKey(t *testing.T) {
	testMnemonic := memguard.NewBufferFromBytes([]byte("test mnemonic"))
	testPassword := memguard.NewBufferFromBytes([]byte("test password"))
	encryptedMnemonic := "encryptedMnemonic"
	salt := []byte("astrongsalt")
	step := "sendMnemonic"
	mockUser := &telebot.User{
		ID:       1,
		Username: "sender",
	}
	mockMessage := &telebot.Message{
		ID:     1,
		Text:   "test password",
		Sender: mockUser,
	}
	mockChat := &telebot.Chat{
		ID: 1,
	}
	// Create mock bot and client
	ctrl := gomock.NewController(t)
	b := mock_internal.NewMockBot(ctrl)
	client := mock_internal.NewMockBotClient(ctrl)
	utilsMock := mock_utils.NewMockUtilsInterface(ctrl)
	ctxMock := mock_internal.NewMockTeleContext(ctrl)
	rdb, mock := redismock.NewClientMock()

	client.EXPECT().GetRedisInstance().Return(rdb)
	utilsMock.EXPECT().GetEncryptedMnemonic(gomock.Any(), gomock.Any()).Return(encryptedMnemonic, salt, nil)
	ctxMock.EXPECT().Message().Return(mockMessage).AnyTimes()
	mock.ExpectHSet(mockUser.Username, "encryptedMnemonic", encryptedMnemonic).SetVal(1)
	mock.ExpectHSet(mockUser.Username, "salt", salt).SetVal(1)
	ctxMock.EXPECT().Chat().Return(mockChat).AnyTimes()
	ctxMock.EXPECT().Send(gomock.Any(), gomock.Any()).MinTimes(1).Return(nil)
	b.EXPECT().Send(mockChat, gomock.Any()).MinTimes(1).Return(mockMessage, nil)

	msg, err := HandleStorePrivateKey(b, ctxMock, client, utilsMock, testMnemonic, testPassword, &step)
	require.NoError(t, err)
	require.Equal(t, 1, msg.ID)
}

func TestHandleConfirmMnemonicStep(t *testing.T) {
	// Create mock bot, client, and context
	b := mock_internal.NewMockBot(gomock.NewController(t))
	c := mock_internal.NewMockTeleContext(gomock.NewController(t))
	botClient := mock_internal.NewMockBotClient(gomock.NewController(t))
	step := "confirmMnemonic"

	// Set test mnemonic with 24 words
	testMnemonic := memguard.NewBufferFromBytes([]byte("word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12 word13 word14 word15 word16 word17 word18 word19 word20 word21 word22 word23 word24"))

	c.EXPECT().Reply(gomock.Any(), gomock.Any()).MinTimes(1).Return(nil)
	err := HandleConfirmMnemonicStep(b, c, botClient, testMnemonic, &step)
	require.NoError(t, err)

	// Verify that step is updated
	require.Equal(t, "receiveMnemonicWords", step)

}

// func TestHandleReceiveMnemonicWords(t *testing.T) {
// 	// Create mock bot, client, and context
// 	b := mock_internal.NewMockBot(gomock.NewController(t))
// 	c := mock_internal.NewMockTeleContext(gomock.NewController(t))
// 	client := mock_internal.NewMockBotClient(gomock.NewController(t))
// 	mockUtils := mock_utils.NewMockUtilsInterface(gomock.NewController(t))
// 	step := "receiveMnemonicWords"

// 	// Set test mnemonic and random indexes
// 	testMnemonic := "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12 word13 word14 word15 word16 word17 word18 word19 word20 word21 word22 word23 word24"
// 	c.EXPECT().Text().Return("word3 word6 word9")
// 	c.EXPECT().Send(gomock.Any(), gomock.Any()).MinTimes(1).Return(nil)
// 	mockUtils.EXPECT().SplitMnemonic(gomock.Any()).Return([]string{"word3", "word6", "word9"})
// 	mockUtils.EXPECT().MnemonicChallenge(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
// 	randomIndexes = [3]int{3, 6, 9}

// 	HandleReceiveMnemonicWords(b, c, client, mockUtils, testMnemonic, randomIndexes, &step)

// 	// Verify that mnemonic is confirmed and not back to the previous step
// 	require.NotEqual(t, "confirmMnemonic", step)
// }

// func TestAfterMnemonicConfirmed(t *testing.T) {
// 	// Create mock bot, client, and context
// 	b := mock_internal.NewMockBot(gomock.NewController(t))
// 	c := mock_internal.NewMockTeleContext(gomock.NewController(t))
// 	// mockUtils := mock_utils.NewMockUtilsInterface(gomock.NewController(t))
// 	exchangeClient := mock_internal.NewMockExchangeClient(gomock.NewController(t))
// 	chainClient := mock_chain.NewMockChainClient(gomock.NewController(t))
// 	step := "confirmMnemonic"

// 	// Set test mnemonic with 24 words
// 	testMnemonic := memguard.NewBufferFromBytes([]byte("pony glide frown crisp unfold lawn cup loan trial govern usual matrix theory wash fresh address pioneer between meadow visa buffalo keep gallery swear"))
// 	exchangeClient.EXPECT().GetChainClient().Return(chainClient)
// 	chainClient.EXPECT().AdjustKeyringFromPrivateKey(gomock.Any()).MinTimes(1)
// 	c.EXPECT().Reply(gomock.Any()).MinTimes(1).Return(nil)
// 	err := AfterMnemonicConfirmed(b, c, exchangeClient, testMnemonic, &step)
// 	require.NoError(t, err)
// 	// Verify that step is updated
// 	require.Equal(t, "", step)
// }

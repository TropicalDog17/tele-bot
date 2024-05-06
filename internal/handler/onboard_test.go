package handler

import (
	"testing"

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
	client := mock_internal.NewMockBotClient(ctrl)
	ctxMock := mock_internal.NewMockTeleContext(ctrl)
	currentStep := "start"

	ctxMock.EXPECT().Reply(gomock.Any(), gomock.Any()).MinTimes(1).Return(nil)
	b.EXPECT().Handle("/start", gomock.Any()).MinTimes(1)
	HandleOnboard(b, client, &currentStep)
	err := HandleStart(ctxMock, client, &currentStep)
	require.NoError(t, err)

	require.Equal(t, "addPassword", currentStep)
}

func TestHandleStorePrivateKey(t *testing.T) {
	testMnemonic := memguard.NewBufferFromBytes([]byte("test mnemonic"))
	testPassword := memguard.NewBufferFromBytes([]byte("test password"))
	encryptedMnemonic := "encryptedMnemonic"
	salt := "astrongsalt"
	step := "sendMnemonic"
	mockUser := &telebot.User{
		ID:       1,
		Username: "sender",
	}
	mockMessage := &telebot.Message{
		Text:   "test password",
		Sender: mockUser,
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
	ctxMock.EXPECT().Send(gomock.Any(), gomock.Any()).MinTimes(1).Return(nil)

	HandleStorePrivateKey(b, ctxMock, client, utilsMock, testMnemonic, testPassword, &step)
}

func TestHandleConfirmMnemonicStep(t *testing.T) {
	// Create mock bot, client, and context
	b := mock_internal.NewMockBot(gomock.NewController(t))
	c := mock_internal.NewMockTeleContext(gomock.NewController(t))
	botClient := mock_internal.NewMockBotClient(gomock.NewController(t))
	step := "confirmMnemonic"

	// Set test mnemonic with 24 words
	testMnemonic := "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12 word13 word14 word15 word16 word17 word18 word19 word20 word21 word22 word23 word24"

	c.EXPECT().Send(gomock.Any(), gomock.Any()).MinTimes(1).Return(nil)
	HandleConfirmMnemonicStep(b, c, botClient, testMnemonic, &step)

	// Verify that step is updated
	require.Equal(t, "receiveMnemonicWords", step)

}

func TestHandleReceiveMnemonicWords(t *testing.T) {
	// Create mock bot, client, and context
	b := mock_internal.NewMockBot(gomock.NewController(t))
	c := mock_internal.NewMockTeleContext(gomock.NewController(t))
	client := mock_internal.NewMockBotClient(gomock.NewController(t))
	mockUtils := mock_utils.NewMockUtilsInterface(gomock.NewController(t))
	step := "receiveMnemonicWords"

	// Set test mnemonic and random indexes
	testMnemonic := "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12 word13 word14 word15 word16 word17 word18 word19 word20 word21 word22 word23 word24"
	c.EXPECT().Text().Return("word3 word6 word9")
	c.EXPECT().Send(gomock.Any(), gomock.Any()).MinTimes(1).Return(nil)
	mockUtils.EXPECT().SplitMnemonic(gomock.Any()).Return([]string{"word3", "word6", "word9"})
	mockUtils.EXPECT().MnemonicChallenge(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	randomIndexes = [3]int{3, 6, 9}

	HandleReceiveMnemonicWords(b, c, client, mockUtils, testMnemonic, randomIndexes, &step)

	// Verify that mnemonic is confirmed and not back to the previous step
	require.NotEqual(t, "confirmMnemonic", step)
}

package zapmixin_test

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/xwjdsh/zapmixin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var conversations = []string{"06602963-d86d-3df3-ac06-000000000000"}

type Suite struct {
	suite.Suite
	client        *mixin.Client
	conversations []string
}

type MockedMixinClient struct {
	mock.Mock
	msgs []*mixin.MessageRequest
}

func (m *MockedMixinClient) SendMessage(ctx context.Context, message *mixin.MessageRequest) error {
	args := m.Called(ctx, message)
	m.msgs = append(m.msgs, message)
	return args.Error(0)
}

func TestLevels(t *testing.T) {
	client := new(MockedMixinClient)
	client.On("SendMessage", mock.Anything, mock.Anything).Return(nil)

	{
		h, _ := zapmixin.New(client, conversations, zapmixin.WithSync())
		logger := getLogger().WithOptions(zap.Hooks(h.Hook()))

		// default level: warn
		logger.Info("1")
		logger.Warn("2")
		logger.Error("3")

		assert.Equal(t, 2, len(client.msgs))
	}

	{
		client.msgs = []*mixin.MessageRequest{}
		h, _ := zapmixin.New(client, conversations, zapmixin.WithSync(), zapmixin.WithThresholdLevel(zapcore.ErrorLevel))
		logger := getLogger().WithOptions(zap.Hooks(h.Hook()))

		logger.Info("1")
		logger.Warn("2")
		logger.Error("3")

		assert.Equal(t, 1, len(client.msgs))
	}
}

func TestSuite(t *testing.T) {
	keyPath := os.Getenv("ZAPMIXIN_KEY_PATH")
	if keyPath == "" {
		t.Skip()
	}

	data, err := ioutil.ReadFile(keyPath)
	require.NoError(t, err)

	var store mixin.Keystore
	require.NoError(t, json.Unmarshal(data, &store))

	client, err := mixin.NewFromKeystore(&store)
	require.NoError(t, err)

	suite.Run(t, &Suite{
		client:        client,
		conversations: strings.Split(os.Getenv("ZAPMIXIN_CONVERSATIONS"), ","),
	})
}

func (s *Suite) TestBasic() {
	h, err := zapmixin.New(s.client, s.conversations, zapmixin.WithSync())
	s.NoError(err)

	logger := getLogger().WithOptions(zap.Hooks(h.Hook()))

	logger.Warn("test event")
}

func getLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(io.Discard),
		zapcore.InfoLevel,
	)

	return zap.New(core)
}

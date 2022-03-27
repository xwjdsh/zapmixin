package zapmixin

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLoggerName = "zapmixin"
)

type MixinClient interface {
	SendMessage(ctx context.Context, message *mixin.MessageRequest) error
}

func (h *Handler) sendMessage(e zapcore.Entry) error {
	for _, id := range h.conversations {
		req := &mixin.MessageRequest{
			ConversationID: id,
			MessageID:      uuid.Must(uuid.NewV4()).String(),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString([]byte(h.formatMessage(e))),
		}

		err := h.client.SendMessage(context.TODO(), req)
		if h.afterFunc != nil {
			return h.afterFunc(e, req, err)
		}

		return err
	}

	return nil
}

func (h *Handler) formatMessage(e zapcore.Entry) string {
	if h.formatter != nil {
		return h.formatter(e)
	}

	loggerName := defaultLoggerName
	if e.LoggerName != "" {
		loggerName = e.LoggerName
	}

	return fmt.Sprintf("Logger: %s\n%s\n%s\n%s", loggerName, e.Time, e.Level, e.Message)
}

package zapmixin

import (
	"github.com/fox-one/mixin-sdk-go"
	"go.uber.org/zap/zapcore"
)

type Handler struct {
	conversations []string
	client        MixinClient
	levels        []zapcore.Level
	async         bool
	formatter     func(zapcore.Entry) string
	filter        func(zapcore.Entry) bool
	afterFunc     func(zapcore.Entry, *mixin.MessageRequest, error) error
}

func New(client MixinClient, conversations []string, opts ...Option) (*Handler, error) {
	h := &Handler{
		client:        client,
		conversations: conversations,
		levels:        getLevelThreshold(zapcore.WarnLevel),
	}

	if err := h.Apply(opts...); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *Handler) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) Client() *mixin.Client {
	return h.client.(*mixin.Client)
}

func (h *Handler) Hook() func(zapcore.Entry) error {
	return func(e zapcore.Entry) error {
		if !h.levelMatched(e.Level) {
			return nil
		}

		if h.filter != nil && !h.filter(e) {
			return nil
		}

		if h.async {
			go h.sendMessage(e)
			return nil
		}

		return h.sendMessage(e)
	}
}

func (h *Handler) levelMatched(l zapcore.Level) bool {
	for _, level := range h.levels {
		if level == l {
			return true
		}
	}
	return false
}

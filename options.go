package zapmixin

import (
	"github.com/fox-one/mixin-sdk-go"
	"go.uber.org/zap/zapcore"
)

type Option func(*Handler) error

func WithThresholdLevel(l zapcore.Level) Option {
	return func(h *Handler) error {
		h.levels = getLevelThreshold(l)
		return nil
	}
}

func WithFixedLevel(l zapcore.Level) Option {
	return func(h *Handler) error {
		h.levels = []zapcore.Level{l}
		return nil
	}
}

func WithSync() Option {
	return func(h *Handler) error {
		h.async = false
		return nil
	}
}

func WithFormatter(f func(e zapcore.Entry) string) Option {
	return func(h *Handler) error {
		h.formatter = f
		return nil
	}
}

func WithFilter(f func(e zapcore.Entry) bool) Option {
	return func(h *Handler) error {
		h.filter = f
		return nil
	}
}

func WithAfter(f func(zapcore.Entry, *mixin.MessageRequest, error) error) Option {
	return func(h *Handler) error {
		h.afterFunc = f
		return nil
	}
}

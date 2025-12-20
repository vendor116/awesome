package httpserver

import (
	"errors"
	"time"
)

type options struct {
	readHeaderTimeout time.Duration
	readTimeout       time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	shutdownTimeout   time.Duration
	maxHeaderBytes    int
}

func getDefaultOptions() options {
	return options{
		readHeaderTimeout: readHeaderTimeout,
		readTimeout:       readTimeout,
		writeTimeout:      writeTimeout,
		idleTimeout:       idleTimeout,
		shutdownTimeout:   shutdownTimeout,
		maxHeaderBytes:    maxHeaderBytes,
	}
}

type OptionFunc func(o *options) error

func WithReadHeaderTimeout(timeout time.Duration) OptionFunc {
	return func(o *options) error {
		o.readHeaderTimeout = timeout
		return nil
	}
}

func WithWriteTimeout(timeout time.Duration) OptionFunc {
	return func(o *options) error {
		o.writeTimeout = timeout
		return nil
	}
}

func WithShutdownTimeout(timeout time.Duration) OptionFunc {
	return func(o *options) error {
		o.shutdownTimeout = timeout
		return nil
	}
}

func WithIdleTimeout(timeout time.Duration) OptionFunc {
	return func(o *options) error {
		o.idleTimeout = timeout
		return nil
	}
}

func WithMaxHeaderBytes(size int) OptionFunc {
	return func(o *options) error {
		if size <= 0 {
			return errors.New("size must be greater than zero")
		}

		o.maxHeaderBytes = size
		return nil
	}
}

func WithReadTimeout(timeout time.Duration) OptionFunc {
	return func(o *options) error {
		o.readTimeout = timeout
		return nil
	}
}

package geoip2

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWithErrorFunc(t *testing.T) {
	assert := assert.New(t)

	ch := make(chan int, 1)
	tests := map[string]struct {
		fun    func(err error)
		ch     chan int
		output int
	}{
		"success": {fun: func(err error) { ch <- 1 }, ch: ch, output: 1},
	}

	for _, t := range tests {
		cfg := &downloadConfig{}
		opt := WithErrorFunc(t.fun)
		opt(cfg)
		cfg.errorFunc(nil)
		assert.Equal(t.output, <-t.ch)
	}
}

func TestWithUpdateInterval(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct{
		interval time.Duration
		output time.Duration
	}{
		"success": {interval: time.Hour, output: time.Hour},
	}

	for _, t := range tests {
		cfg := &downloadConfig{}
		opt := WithUpdateInterval(t.interval)
		opt(cfg)
		assert.Equal(t.output, cfg.updateInterval)
	}
}

func TestWithSuccessFunc(t *testing.T) {
	assert := assert.New(t)

	ch := make(chan int, 1)
	tests := map[string]struct{
		fun    func()
		ch     chan int
		output int
	}{
		"success": {fun: func() { ch <- 1 }, ch: ch, output: 1},
	}

	for _, t := range tests {
		cfg := &downloadConfig{}
		opt := WithSuccessFunc(t.fun)
		opt(cfg)
		cfg.successFunc()
		assert.Equal(t.output, <-t.ch)
	}
}

func TestWithRetries(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct{
		retries int
		output int
	}{
		"success": {retries: 10,  output: 10},
	}

	for _, t := range tests {
		cfg := &downloadConfig{}
		opt := WithRetries(t.retries)
		opt(cfg)
		assert.Equal(t.output, cfg.retries)
	}
}

func TestWithFirstDownloadWait(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct{
		first time.Duration
		output  time.Duration
	}{
		"success": {first: time.Hour, output: time.Hour},
	}

	for _, t := range tests {
		cfg := &downloadConfig{}
		opt := WithFirstDownloadWait(t.first)
		opt(cfg)
		assert.Equal(t.output, cfg.firstDownloadWait)
	}
}

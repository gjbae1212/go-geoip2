package geoip2

import "time"

type DownloadOption interface {
	apply(dc *downloadConfig)
}

type DownloadOptionFunc func(cfg *downloadConfig)

func (dof DownloadOptionFunc) apply(cfg *downloadConfig) { dof(cfg) }

type downloadConfig struct {
	downloadURL       string
	checksumURL       string
	targetPath        string
	firstDownloadWait time.Duration
	updateInterval    time.Duration
	retries           int
	successFunc       func()
	errorFunc         func(err error)
	checksum          string
}

// WithUpdateInterval returns a function for setting download time interval.
func WithUpdateInterval(d time.Duration) DownloadOptionFunc {
	return func(cfg *downloadConfig) { cfg.updateInterval = d }
}

// WithUpdateInterval returns a function for setting download retry count if a download is failed.
func WithRetries(retries int) DownloadOptionFunc {
	return func(cfg *downloadConfig) { cfg.retries = retries }
}

// WithSuccessFunc returns a function for setting a method to call if a download succeeded.
func WithSuccessFunc(f func()) DownloadOptionFunc {
	return func(cfg *downloadConfig) { cfg.successFunc = f }
}

// WithErrorFunc returns a function for setting a method to call if a download failed.
func WithErrorFunc(f func(error)) DownloadOptionFunc {
	return func(cfg *downloadConfig) { cfg.errorFunc = f }
}

// WithFirstDownloadWait returns a function for setting first download wait time.
func WithFirstDownloadWait(d time.Duration) DownloadOptionFunc {
	return func(cfg *downloadConfig) { cfg.firstDownloadWait = d }
}

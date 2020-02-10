package geoip2

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"net"
	"time"

	geoip2_golang "github.com/oschwald/geoip2-golang"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

type MaxmindDownloadSuffix string

const (
	MaxmindDownloadFormat = "https://download.maxmind.com/app/geoip_download?license_key=%s&edition_id=%s&suffix=%s"
	GZIP                  = MaxmindDownloadSuffix("tar.gz")
	MD5                   = MaxmindDownloadSuffix("tar.gz.md5")
)

var (
	ErrInvalidParameters = fmt.Errorf("[err] invalid parameters")
	ErrNotFoundDatabase  = fmt.Errorf("[err] not found database")
	ErrFirstDownloadFail = fmt.Errorf("[err] first download fail")
)

// support to interface for oschwald/geoip2-golang.
type Reader interface {
	ASN(ipAddress net.IP) (*geoip2_golang.ASN, error)
	AnonymousIP(ipAddress net.IP) (*geoip2_golang.AnonymousIP, error)
	City(ipAddress net.IP) (*geoip2_golang.City, error)
	ConnectionType(ipAddress net.IP) (*geoip2_golang.ConnectionType, error)
	Country(ipAddress net.IP) (*geoip2_golang.Country, error)
	Domain(ipAddress net.IP) (*geoip2_golang.Domain, error)
	Enterprise(ipAddress net.IP) (*geoip2_golang.Enterprise, error)
	ISP(ipAddress net.IP) (*geoip2_golang.ISP, error)
	Metadata() maxminddb.Metadata
	Close() error
}

// Open returns geoip Reader from a local file.
func Open(file string) (Reader, error) {
	db, err := geoip2_golang.Open(file)
	if err != nil {
		return nil, err
	}

	return &fileReader{db}, nil
}

// OpenURL returns geoip Reader from maxmind download URL and updates automatically the latest maxmind databases.
// reference: maxmind URL https://dev.maxmind.com/geoip/geoipupdate/#Direct_Downloads
func OpenURL(licenseKey, editionId, targetPath string, opts ...DownloadOption) (Reader, error) {
	if licenseKey == "" || editionId == "" || targetPath == "" {
		return nil, fmt.Errorf("[err] OpenURL %w", ErrInvalidParameters)
	}

	// generate maxmind download URL
	downloadURL, err := MaxmindDownloadURL(licenseKey, editionId, GZIP)
	if err != nil {
		return nil, fmt.Errorf("[err] OpenURL %w", err)
	}

	// generate maxmind checksum URL
	checkSumURL, err := MaxmindDownloadURL(licenseKey, editionId, MD5)
	if err != nil {
		return nil, fmt.Errorf("[err] OpenURL %w", err)
	}

	cfg := &downloadConfig{
		downloadURL:       downloadURL,
		checksumURL:       checkSumURL,
		targetPath:        targetPath,
		firstDownloadWait: 10 * time.Second,
		updateInterval:    time.Hour,
		retries:           1,
		successFunc:       func() {},
		errorFunc:         func(err error) {},
	}

	// dependency injection.
	for _, opt := range opts {
		opt.apply(cfg)
	}

	reader := &downloadReader{cfg: cfg, backoff: backoff.NewExponentialBackOff()}

	// if maxmind database is already exist, using it.
	reader.databaseReload(cfg.targetPath, "")

	// run update and download logic async
	go reader.runDownloadURL()

	// if default db exists, returning.
	if reader.db != nil {
		return reader, nil
	}

	// wait first download success
Wait:
	for {
		select {
		case <-time.After(time.Second):
			reader.RLock()
			if reader.db != nil {
				reader.RUnlock()
				break Wait
			}
			reader.RUnlock()
		case <-time.After(reader.cfg.firstDownloadWait):
			break Wait
		}
	}
	if reader.db == nil {
		return nil, ErrFirstDownloadFail
	}

	return reader, nil
}

// maxmindDownloadURL returns Maxmind download URL
// reference: maxmind URL https://dev.maxmind.com/geoip/geoipupdate/#Direct_Downloads
func MaxmindDownloadURL(licenseKey, editionId string, suffix MaxmindDownloadSuffix) (string, error) {
	if licenseKey == "" || editionId == "" || suffix == "" {
		return "", fmt.Errorf("[err] maxmindDownloadURL %w", ErrInvalidParameters)
	}
	return fmt.Sprintf(MaxmindDownloadFormat, licenseKey, editionId, suffix), nil
}

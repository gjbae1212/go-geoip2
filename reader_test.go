package geoip2

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	geoip2_golang "github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"
)

var (
	testReader Reader
	testOnce   sync.Once
)

func TestDownloadReader_AnonymousIP(t *testing.T) {
	assert := assert.New(t)

	reader := testDownloadReader()
	if reader != nil {
		tests := map[string]struct {
			ip net.IP
		}{
			"success": {ip: net.ParseIP("8.8.8.8")},
		}

		for _, t := range tests {
			_, err := reader.AnonymousIP(t.ip)
			if err != nil {
				assert.True(errors.As(err, &geoip2_golang.InvalidMethodError{}))
			}
		}
	}
}

func TestDownloadReader_ASN(t *testing.T) {
	assert := assert.New(t)

	reader := testDownloadReader()
	if reader != nil {
		tests := map[string]struct {
			ip net.IP
		}{
			"success": {ip: net.ParseIP("8.8.8.8")},
		}

		for _, t := range tests {
			if _, err := reader.ASN(t.ip); err != nil {
				assert.True(errors.As(err, &geoip2_golang.InvalidMethodError{}))
			}
		}
	}
}

func TestDownloadReader_City(t *testing.T) {
	assert := assert.New(t)

	reader := testDownloadReader()
	if reader != nil {
		tests := map[string]struct {
			ip net.IP
		}{
			"success": {ip: net.ParseIP("8.8.8.8")},
		}

		for _, t := range tests {
			if _, err := reader.City(t.ip); err != nil {
				assert.True(errors.As(err, &geoip2_golang.InvalidMethodError{}))
			}
		}
	}
}

func TestDownloadReader_ConnectionType(t *testing.T) {
	assert := assert.New(t)

	reader := testDownloadReader()
	if reader != nil {
		tests := map[string]struct {
			ip net.IP
		}{
			"success": {ip: net.ParseIP("8.8.8.8")},
		}

		for _, t := range tests {
			if _, err := reader.ConnectionType(t.ip); err != nil {
				assert.True(errors.As(err, &geoip2_golang.InvalidMethodError{}))
			}
		}
	}
}

func TestDownloadReader_Country(t *testing.T) {
	assert := assert.New(t)

	reader := testDownloadReader()
	if reader != nil {
		tests := map[string]struct {
			ip net.IP
		}{
			"success": {ip: net.ParseIP("8.8.8.8")},
		}

		for _, t := range tests {
			if c, err := reader.Country(t.ip); err != nil {
				assert.True(errors.As(err, &geoip2_golang.InvalidMethodError{}))
			} else {
				log.Println(c)
			}
		}
	}
}

func TestDownloadReader_Domain(t *testing.T) {
	assert := assert.New(t)

	reader := testDownloadReader()
	if reader != nil {
		tests := map[string]struct {
			ip net.IP
		}{
			"success": {ip: net.ParseIP("8.8.8.8")},
		}

		for _, t := range tests {
			if _, err := reader.Domain(t.ip); err != nil {
				assert.True(errors.As(err, &geoip2_golang.InvalidMethodError{}))
			}
		}
	}
}

func TestDownloadReader_Enterprise(t *testing.T) {
	assert := assert.New(t)

	reader := testDownloadReader()
	if reader != nil {
		tests := map[string]struct {
			ip net.IP
		}{
			"success": {ip: net.ParseIP("8.8.8.8")},
		}

		for _, t := range tests {
			if _, err := reader.Enterprise(t.ip); err != nil {
				assert.True(errors.As(err, &geoip2_golang.InvalidMethodError{}))
			}
		}
	}
}

func TestDownloadReader_Metadata(t *testing.T) {
	assert := assert.New(t)

	reader := testDownloadReader()
	if reader != nil {
		tests := map[string]struct {
		}{
			"success": {},
		}

		for _, _ = range tests {
			meta := reader.Metadata()
			log.Println(meta)
			_ = assert
		}
	}
}

func TestDownloadReader_ISP(t *testing.T) {
	assert := assert.New(t)

	reader := testDownloadReader()
	if reader != nil {
		tests := map[string]struct {
			ip net.IP
		}{
			"success": {ip: net.ParseIP("8.8.8.8")},
		}

		for _, t := range tests {
			if _, err := reader.ISP(t.ip); err != nil {
				assert.True(errors.As(err, &geoip2_golang.InvalidMethodError{}))
			}
		}
	}
}

func TestDownloadReader_Close(t *testing.T) {
	assert := assert.New(t)

	licenseKey := os.Getenv("MAXMIND_LICENSE_KEY")
	editionId := os.Getenv("MAXMIND_EDITION_ID")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dir := filepath.Dir(filename)
	if licenseKey != "" && editionId != "" {
		targetPath := filepath.Join(dir, "test.mmdb")
		os.RemoveAll(targetPath)
		r, err := OpenURL(licenseKey, editionId, targetPath)
		assert.NoError(err)
		assert.NoError(r.Close())
	}
}

func TestCheckUpdateDownload(t *testing.T) {
	assert := assert.New(t)
	licenseKey := os.Getenv("MAXMIND_LICENSE_KEY")
	editionId := os.Getenv("MAXMIND_EDITION_ID")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dir := filepath.Dir(filename)

	if licenseKey != "" && editionId != "" {
		targetPath := filepath.Join(dir, "test2.mmdb")
		reader, err := OpenURL(licenseKey, editionId, targetPath,
			WithUpdateInterval(1*time.Second), WithSuccessFunc(func() {
				fmt.Println("")
				fmt.Println("update success")
			}), WithErrorFunc(func(err error) {
				fmt.Println(err)
			}))
		if err != nil {
			panic(err)
		}

		_ = reader
		_ = assert
		time.Sleep(time.Second * 5)
	}
}

func testDownloadReader() (reader Reader) {
	testOnce.Do(func() {
		licenseKey := os.Getenv("MAXMIND_LICENSE_KEY")
		editionId := os.Getenv("MAXMIND_EDITION_ID")
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("No caller information")
		}
		dir := filepath.Dir(filename)
		if licenseKey != "" && editionId != "" {
			targetPath := filepath.Join(dir, "test.mmdb")
			os.RemoveAll(targetPath)
			r, err := OpenURL(licenseKey, editionId, targetPath)
			if err != nil {
				panic(err)
			}
			testReader = r
		}
	})
	return testReader
}

func BenchmarkDefaultReader_Country(b *testing.B) {
	dbpath := os.Getenv("MAXMIND_DB_PATH")
	if dbpath != "" {
		reader, err := Open(dbpath)
		if err != nil {
			panic(err)
		}
		for i := 0; i < b.N; i++ {
			_, err := reader.Country(net.ParseIP("1.1.1.1"))
			if err != nil {
				panic(err)
			}
		}
	}
}

func BenchmarkDownloadReader_Country(b *testing.B) {
	licenseKey := os.Getenv("MAXMIND_LICENSE_KEY")
	editionId := os.Getenv("MAXMIND_EDITION_ID")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dir := filepath.Dir(filename)

	if licenseKey != "" && editionId != "" {
		targetPath := filepath.Join(dir, "test.mmdb")
		os.RemoveAll(targetPath)
		reader, err := OpenURL(licenseKey, editionId, targetPath,
			WithUpdateInterval(3*time.Second), WithSuccessFunc(func() {
				fmt.Println("")
				fmt.Println("update success")
			}), WithErrorFunc(func(err error) {
				fmt.Println(err)
			}))
		if err != nil {
			panic(err)
		}
		for i := 0; i < b.N; i++ {
			_, err := reader.Country(net.ParseIP("1.1.1.1"))
			if err != nil {
				panic(err)
			}
		}
	}
}

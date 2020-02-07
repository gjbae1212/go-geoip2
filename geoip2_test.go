package geoip2

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMaxmindDownloadURL(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		licenseKey string
		editionId  string
		suffix     MaxmindDownloadSuffix
		output     string
	}{
		"fail": {},
		"success": {licenseKey: "license-key", editionId: "edition", suffix: MD5,
			output: "https://download.maxmind.com/app/geoip_download?license_key=license-key&edition_id=edition&suffix=tar.gz.md5"},
	}

	for _, t := range tests {
		url, _ := MaxmindDownloadURL(t.licenseKey, t.editionId, t.suffix)
		assert.Equal(t.output, url)
	}
}

func TestOpen(t *testing.T) {
	assert := assert.New(t)

	dbpath := os.Getenv("MAXMIND_DB_PATH")
	if dbpath != "" {
		tests := map[string]struct {
			input string
			isErr bool
		}{
			"success": {input: dbpath, isErr: false},
		}

		for _, t := range tests {
			_, err := Open(t.input)
			assert.Equal(t.isErr, err != nil)
		}
	}
}

func TestOpenURL(t *testing.T) {
	assert := assert.New(t)

	licenseKey := os.Getenv("MAXMIND_LICENSE_KEY")
	editionId := os.Getenv("MAXMIND_EDITION_ID")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	case1 := filepath.Join(path.Dir(filename), "case1.mmdb")
	case2 := filepath.Join(path.Dir(filename), "case2.mmdb")
	os.RemoveAll(case1)

	if licenseKey != "" && editionId != "" {
		tests := map[string]struct {
			licenseKey string
			editionId  string
			targetPath string
			isErr      bool
		}{
			"case 1": {licenseKey: licenseKey, editionId: editionId, targetPath: case1},
			"case 2": {licenseKey: licenseKey, editionId: editionId, targetPath: case2},
		}

		for k, tc := range tests {
			t.Run(k, func(t *testing.T) {
				_, err := OpenURL(tc.licenseKey, tc.editionId, tc.targetPath, WithSuccessFunc(func() {
					assert.True(true)
				}), WithErrorFunc(func(err error) {
					panic(fmt.Errorf("[TestOpenURL] fail %w", err))
				}))
				assert.Equal(tc.isErr, err != nil)
			})

		}
	}
}

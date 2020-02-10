package geoip2

import (
	"archive/tar"
	"compress/gzip"
	"fmt"

	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	geoip2_golang "github.com/oschwald/geoip2-golang"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

type fileReader struct {
	*geoip2_golang.Reader
}

type downloadReader struct {
	sync.RWMutex
	db               *geoip2_golang.Reader
	cfg              *downloadConfig
	runDownloadClose chan bool
	backoff          *backoff.ExponentialBackOff
}

// ASN is the same method as that "github.com/oschwald/geoip2-golang" is.
func (r *downloadReader) ASN(ipAddress net.IP) (*geoip2_golang.ASN, error) {
	r.RLock()
	defer r.RUnlock()

	return r.db.ASN(ipAddress)
}

// AnonymousIP is the same method as that "github.com/oschwald/geoip2-golang" is.
func (r *downloadReader) AnonymousIP(ipAddress net.IP) (*geoip2_golang.AnonymousIP, error) {
	r.RLock()
	defer r.RUnlock()

	return r.db.AnonymousIP(ipAddress)
}

// City is the same method as that "github.com/oschwald/geoip2-golang" is.
func (r *downloadReader) City(ipAddress net.IP) (*geoip2_golang.City, error) {
	r.RLock()
	defer r.RUnlock()

	return r.db.City(ipAddress)
}

// ConnectionType is the same method as that "github.com/oschwald/geoip2-golang" is.
func (r *downloadReader) ConnectionType(ipAddress net.IP) (*geoip2_golang.ConnectionType, error) {
	r.RLock()
	defer r.RUnlock()

	return r.db.ConnectionType(ipAddress)
}

// Country is the same method as that "github.com/oschwald/geoip2-golang" is.
func (r *downloadReader) Country(ipAddress net.IP) (*geoip2_golang.Country, error) {
	r.RLock()
	defer r.RUnlock()

	return r.db.Country(ipAddress)
}

// Domain is the same method as that "github.com/oschwald/geoip2-golang" is.
func (r *downloadReader) Domain(ipAddress net.IP) (*geoip2_golang.Domain, error) {
	r.RLock()
	defer r.RUnlock()

	return r.db.Domain(ipAddress)
}

// Enterprise is the same method as that "github.com/oschwald/geoip2-golang" is.
func (r *downloadReader) Enterprise(ipAddress net.IP) (*geoip2_golang.Enterprise, error) {
	r.RLock()
	defer r.RUnlock()

	return r.db.Enterprise(ipAddress)
}

// ISP is the same method as that "github.com/oschwald/geoip2-golang" is.
func (r *downloadReader) ISP(ipAddress net.IP) (*geoip2_golang.ISP, error) {
	r.RLock()
	defer r.RUnlock()

	return r.db.ISP(ipAddress)
}

// Metadata is the same method as that "github.com/oschwald/geoip2-golang" is.
func (r *downloadReader) Metadata() maxminddb.Metadata {
	r.RLock()
	defer r.RUnlock()

	return r.db.Metadata()
}

// Close is the same method as that "github.com/oschwald/geoip2-golang" is.
func (r *downloadReader) Close() error {
	r.Lock()
	defer r.Unlock()
	close(r.runDownloadClose)

	return r.db.Close()
}

func (r *downloadReader) runDownloadURL() {
	for {
		// getting checksum
		var remoteChecksum string
		for i := 0; i < r.cfg.retries; i++ {
			// wait for backoff interval.
			time.Sleep(r.backoff.NextBackOff())

			c, err := r.downloadChecksum()
			if err != nil {
				r.cfg.errorFunc(fmt.Errorf("[err] runDownloadURL %w", err))
				continue
			}
			remoteChecksum = strings.TrimSpace(c)
		}
		// reset backoff.
		r.backoff.Reset()

		if remoteChecksum == "" {
			r.cfg.errorFunc(fmt.Errorf("[err] runDownloadURL checksum download fail"))
		} else {
			// if local checksum is equal to remote checksum, updating maxmind database.
			if remoteChecksum != r.cfg.checksum {
				for i := 0; i < r.cfg.retries; i++ {
					// wait for backoff interval.
					time.Sleep(r.backoff.NextBackOff())

					// downloading database.
					tempPath, err := r.downloadDatabase()
					if err != nil {
						r.cfg.errorFunc(fmt.Errorf("[err] runDownloadURL %w", err))
						continue
					}

					// reload new database.
					if err := r.databaseReload(tempPath, remoteChecksum); err != nil {
						r.cfg.errorFunc(fmt.Errorf("[err] runDownloadURL %w", err))
						continue
					}

					// call a success function.
					r.cfg.successFunc()
				}
				// reset backoff.
				r.backoff.Reset()
			}
		}

		select {
		case <-r.runDownloadClose:
			return
		case <-time.After(r.cfg.updateInterval):
		}
	}
}

// databaseReload reloads maxmind database.
func (r *downloadReader) databaseReload(tempPath, checksum string) error {
	if tempPath == "" {
		return fmt.Errorf("[err] databaseReload %w", ErrInvalidParameters)
	}
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		return fmt.Errorf("[err] databaseReload %w", ErrNotFoundDatabase)
	}
	r.Lock()
	defer r.Unlock()

	// make directory.
	if _, err := os.Stat(r.cfg.storeDir); os.IsNotExist(err) {
		if err := os.MkdirAll(r.cfg.storeDir, os.ModePerm); err != nil {
			return fmt.Errorf("[err] databaseReload %w", err)
		}
	}

	// backup old database.
	dbpath := r.cfg.dbPath()
	dbBackupPath := r.cfg.dbBackupPath()
	checksumPath := r.cfg.checksumPath()
	if info, err := os.Stat(dbpath); info != nil || os.IsExist(err) {
		if dbpath != tempPath {
			// make backup old database
			os.Rename(dbpath, dbBackupPath)
		}
	}

	// move tempPath to dbpath.
	if err := os.Rename(tempPath, dbpath); err != nil {
		// delete new database
		os.RemoveAll(tempPath)
		// rollback old database
		os.Rename(dbBackupPath, dbpath)
		return fmt.Errorf("[err] databaseReload %w", err)
	}

	// open new database.
	db, err := geoip2_golang.Open(dbpath)
	if err != nil {
		// delete new database
		os.RemoveAll(dbpath)
		// rollback old database
		os.Rename(dbBackupPath, dbpath)
		return fmt.Errorf("[err] databaseReload %w", err)
	}

	// delete back old database
	if info, err := os.Stat(dbBackupPath); info != nil || os.IsExist(err) {
		os.RemoveAll(dbBackupPath)
	}

	// release old database
	if r.db != nil {
		if err := r.db.Close(); err != nil {
			fmt.Printf("[err] databaseReload old database close %v", err)
		}
		r.db = nil
		r.cfg.checksum = ""
	}

	if checksum == "" {
		// read md5 file.
		if bys, err := ioutil.ReadFile(checksumPath); err == nil {
			checksum = string(bys)
		}
	} else {
		// write md5 to file.
		if cpath, err := os.Create(r.cfg.checksumPath()); err == nil {
			cpath.WriteString(checksum)
			cpath.Close()
		}
	}

	r.db = db
	r.cfg.checksum = checksum
	return nil
}

// requestChecksum requests checksum data.
func (r *downloadReader) downloadChecksum() (checksum string, err error) {
	resp, suberr := http.Get(r.cfg.checksumURL)
	if resp != nil {
		defer resp.Body.Close()
	} else {
		err = fmt.Errorf("[err] downloadChecksum resp nil")
		return
	}
	if suberr != nil {
		err = fmt.Errorf("[err] downloadChecksum %w", suberr)
		return
	}

	status := resp.StatusCode
	if resp.StatusCode/100 != 2 {
		err = fmt.Errorf("[err] downloadChecksum status %d", status)
		return

	}

	data, suberr := ioutil.ReadAll(resp.Body)
	if suberr != nil {
		err = fmt.Errorf("[err] downloadChecksum status %w", suberr)
		return
	}

	checksum = string(data)
	return
}

// request requests checksum data.
func (r *downloadReader) downloadDatabase() (tempPath string, err error) {
	// download database
	resp, suberr := http.Get(r.cfg.downloadURL)
	if resp != nil {
		defer resp.Body.Close()
	} else {
		err = fmt.Errorf("[err] downloadDatabase resp nil")
		return
	}
	if suberr != nil {
		err = fmt.Errorf("[err] downloadDatabase %w", suberr)
		return
	}

	status := resp.StatusCode
	if resp.StatusCode/100 != 2 {
		err = fmt.Errorf("[err] downloadDatabase status %d", status)
		return
	}

	// save database to temporary path
	fpath := filepath.Join(os.TempDir(), fmt.Sprintf("maxmind-%d.mmdb", time.Now().UnixNano()))
	f, suberr := os.Create(fpath)
	if suberr != nil {
		err = fmt.Errorf("[err] downloadDatabase %w", suberr)
		return
	}
	defer f.Close()

	// wrapping unzip reader
	gr, suberr := gzip.NewReader(resp.Body)
	if suberr != nil {
		err = fmt.Errorf("[err] downloadDatabase %w", suberr)
		return
	}
	defer gr.Close()

	// wrapping tar reader
	tr := tar.NewReader(gr)

Search:
	for {
		header, suberr := tr.Next()
		// check error
		switch {
		case suberr == io.EOF:
			err = fmt.Errorf("[err] downloadDatabase not found mmdb in gzip.")
			return
		case suberr != nil:
			err = fmt.Errorf("[err] downloadDatabase read gzip %w", suberr)
			return
		case header == nil:
			continue
		}

		// search mmdb
		switch header.Typeflag {
		case tar.TypeReg:
			if strings.HasSuffix(strings.ToLower(header.Name), "mmdb") {
				if _, suberr := io.Copy(f, gr); suberr != nil {
					err = fmt.Errorf("[err] downloadDatabase read gzip %w", suberr)
					return
				}
				break Search
			}
		}
	}

	tempPath = fpath
	return
}

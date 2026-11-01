package certpool

import (
	"crypto/x509"
	_ "embed"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"
)

//go:embed cacert.pem
var embeddedMozillaRoots []byte

var initOnce sync.Once

// EnsureSystemCertPool ensures that a valid root certificate pool is available on all platforms,
// with special handling for iOS / Darwin where /etc/ssl/cert.pem may not exist.
func EnsureSystemCertPool() {
	initOnce.Do(func() {
		// 1. If SSL_CERT_FILE is already specified and exists, honor it
		if customFile := os.Getenv("SSL_CERT_FILE"); customFile != "" {
			if _, err := os.Stat(customFile); err == nil {
				return
			}
		}

		// 2. Search candidate local paths
		candidatePaths := []string{
			"./config/cacert.pem",
			"./cacert.pem",
			"/var/mobile/Documents/service/cliproxyapi/config/cacert.pem",
			"/var/mobile/Documents/service/cliproxyapi/cacert.pem",
		}

		if exe, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exe)
			candidatePaths = append(candidatePaths,
				filepath.Join(exeDir, "config", "cacert.pem"),
				filepath.Join(exeDir, "cacert.pem"),
			)
		}

		for _, p := range candidatePaths {
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() && fi.Size() > 0 {
				_ = os.Setenv("SSL_CERT_FILE", p)
				log.Debugf("certpool: automatically loaded SSL_CERT_FILE from %s", p)
				return
			}
		}

		// 3. If running on iOS/Darwin or system pool fails, write embedded roots to a fallback file and set SSL_CERT_FILE
		if runtime.GOOS == "darwin" || runtime.GOOS == "ios" {
			pool, err := x509.SystemCertPool()
			if err != nil || pool == nil {
				if len(embeddedMozillaRoots) > 0 {
					fallbackFile := filepath.Join(os.TempDir(), "cliproxyapi_cacert.pem")
					if errWrite := os.WriteFile(fallbackFile, embeddedMozillaRoots, 0o644); errWrite == nil {
						_ = os.Setenv("SSL_CERT_FILE", fallbackFile)
						log.Infof("certpool: system cert pool unavailable on %s; loaded embedded Mozilla roots to %s", runtime.GOOS, fallbackFile)
					}
				}
			}
		}
	})
}

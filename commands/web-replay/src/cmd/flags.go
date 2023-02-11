package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	fsExtras "github.com/alshdavid/web-replay/src/platform/fs-extras"
)

type flags struct {
	sslPrivateKeyFilepath  string
	sslCertificateFilepath string
	harFilePath            string
	patchesFolder          string
	domain                 string
	host                   string
	startingPort           int
	logLevel               string
	xhrBlockingTime        time.Duration
}

func parseFlags(args []string) *flags {
	var sslPrivateKeyFilepath string
	var sslCertificateFilepath string
	var harFilePath string
	var patchesFolder string
	var domain string
	var host string
	var logLevel string
	var startingPort int
	var xhrBlockingTime int64

	ex, _ := os.Executable()
	f := flag.NewFlagSet("serve", flag.ExitOnError)

	// Optional Args
	f.StringVar(&sslPrivateKeyFilepath, "ssl-key", filepath.Join(filepath.Dir(ex), "server.key"), "Path to the private key for the SSL server (server.key)")
	f.StringVar(&sslCertificateFilepath, "ssl-cert", filepath.Join(filepath.Dir(ex), "server.crt"), "Path to the certificate file key for the SSL server (server.crt)")
	f.StringVar(&host, "host", "127.0.0.1", "Host address to use for local server")
	f.StringVar(&domain, "domain", "localhost", "Host domain to use for local server")
	f.StringVar(&logLevel, "log-level", "normal", "Verbosity ('normal', 'silent')")
	f.IntVar(&startingPort, "port", 3000, "Starting port")
	f.Int64Var(&xhrBlockingTime, "xhr-blocking-time", int64(time.Millisecond*500), "How long mock xhr request should take to complete")

	f.Parse(args)

	harFilePath = f.Arg(0)

	if harFilePath == "" {
		fmt.Println("Please use \"--har\" flag to set path to har file")
		os.Exit(1)
	}

	sslPrivateKeyFilepath, _ = filepath.Abs(sslPrivateKeyFilepath)
	sslCertificateFilepath, _ = filepath.Abs(sslCertificateFilepath)
	harFilePath, _ = filepath.Abs(harFilePath)
	patchesFolder = filepath.Join(filepath.Dir(ex), "..", "patches", "enabled")

	return &flags{
		sslPrivateKeyFilepath:  sslPrivateKeyFilepath,
		sslCertificateFilepath: sslCertificateFilepath,
		harFilePath:            harFilePath,
		patchesFolder:          patchesFolder,
		startingPort:           startingPort,
		host:                   host,
		logLevel:               logLevel,
		domain:                 domain,
		xhrBlockingTime:        time.Duration(xhrBlockingTime),
	}
}

func validateFlags(f *flags) error {
	if !fsExtras.Exists(f.harFilePath) {
		return errors.New("no har file supplied")
	}

	if !fsExtras.Exists(f.sslPrivateKeyFilepath) {
		return errors.New("no SSL private key file supplied")
	}

	if !fsExtras.Exists(f.sslCertificateFilepath) {
		return errors.New("no SSL public file supplied")
	}

	return nil
}

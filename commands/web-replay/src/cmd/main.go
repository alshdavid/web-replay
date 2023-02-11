package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	internal_serve "github.com/alshdavid/web-replay/src/cmd/internal"
	"github.com/alshdavid/web-replay/src/platform/extras"
	"github.com/alshdavid/web-replay/src/platform/har"
)

var TriggerShutdownSignal = make(chan bool)
var OnServersClosedSignal = make(chan bool)
var OnServersUpSignal = make(chan bool)

func main() {
	env := parseFlags(os.Args[1:])

	if err := validateFlags(env); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger := internal_serve.NewLogger(env.logLevel)

	serverPort := env.startingPort
	blockXhrRequestTime := env.xhrBlockingTime
	harEntries := har.MustParseFile(env.harFilePath)

	// Parse har file and generate settings that correspond to each server
	serverSettings := internal_serve.GenerateSettings(serverPort, env.domain, harEntries)
	domainMap := internal_serve.GenerateDomainMap(serverSettings)
	servicePatches, loaded, err := internal_serve.LoadPatchesFromDir(env.patchesFolder)
	if err != nil {
		fmt.Println("Unable to parse patch YAML")
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Println("FILENAME:")
	logger.Printf(" %s\n", extras.Must(filepath.Rel(extras.Must(os.Getwd()), env.harFilePath)))
	logger.Println("")

	logger.Println("SSL:")
	logger.Printf(" Cert File: %s\n", env.sslCertificateFilepath)
	logger.Printf(" Key File:  %s\n", env.sslPrivateKeyFilepath)
	logger.Println("")

	logger.Println("SERVICES:")
	for _, setting := range serverSettings {
		logger.Printf(" https://%s:%d \t %s\n", env.domain, setting.Port, setting.OriginalHost)
	}

	logger.Println("")
	logger.Println("PATCHES APPLIED:")
	if len(loaded) == 0 {
		logger.Println(" <none>")
	}
	for _, patch := range loaded {
		logger.Printf(" %s\n", patch.Name())
	}

	servers := []*http.Server{}

	for _, setting := range serverSettings {
		address := fmt.Sprintf("%s:%d", env.host, setting.Port)
		srv := &http.Server{
			Addr:     address,
			ErrorLog: log.New(io.Discard, "", 0),
			Handler: internal_serve.Handler(
				setting,
				domainMap,
				blockXhrRequestTime,
				servicePatches,
				logger,
			),
		}
		servers = append(servers, srv)
	}

	go func() {
		<-TriggerShutdownSignal
		for _, srv := range servers {
			srv.Close()
		}
		OnServersClosedSignal <- true
	}()

	for _, srv := range servers {
		go srv.ListenAndServeTLS(
			env.sslCertificateFilepath,
			env.sslPrivateKeyFilepath,
		)
	}

	logger.Printf("\nLOGS:\n")

	OnServersUpSignal <- true
	<-OnServersClosedSignal
	os.Exit(0)
}

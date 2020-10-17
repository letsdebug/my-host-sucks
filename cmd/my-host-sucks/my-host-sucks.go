package main

// Rough program outline:
//
// 1. When NOT running on the cPanel server (e.g. running on the user's workstation):
//   a. Ask for the details of the cPanel service (URL, credentials).
//   b. Enumerate all of the domain names on the cPanel server.
//   c. Create an order for all eligible domains and attempt to fulfill as many authorizations
//      as possible.
//   d. Create the appropriate certificates (if needed) and install them to each virtual host.
//   e. Upload the program via API to the server.
//   f. Create a crontask on the server if it does not exist already.
//
// 2. If the program IS running on the cPanel server:
//   a. Run [1b..1d] and 1f using the local cPanel UAPI.

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/letsdebug/my-host-sucks/cpanel"
)

var (
	cpanelClient *cpanel.Client
	cpanelCtx    context.Context = context.Background()
)

func main() {
	// Collect credentials if required
	var (
		cpURL, cpUser, cpPassword string
		cpInsecure, verbose       bool
		err                       error
	)
	flag.StringVar(&cpURL, "cpanel-url", "", "The URL you use to access cPanel")
	flag.StringVar(&cpUser, "cpanel-username", "", "The username you use to access cPanel")
	flag.StringVar(&cpPassword, "cpanel-password", "", "The password you use to access cPanel")
	flag.BoolVar(&cpInsecure, "cpanel-insecure", false, "Whether the cPanel URL needs to be accessed insecurely")
	flag.BoolVar(&verbose, "verbose", false, "Whether to be very noisy")
	flag.Parse()

	// Choose cPanel API client based on environment
	if cpanel.IsLocal() {
		cpanelClient = cpanel.NewLocalClient()
	} else if cpanelClient, err = cpanel.NewRemoteClient(
		cpURL, cpUser, cpPassword, makeHTTPClient(cpInsecure)); err != nil {
		log.Fatalf("Couldn't create remote cPanel API client: %v.", err)
	}

	if verbose {
		cpanelCtx = context.WithValue(cpanelCtx, cpanel.LogRequestsAndResponses, true)
	}

	// Make sure that the cPanel account is valid and has the right features
	if err = ensureCpanelPrereqs(); err != nil {
		log.Fatalf("cPanel credentials did not work: %v.", err)
		return
	}
}

func ensureCpanelPrereqs() error {
	ctx, cancel := context.WithTimeout(cpanelCtx, 10*time.Second)
	defer cancel()

	features, err := cpanelClient.ListFeatures(ctx)
	if err != nil {
		return err
	}

	if !features.HasFeature("sslinstall") {
		return errors.New("cPanel account doesn't allow the installation of " +
			"SSL certificates, which prevents this program from working")
	}
	if !features.HasFeature("filemanager") {
		return errors.New("cPanel account doesn't allow use of the file manager, " +
			", which prevents this program from working")
	}
	return nil
}

func makeHTTPClient(insecure bool) *http.Client {
	if insecure {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}
	return http.DefaultClient
}

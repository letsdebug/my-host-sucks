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
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/letsdebug/my-host-sucks/cpanel"
)

var (
	api    cpanel.API
	apiCtx context.Context = context.Background()
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
		api = cpanel.NewLocalAPI()
	} else if api, err = cpanel.NewRemoteAPI(cpURL, cpUser, cpPassword, makeHTTPClient(cpInsecure)); err != nil {
		log.Fatalf("Couldn't create remote cPanel API client: %v. Make sure the details are correct.", err)
	}

	if verbose {
		apiCtx = context.WithValue(apiCtx, cpanel.LogRequestsAndResponses, true)
	}

	// Make sure we can talk to cPanel
	if err = testCpanel(); err != nil {
		log.Fatalf("cPanel credentials did not work: %v. Make sure the details are correct.", err)
		return
	}
}

func testCpanel() error {
	ctx, cancel := context.WithTimeout(apiCtx, 10*time.Second)
	defer cancel()
	if _, err := cpanel.DomainsData(ctx, api); err != nil {
		return err
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

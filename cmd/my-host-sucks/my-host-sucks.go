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
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"github.com/letsdebug/my-host-sucks/cpanel"
)

var state struct {
	CpanelURL      string
	CpanelUsername string
	CpanelPassword string
	CpanelInsecure bool

	api cpanel.API
}

func main() {
	// Collect credentials if required
	flag.StringVar(&state.CpanelURL, "cpanel-url", "", "The URL you use to access cPanel")
	flag.StringVar(&state.CpanelUsername, "cpanel-username", "", "The URL you use to access cPanel")
	flag.StringVar(&state.CpanelPassword, "cpanel-password", "", "The URL you use to access cPanel")
	flag.BoolVar(&state.CpanelInsecure, "cpanel-insecure", false, "Whether the cPanel URL needs to be accessed inscurely")
	flag.Parse()

	// Choose cPanel API client based on environment
	if cpanel.IsLocal() {
		state.api = cpanel.NewLocalAPI()
	} else {
		api, err := cpanel.NewRemoteAPI(state.CpanelURL, state.CpanelUsername, state.CpanelPassword, makeHTTPClient())
		if err != nil {
			log.Fatalf("Couldn't create remote cPanel API client: %v. Make sure the details are correct.", err)
		}
		state.api = api
	}

	// Make sure we can talk to cPanel
	if err := testCpanel(); err != nil {
		log.Fatalf("cPanel credentials did not work: %v. Make sure the details are correct.", err)
		return
	}
}

func testCpanel() error {
	if _, err := cpanel.DomainsData(state.api); err != nil {
		return err
	}
	return nil
}

func makeHTTPClient() *http.Client {
	if state.CpanelInsecure {
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

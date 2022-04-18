package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"

	"github.com/algorand/node-ui/messages"
	"github.com/algorand/node-ui/tui"
)

var uiPort uint64
var url string
var token string
var algodDataDir string

// TODO "r" to set the refresh rate

func init() {
	flag.Uint64Var(&uiPort, "a", 0, "Port address to host TUI from, set to 0 to run directly")
	flag.StringVar(&url, "u", "", "Algod URL and port formatted like localhost:1234")
	flag.StringVar(&token, "t", "", "Algod REST API Token")
	flag.StringVar(&algodDataDir, "d", "", "Path to algorand data directory, used to override ALGORAND_DATA environment variable")
}

func getRequestorOrExit() *messages.Requestor {
	// Initialize from -d, ALGORAND_DATA, or provided URL/Token
	if algodDataDir == "" {
		algodDataDir = os.Getenv("ALGORAND_DATA")
		if algodDataDir != "" {
			fmt.Println("Using ALGORAND_DATA environment variable.")
		}
	}

	// Lookup URL/Token
	if algodDataDir != "" {
		if url != "" || token != "" {
			fmt.Fprintln(os.Stderr, "Do not use -u/-t with -d or the ALGORAND_DATA environment variable.")
			os.Exit(1)
		}

		netpath := filepath.Join(algodDataDir, "algod.net")
		tokenpath := filepath.Join(algodDataDir, "algod.token")

		var netaddrbytes []byte
		netaddrbytes, err := ioutil.ReadFile(netpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read URL from file (%s): %s\n", netpath, err.Error())
			os.Exit(1)
		}
		url = string(netaddrbytes)
		netaddr := strings.TrimSpace(string(netaddrbytes))
		if !strings.HasPrefix(netaddr, "http") {
			netaddr = "http://" + netaddr
		}
		tokenBytes, err := ioutil.ReadFile(tokenpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read Token from file (%s): %s\n", tokenpath, err.Error())
			os.Exit(1)
		}
		token = string(tokenBytes)
	}

	if url == "" || token == "" {
		fmt.Fprintln(os.Stderr, "Must provide a way to get the algod REST API.")
		os.Exit(1)
	}

	client, err := algod.MakeClient(url, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem creating client connection: %s\n", err.Error())
		os.Exit(1)
	}

	return messages.MakeRequestor(client, algodDataDir)
}

func main() {
	flag.Parse()
	tui.Start(uiPort, getRequestorOrExit())
}

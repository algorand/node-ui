package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/types"

	"github.com/algorand/node-ui/messages"
	"github.com/algorand/node-ui/tui"
)

var command *cobra.Command

func main() {
	err := command.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem running command: %s\n", err.Error())
	}
}

// TODO "r" to set the refresh rate
type arguments struct {
	tuiPort          uint64
	algodURL         string
	algodToken       string
	algodDataDir     string
	addressWatchList []string
}

func run(args arguments) {
	request := getRequestorOrExit(args.algodDataDir, args.algodURL, args.algodToken)
	addresses := getAddressesOrExit(args.addressWatchList)
	tui.Start(args.tuiPort, request, addresses)
}

func init() {
	var args arguments

	command = &cobra.Command{
		Use:   "",
		Short: "Launch terminal user interface",
		Long:  "Node UI is a terminal user interface that displays information about a target algod instance.",
		Run: func(_ *cobra.Command, _ []string) {
			run(args)
		},
	}

	command.Flags().Uint64VarP(&args.tuiPort, "tui-port", "p", 0, "Port address to host TUI from, set to 0 to run directly.")
	command.Flags().StringVarP(&args.algodURL, "algod-url", "u", "", "Algod URL and port to monitor, formatted like localhost:1234.")
	command.Flags().StringVarP(&args.algodToken, "algod-token", "t", "", "Algod REST API token.")
	command.Flags().StringVarP(&args.algodDataDir, "algod-data-dir", "d", "", "Path to Algorand data directory, when set it overrides the ALGORAND_DATA environment variable.")
	command.Flags().StringArrayVarP(&args.addressWatchList, "watch-list", "w", nil, "Account addresses to watch in the accounts tab, may provide more than once to watch multiple accounts.")
}

func getRequestorOrExit(algodDataDir, url, token string) *messages.Requestor {
	// Initialize from -d, ALGORAND_DATA, or provided URL/Token

	if algodDataDir != "" && (url != "" || token != "") {
			fmt.Fprintln(os.Stderr, "Do not use -u/-t with -d.")
			os.Exit(1)
	}

	// If url/token are missing, attempt to use environment variable.
	if url == "" && token == "" {
		if algodDataDir == "" {
			algodDataDir = os.Getenv("ALGORAND_DATA")
			if algodDataDir != "" {
				fmt.Println("Using ALGORAND_DATA environment variable.")
			}
		}

		if algodDataDir == "" {
			fmt.Fprintln(os.Stderr, "Algod is not available.\nMust provide url and token with -u/-t or a data directory with -d or the ALGORAND_DATA environment variable.")
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
		url = strings.TrimSpace(string(netaddrbytes))
		tokenBytes, err := ioutil.ReadFile(tokenpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read Token from file (%s): %s\n", tokenpath, err.Error())
			os.Exit(1)
		}
		token = string(tokenBytes)
	}

	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
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

func getAddressesOrExit(addrs []string) (result []types.Address) {
	failed := false
	for _, addr := range addrs {
		converted, err := types.DecodeAddress(addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode address '%s': %s\n", addr, err.Error())
			failed = true
		}
		result = append(result, converted)
	}

	if failed {
		os.Exit(1)
	}

	return result
}

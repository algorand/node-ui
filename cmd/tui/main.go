package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/types"

	"github.com/algorand/node-ui/messages"
	"github.com/algorand/node-ui/tui"
	"github.com/algorand/node-ui/version"
)

func main() {
	err := makeCommand().Run(context.Background(), os.Args)
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
	versionFlag      bool
}

func run(args arguments) {
	if args.versionFlag {
		fmt.Println(version.LongVersion())
		os.Exit(0)
	}
	request := getRequestorOrExit(args.algodDataDir, args.algodURL, args.algodToken)
	addresses := getAddressesOrExit(args.addressWatchList)
	tui.Start(args.tuiPort, request, addresses)
}

func makeCommand() *cli.Command {
	var args arguments
	return &cli.Command{
		Name:  "node-ui",
		Usage: "Launch the Algorand Node UI.",
		Flags: []cli.Flag{
			&cli.Uint64Flag{
				Name:        "tui-port",
				Aliases:     []string{"p"},
				Usage:       "Port address to host TUI from, set to 0 to run directly.",
				Value:       0,
				Sources:     cli.EnvVars("TUI_PORT"),
				Destination: &args.tuiPort,
			},
			&cli.StringFlag{
				Name:        "algod-url",
				Aliases:     []string{"u"},
				Usage:       "Algod URL and port to monitor, formatted like localhost:1234.",
				Value:       "",
				Sources:     cli.EnvVars("ALGOD_URL"),
				Destination: &args.algodURL,
			},
			&cli.StringFlag{
				Name:        "algod-token",
				Aliases:     []string{"t"},
				Usage:       "Algod REST API token.",
				Value:       "",
				Sources:     cli.EnvVars("ALGOD_TOKEN"),
				Destination: &args.algodToken,
			},
			&cli.StringFlag{
				Name:        "algod-data-dir",
				Aliases:     []string{"d"},
				Usage:       "Path to Algorand data directory.",
				Value:       "",
				Sources:     cli.EnvVars("ALGORAND_DATA"),
				Destination: &args.algodDataDir,
			},
			&cli.StringSliceFlag{
				Name:        "watch-list",
				Aliases:     []string{"w"},
				Usage:       "Account addresses to watch in the accounts tab, may provide more than once to watch multiple accounts. Use comma separated values if providing more than one account with an environment variable.",
				Value:       nil,
				Sources:     cli.EnvVars("WATCH_LIST"),
				Destination: &args.addressWatchList,
			},
			&cli.BoolFlag{
				Name:        "version",
				Aliases:     []string{"v"},
				Usage:       "Print version information and exit.",
				Value:       false,
				Destination: &args.versionFlag,
			},
		},
		Action: func(c *cli.Context) error {
			run(args)
			return nil
		},
	}
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
		netaddrbytes, err := os.ReadFile(netpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read URL from file (%s): %s\n", netpath, err.Error())
			os.Exit(1)
		}
		url = strings.TrimSpace(string(netaddrbytes))
		tokenBytes, err := os.ReadFile(tokenpath)
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

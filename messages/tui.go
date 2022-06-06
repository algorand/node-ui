// Copyright (C) 2019-2022 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package messages

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Requestor provides an opaque pointer for an algod client.
type Requestor struct {
	Client  *algod.Client
	dataDir string
}

// MakeRequestor builds the requestor object.
func MakeRequestor(client *algod.Client, dataDir string) *Requestor {
	return &Requestor{
		Client:  client,
		dataDir: dataDir,
	}
}

// NetworkMsg holds network information.
type NetworkMsg struct {
	GenesisID   string
	GenesisHash types.Digest
	NodeVersion string
	Err         error
}

func formatVersion(ver models.Version) string {
	return fmt.Sprintf("%s %d.%d.%d (%s)",
		ver.Build.Channel,
		ver.Build.Major,
		ver.Build.Major,
		ver.Build.BuildNumber,
		ver.Build.CommitHash)
}

// GetNetworkCmd provides a tea.Cmd for fetching a NetworkMsg.
func (r Requestor) GetNetworkCmd() tea.Cmd {
	return func() tea.Msg {
		ver, err := r.Client.Versions().Do(context.Background())
		if err != nil {
			return NetworkMsg{
				Err: err,
			}
		}

		var digest types.Digest
		if len(ver.GenesisHash) != len(digest) {
			return NetworkMsg{
				Err: fmt.Errorf("unexpected genesis hash, wrong number of bytes"),
			}
		}
		copy(digest[:], ver.GenesisHash)

		return NetworkMsg{
			GenesisID:   ver.GenesisID,
			GenesisHash: digest,
			NodeVersion: formatVersion(ver),
		}
	}
}

// StatusMsg has node status information.
type StatusMsg struct {
	Status models.NodeStatus
	Error  error
}

// GetStatusCmd provides a tea.Cmd for fetching a StatusMsg.
func (r Requestor) GetStatusCmd() tea.Cmd {
	return func() tea.Msg {
		resp, err := r.Client.Status().Do(context.Background())
		//s, err := s.node.Status()
		return StatusMsg{
			Status: resp,
			Error:  err,
		}
	}
}

// AccountStatusMsg has account balance information.
type AccountStatusMsg struct {
	Balances map[types.Address]map[uint64]uint64
	Err      error
}

// GetAccountStatusCmd provides a tea.Cmd for fetching a AccountStatusMsg.
func (r Requestor) GetAccountStatusCmd(accounts []types.Address) tea.Cmd {
	return func() tea.Msg {
		var rval AccountStatusMsg
		rval.Balances = make(map[types.Address]map[uint64]uint64)

		for _, acct := range accounts {
			resp, err := r.Client.AccountInformation(acct.String()).Do(context.Background())
			if err != nil {
				return AccountStatusMsg{
					Err: err,
				}
			}
			rval.Balances[acct] = make(map[uint64]uint64)

			// algos at the special index
			rval.Balances[acct][0] = resp.Amount

			// everything else
			for _, holding := range resp.Assets {
				rval.Balances[acct][holding.AssetId] = holding.Amount
			}
		}

		return rval
	}
}

func doFastCatchupRequest(verb, network string) error {
	resp, err := http.Get(fmt.Sprintf("https://algorand-catchpoints.s3.us-east-2.amazonaws.com/channel/%s/latest.catchpoint", network))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	catchpoint := strings.Replace(string(body), "#", "%23", 1)

	//start fast catchup
	url := fmt.Sprintf("http://localhost:8080/v2/catchup/%s", catchpoint)
	url = url[:len(url)-1] // remove \n
	apiToken, err := os.ReadFile(path.Join(os.Getenv("ALGORAND_DATA"), "algod.admin.token"))
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(verb, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Algo-Api-Token", string(apiToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// StartFastCatchup attempts to start fast catchup for a given network.
func StartFastCatchup(network string) tea.Cmd {
	return func() tea.Msg {
		err := doFastCatchupRequest(http.MethodPost, network)
		if err != nil {
			panic(err)
		}
		return nil
	}
}

// StopFastCatchup attempts to stop fast catchup for a given network.
func StopFastCatchup(network string) tea.Cmd {
	return func() tea.Msg {
		err := doFastCatchupRequest(http.MethodDelete, network)
		if err != nil {
			panic(err)
		}
		return nil
	}
}

// GetConfigs returns the node config.json file if possible.
func GetConfigs() string {
	// TODO: Optional
	configs, err := os.ReadFile(path.Join(os.Getenv("ALGORAND_DATA"), "config.json"))
	if err != nil {
		return "config.json file not found"
	}
	return string(configs)
}

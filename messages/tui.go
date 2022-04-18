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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"

	tea "github.com/charmbracelet/bubbletea"
	//"github.com/algorand/go-algorand/crypto"
	//"github.com/algorand/go-algorand/data/basics"
	//"github.com/algorand/go-algorand/node"
)

var (
	AddressList = []string{
		"ZONEGRWBV3Q7JA6RHAN4EMAX6ICIVZX2C6U65DCNHYLIL4PUBB7O6DOSBI",
		"XFYAYSEGQIY2J3DCGGXCPXY5FGHSVKM3V4WCNYCLKDLHB7RYDBU233QB5M",
		"WNEJFT6HTAX3CQ6YOPIY65AKYCBQM6BLV4S5OP54VH76OP33LOL2MYGSIM",
		"GULDQIEZ2CUPBSHKXRWUW7X3LCYL44AI5GGSHHOQDGKJAZ2OANZJ43S72U",
		"57QZ4S7YHTWPRAM3DQ2MLNSVLAQB7DTK4D7SUNRIEFMRGOU7DMYFGF55BY",
		"ETGSQKACKC56JWGMDAEP5S2JVQWRKTQUVKCZTMPNUGZLDVCWPY63LSI3H4",
		"J4AEINCSSLDA7LNBNWM4ZXFCTLTOZT5LG3F5BLMFPJYGFWVCMU37EZI2AM",
	}
)

type NetworkMsg struct {
	GenesisID   string
	GenesisHash types.Digest
}

// TODO: client instead of server
func GetNetworkCmd() tea.Cmd {
	return func() tea.Msg {
		return NetworkMsg{
			GenesisID:   "test",
			GenesisHash: [32]byte{},
		}
	}
}

type StatusMsg struct {
	Status models.NodeStatusResponse
	Error  error
}

// TODO: client instead of server
func GetStatusCmd() tea.Cmd {
	return func() tea.Msg {
		//s, err := s.node.Status()
		return StatusMsg{
			Status: models.NodeStatusResponse{},
			Error:  nil,
		}
	}
}

type AccountStatusMsg map[types.Address]uint64

// TODO: client instead of server
func GetAccountStatusMsg() tea.Cmd {
	return func() tea.Msg {
		/*
			currentNode := GetNode(s)
			ledger := currentNode.Ledger()

			rval := make(AccountStatusMsg)

			for _, a := range AddressList {
				currentAddress, err := basics.UnmarshalChecksumAddress(a)

				if err != nil {
					continue
				}

				record, _, _, err := ledger.LookupLatest(currentAddress)

				rval[currentAddress] = record.MicroAlgos.Raw
			}

			return rval
		*/
		return nil
	}
}

// TODO: client instead of server
func StartFastCatchup(network string) tea.Cmd {
	return func() tea.Msg {
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
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("X-Algo-Api-Token", string(apiToken))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		/*
			s, err := s.node.Status()
			return StatusMsg{
				Status: s,
				Error:  err,
			}
		*/
		return nil
	}
}

// TODO: client instead of server
func StopFastCatchup(network string) tea.Cmd {
	return func() tea.Msg {
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
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("X-Algo-Api-Token", string(apiToken))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		return nil
	}
}

func GetConfigs() string {
	// TODO: Optional
	configs, err := os.ReadFile(path.Join(os.Getenv("ALGORAND_DATA"), "config.json"))
	if err != nil {
		return "config.json file not found"
	}
	return string(configs)
}

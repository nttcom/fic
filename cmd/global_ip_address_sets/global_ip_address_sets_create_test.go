/*
Copyright 2020 NTT Limited and NTT Communications Corporation All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package global_ip_address_sets

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/nttcom/fic/cmd/global_ip_address_sets/testdata"
	"github.com/nttcom/go-fic"
	th "github.com/nttcom/go-fic/testhelper"
	fakeClient "github.com/nttcom/go-fic/testhelper/client"
	"github.com/stretchr/testify/require"
)

func TestCreateCmd(t *testing.T) {
	const routerID = "F022000000013"
	const natID = "F052000000013"

	testCases := []struct {
		name        string
		clientError error
		args        []string
		flags       map[string]string
		request     string
		statusCode  int
		response    string
		output      string
		errorMsg    string
	}{
		{
			"success with destination NAT",
			nil,
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"router": routerID, "nat": natID, "type": "destinationNat", "num-of-addresses": "1"},
			testdata.CreateDstRequest,
			http.StatusOK,
			testdata.CreateDstResponse,
			testdata.CreateDstOutput,
			"",
		},
		{
			"success with source NAPT",
			nil,
			[]string{"testSrc&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"router": routerID, "nat": natID, "type": "sourceNapt", "num-of-addresses": "5"},
			testdata.CreateSrcRequest,
			http.StatusOK,
			testdata.CreateSrcResponse,
			testdata.CreateSrcOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"router": routerID, "nat": natID, "type": "destinationNat", "num-of-addresses": "1"},
			"",
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"global IP address set already exists",
			nil,
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"router": routerID, "nat": natID, "type": "destinationNat", "num-of-addresses": "1"},
			testdata.CreateDstRequest,
			http.StatusConflict,
			"",
			"",
			"calling Create global IP address set API:",
		},
		{
			"too many arguments",
			nil,
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901", "too-many-arguments"},
			map[string]string{"router": routerID, "nat": natID, "type": "destinationNat", "num-of-addresses": "1"},
			"",
			0,
			"",
			"",
			"accepts 1 arg(s), received 2",
		},
		{
			"invalid name",
			nil,
			[]string{"testDst!"},
			map[string]string{"router": routerID, "nat": natID, "type": "destinationNat", "num-of-addresses": "1"},
			"",
			0,
			"",
			"",
			"name must consist of half-width alphanumeric characters and some symbols &()-_ less than 64: received testDst!",
		},
		{
			"too long name(65)",
			nil,
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901_"},
			map[string]string{"router": routerID, "nat": natID, "type": "destinationNat", "num-of-addresses": "1"},
			"",
			0,
			"",
			"",
			"name must consist of half-width alphanumeric characters and some symbols &()-_ less than 64: received testDst&()-_0123456789012345678901234567890123456789012345678901_",
		},
		{
			"invalid type",
			nil,
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"router": routerID, "nat": natID, "type": "invalidType", "num-of-addresses": "1"},
			"",
			0,
			"",
			"",
			"type must be sourceNapt or destinationNat: received invalidType",
		},
		{
			"too many num-of-addresses with destinationNat",
			nil,
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"router": routerID, "nat": natID, "type": "destinationNat", "num-of-addresses": "2"},
			"",
			0,
			"",
			"",
			"num-of-addresses must be 1 with destinationNat type: received 2",
		},
		{
			"too many num-of-addresses with sourceNapt",
			nil,
			[]string{"testSrc&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"router": routerID, "nat": natID, "type": "sourceNapt", "num-of-addresses": "6"},
			"",
			0,
			"",
			"",
			"num-of-addresses must be [1-5] with sourceNapt type: received 6",
		},
		{
			"router flag is not set",
			nil,
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"nat": natID, "type": "destinationNat", "num-of-addresses": "1"},
			"",
			0,
			"",
			"",
			"required flag(s) \"router\" not set",
		},
		{
			"nat flag is not set",
			nil,
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"router": routerID, "type": "destinationNat", "num-of-addresses": "1"},
			"",
			0,
			"",
			"",
			"required flag(s) \"nat\" not set",
		},
		{
			"type flag is not set",
			nil,
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"router": routerID, "nat": natID, "num-of-addresses": "1"},
			"",
			0,
			"",
			"",
			"required flag(s) \"type\" not set",
		},
		{
			"num-of-addresses flag is not set",
			nil,
			[]string{"testDst&()-_0123456789012345678901234567890123456789012345678901"},
			map[string]string{"router": routerID, "nat": natID, "type": "destinationNat"},
			"",
			0,
			"",
			"",
			"required flag(s) \"num-of-addresses\" not set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.statusCode != 0 {
				th.SetupHTTP()
				defer th.TeardownHTTP()

				th.Mux.HandleFunc(
					fmt.Sprintf("/routers/%s/nats/%s/global-ip-address-sets", routerID, natID),
					func(w http.ResponseWriter, r *http.Request) {
						th.TestMethod(t, r, "POST")
						th.TestJSONRequest(t, r, tc.request)

						w.Header().Add("Content-Type", "application/json")
						w.WriteHeader(tc.statusCode)

						fmt.Fprintf(w, tc.response)
					},
				)
			}

			buf := new(bytes.Buffer)
			cmd := newCreateCmd(
				func() (client *fic.ServiceClient, err error) {
					return fakeClient.ServiceClient(), tc.clientError
				},
				buf,
			)
			cmd.SetArgs(tc.args)
			for key, value := range tc.flags {
				cmd.Flags().Set(key, value)
			}
			err := cmd.Execute()

			if tc.errorMsg == "" {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			}

			require.Equal(t, strings.TrimPrefix(tc.output, "\n"), buf.String())
		})
	}
}

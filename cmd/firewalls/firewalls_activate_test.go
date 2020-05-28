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

package firewalls

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/nttcom/fic/cmd/firewalls/testdata"
	"github.com/nttcom/go-fic"
	th "github.com/nttcom/go-fic/testhelper"
	fakeClient "github.com/nttcom/go-fic/testhelper/client"
	"github.com/stretchr/testify/require"
)

func TestActivateCmd(t *testing.T) {
	const firewallID = "F040123456789"
	const routerID = "F042000000118"
	userIPAddresses := "192.168.0.0/30,192.168.0.4/30,192.168.0.8/30,192.168.0.12/30"

	testCases := []struct {
		name        string
		clientError error
		flags       map[string]string
		request     string
		statusCode  int
		response    string
		output      string
		errorMsg    string
	}{
		{
			"success",
			nil,
			map[string]string{
				"router":            routerID,
				"user-ip-addresses": userIPAddresses,
			},
			testdata.ActivateRequest,
			http.StatusAccepted,
			testdata.ActivateResponse,
			testdata.ActivateOutput,
			"",
		},
		{
			"Validation: UserIPAddresses (Invalid length)",
			nil,
			map[string]string{
				"router":            routerID,
				"user-ip-addresses": "192.168.0.0/30,192.168.0.4/30,192.168.0.8/30",
			},
			"",
			0,
			"",
			"",
			"user-ip-address must have length of 4 or 8",
		},
		{
			"Validation: UserIPAddresses (Invalid format)",
			nil,
			map[string]string{
				"router":            routerID,
				"user-ip-addresses": "192.168.0.0/30,192.168.0.4/30,192.168.0.8/30,foobar",
			},
			"",
			0,
			"",
			"",
			"user-ip-addresses must be CIDR whose subnet mask is 30, e.g. 192.168.0.0/30: received foobar",
		},
		{
			"Validation: UserIPAddresses (Invalid subnet)",
			nil,
			map[string]string{
				"router":            routerID,
				"user-ip-addresses": "192.168.0.0/31,192.168.0.4/30,192.168.0.8/30,192.168.0.12/30",
			},
			"",
			0,
			"",
			"",
			"subnet mask of user-ip-addresses must be 30, e.g. 192.168.0.0/30: " +
				"received 192.168.0.0/31",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			map[string]string{
				"router":            routerID,
				"user-ip-addresses": userIPAddresses,
			},
			"",
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"resource not found",
			nil,
			map[string]string{
				"router":            routerID,
				"user-ip-addresses": userIPAddresses,
			},
			testdata.ActivateRequest,
			http.StatusNotFound,
			"",
			"",
			"calling Activate firewall API:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.statusCode != 0 {
				th.SetupHTTP()
				defer th.TeardownHTTP()

				th.Mux.HandleFunc(
					fmt.Sprintf("/routers/%s/firewalls/%s/activate",
						routerID, firewallID),
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
			cmd := newActivateCmd(
				func() (client *fic.ServiceClient, err error) {
					return fakeClient.ServiceClient(), tc.clientError
				},
				buf,
			)
			cmd.SetArgs([]string{firewallID})
			for k, v := range tc.flags {
				cmd.Flags().Set(k, v)
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

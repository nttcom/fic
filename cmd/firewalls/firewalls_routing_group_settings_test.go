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

func TestRoutingGroupSettingsCmd(t *testing.T) {
	const firewallID = "F040123456789"
	const routerID = "F042000000118"

	testCases := []struct {
		name        string
		clientError error
		flags       map[string]string
		statusCode  int
		response    string
		output      string
		errorMsg    string
	}{
		{
			"success",
			nil,
			map[string]string{"router": routerID},
			http.StatusOK,
			testdata.ShowResponse,
			testdata.RoutingGroupSettingsOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			map[string]string{"router": routerID},
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"resource not found",
			nil,
			map[string]string{"router": routerID},
			http.StatusNotFound,
			"",
			"",
			"calling Get firewall API:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.statusCode != 0 {
				th.SetupHTTP()
				defer th.TeardownHTTP()

				th.Mux.HandleFunc(
					fmt.Sprintf("/routers/%s/firewalls/%s",
						routerID, firewallID),
					func(w http.ResponseWriter, r *http.Request) {
						th.TestMethod(t, r, "GET")

						w.Header().Add("Content-Type", "application/json")
						w.WriteHeader(tc.statusCode)

						fmt.Fprintf(w, tc.response)
					},
				)
			}

			buf := new(bytes.Buffer)
			cmd := newRoutingGroupSettingsCmd(
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

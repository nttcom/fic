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

func TestShowCmd(t *testing.T) {
	const id = "3e262888541d404f97546af412ed1dab"
	const routerID = "F022000000013"
	const natID = "F052000000013"

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
			map[string]string{"router": routerID, "nat": natID},
			http.StatusOK,
			testdata.ShowResponse,
			testdata.ShowOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			map[string]string{"router": routerID, "nat": natID},
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"resource not found",
			nil,
			map[string]string{"router": routerID, "nat": natID},
			http.StatusNotFound,
			"",
			"",
			"calling Show global IP address set API:",
		},
		{
			"router flag is not set",
			nil,
			map[string]string{"nat": natID},
			0,
			"",
			"",
			"required flag(s) \"router\" not set",
		},
		{
			"nat flag is not set",
			nil,
			map[string]string{"router": routerID},
			0,
			"",
			"",
			"required flag(s) \"nat\" not set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			th.SetupHTTP()
			defer th.TeardownHTTP()

			th.Mux.HandleFunc(
				fmt.Sprintf("/routers/%s/nats/%s/global-ip-address-sets/%s", routerID, natID, id),
				func(w http.ResponseWriter, r *http.Request) {
					th.TestMethod(t, r, "GET")

					w.Header().Add("Content-Type", "application/json")
					w.WriteHeader(tc.statusCode)

					fmt.Fprintf(w, tc.response)
				},
			)

			buf := new(bytes.Buffer)
			cmd := newShowCmd(
				func() (client *fic.ServiceClient, err error) {
					return fakeClient.ServiceClient(), tc.clientError
				},
				buf,
			)
			cmd.SetArgs([]string{id})
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

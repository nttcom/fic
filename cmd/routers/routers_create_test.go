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

package routers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/nttcom/fic/cmd/routers/testdata"
	"github.com/nttcom/go-fic"
	th "github.com/nttcom/go-fic/testhelper"
	fakeClient "github.com/nttcom/go-fic/testhelper/client"
	"github.com/stretchr/testify/require"
)

func TestCreateCmd(t *testing.T) {
	const routerName = "DummyRouter"

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
			"success",
			nil,
			[]string{routerName},
			map[string]string{"area": "JPEAST", "redundant": "true", "user-ip-address": "10.100.0.0/27"},
			testdata.CreateRequest,
			http.StatusAccepted,
			testdata.CreateResponse,
			testdata.CreateOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			[]string{routerName},
			map[string]string{"area": "JPEAST", "redundant": "true", "user-ip-address": "10.100.0.0/27"},
			"",
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"router already exists",
			nil,
			[]string{routerName},
			map[string]string{"area": "JPEAST", "redundant": "true", "user-ip-address": "10.100.0.0/27"},
			testdata.CreateRequest,
			http.StatusConflict,
			"",
			"",
			"calling Create router API:",
		},
		{
			"too many arguments",
			nil,
			[]string{routerName, "too-many-arguments"},
			map[string]string{"area": "JPEAST", "redundant": "true", "user-ip-address": "10.100.0.0/27"},
			"",
			0,
			"",
			"",
			"accepts 1 arg(s), received 2",
		},
		{
			"user-ip-address is not cidr format",
			nil,
			[]string{routerName},
			map[string]string{"area": "JPEAST", "redundant": "true", "user-ip-address": "10.100.0.0<^v^>27"},
			"",
			0,
			"",
			"",
			"user-ip-address must be CIDR whose subnet mask is 27, e.g. 10.100.0.0/27: received 10.100.0.0<^v^>27",
		},
		{
			"subnet mask of user-ip-address is not 27",
			nil,
			[]string{routerName},
			map[string]string{"area": "JPEAST", "redundant": "true", "user-ip-address": "10.100.0.0/24"},
			"",
			0,
			"",
			"",
			"subnet mask of user-ip-address must be 27, e.g. 10.100.0.0/27: received 10.100.0.0/24",
		},
		{
			"area flag is not set",
			nil,
			[]string{routerName},
			map[string]string{"redundant": "true", "user-ip-address": "10.100.0.0/27"},
			"",
			0,
			"",
			"",
			"required flag(s) \"area\" not set",
		},
		{
			"user-ip-address flag is not set",
			nil,
			[]string{routerName},
			map[string]string{"area": "JPEAST", "redundant": "true"},
			"",
			0,
			"",
			"",
			"required flag(s) \"user-ip-address\" not set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.statusCode != 0 {
				th.SetupHTTP()
				defer th.TeardownHTTP()

				th.Mux.HandleFunc(
					fmt.Sprintf("/routers"),
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

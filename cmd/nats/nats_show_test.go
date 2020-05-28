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

package nats

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/nttcom/fic/cmd/nats/testdata"
	"github.com/nttcom/go-fic"
	th "github.com/nttcom/go-fic/testhelper"
	fakeClient "github.com/nttcom/go-fic/testhelper/client"
	"github.com/stretchr/testify/require"
)

func TestShowCmd(t *testing.T) {
	const routerID = "F022000000335"
	const natID = "F052000000335"

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
			map[string]string{
				"router": routerID,
			},
			http.StatusOK,
			testdata.ShowResponse,
			testdata.ShowOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			map[string]string{
				"router": routerID,
			},
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"resource not found",
			nil,
			map[string]string{
				"router": routerID,
			},
			http.StatusNotFound,
			"",
			"",
			"calling Show nat API:",
		},
		{
			"router flag is not specified",
			nil,
			map[string]string{},
			0,
			"",
			"",
			"required flag(s) \"router\" not set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.statusCode != 0 {
				th.SetupHTTP()
				defer th.TeardownHTTP()

				th.Mux.HandleFunc(
					fmt.Sprintf("/routers/%s/nats/%s", routerID, natID),
					func(w http.ResponseWriter, r *http.Request) {
						th.TestMethod(t, r, "GET")

						w.Header().Add("Content-Type", "application/json")
						w.WriteHeader(tc.statusCode)

						fmt.Fprintf(w, tc.response)
					},
				)
			}

			buf := new(bytes.Buffer)
			cmd := newShowCmd(
				func() (client *fic.ServiceClient, err error) {
					return fakeClient.ServiceClient(), tc.clientError
				},
				buf,
			)
			cmd.SetArgs([]string{natID})
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

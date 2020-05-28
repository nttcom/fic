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

package ports

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/nttcom/fic/cmd/ports/testdata"
	"github.com/nttcom/go-fic"
	th "github.com/nttcom/go-fic/testhelper"
	fakeClient "github.com/nttcom/go-fic/testhelper/client"
	"github.com/stretchr/testify/require"
)

func TestShowCmd(t *testing.T) {
	const id = "F012000000172"

	testCases := []struct {
		name        string
		clientError error
		statusCode  int
		response    string
		output      string
		errorMsg    string
	}{
		{
			"success",
			nil,
			http.StatusOK,
			testdata.ShowResponse,
			testdata.ShowOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"resource not found",
			nil,
			http.StatusNotFound,
			"",
			"",
			"calling Show port API:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			th.SetupHTTP()
			defer th.TeardownHTTP()

			th.Mux.HandleFunc(
				fmt.Sprintf("/ports/%s", id),
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

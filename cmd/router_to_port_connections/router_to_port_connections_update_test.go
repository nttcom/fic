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

package router_to_port_connections

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/nttcom/fic/cmd/router_to_port_connections/testdata"
	"github.com/nttcom/go-fic"
	th "github.com/nttcom/go-fic/testhelper"
	fakeClient "github.com/nttcom/go-fic/testhelper/client"
	"github.com/stretchr/testify/require"
)

func TestUpdateCmd(t *testing.T) {
	const connectionID = "F032000001027"

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
			"success in case of non paired-configuration",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-route-filter": "noRoute,fullRoute",
			},
			testdata.UpdateRequest,
			http.StatusAccepted,
			testdata.UpdateResponse,
			testdata.UpdateOutput,
			"",
		},
		{
			"success in case of paired-configuration",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-secondary":    "null,2",
				"source-route-filter": "noRoute,fullRoute",
			},
			testdata.UpdatePairedRequest,
			http.StatusAccepted,
			testdata.UpdatePairedResponse,
			testdata.UpdatePairedOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"too many arguments",
			nil,
			[]string{connectionID, "too-many-arguments"},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"accepts 1 arg(s), received 2",
		},
		{
			"too many comma separated item to source-primary",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10,2",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"source-primary must have format like <asPathPrepend.in>,<asPathPrepend.out>,<med.out>: received null,2,10,2",
		},
		{
			"too many comma separated item to source-secondary",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-secondary":    "null,2,1",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"source-secondary must have format like <asPathPrepend.in>,<asPathPrepend.out>: received null,2,1",
		},
		{
			"invalid value of primary.asPathPrepend.in",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      ",2,10",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"in source-primary, asPathPrepend.in must be one of [null 1 2 3 4 5]: received ",
		},
		{
			"invalid value of primary.asPathPrepend.out",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,11,10",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"in source-primary, asPathPrepend.out must be one of [null 1 2 3 4 5]: received 11",
		},
		{
			"invalid type of primary.med.out",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,null",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"in source-primary, med.out must be numeric value: received null",
		},
		{
			"invalid value of primary.med.out with single connection",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,50",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"in source-primary, med.out in single connection must be one of [10 20 30 40]: received 50",
		},
		{
			"invalid value of primary.med.out with paired connection",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,20",
				"source-secondary":    "null,2",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"in source-primary, med.out in paired connection must be one of [10 30]: received 20",
		},
		{
			"invalid value of secondary.asPathPrepend.in",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-secondary":    "0,2",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"in source-secondary, asPathPrepend.in must be one of [null 1 2 3 4 5]: received 0",
		},
		{
			"invalid value of secondary.asPathPrepend.out",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-secondary":    "null,8",
				"source-route-filter": "noRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"in source-secondary, asPathPrepend.out must be one of [null 1 2 3 4 5]: received 8",
		},
		{
			"invalid value of source.routeFilter.in",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-route-filter": "invalidRoute,fullRoute",
			},
			"",
			0,
			"",
			"",
			"routeFilter.in must be one of [fullRoute noRoute]: received invalidRoute",
		},
		{
			"invalid value of source.routeFilter.out",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-route-filter": "noRoute,invalidRoute",
			},
			"",
			0,
			"",
			"",
			"routeFilter.out must be one of [fullRoute fullRouteWithDefaultRoute noRoute]: received invalidRoute",
		},
		{
			"internal server error in case of non paired-connection",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-route-filter": "noRoute,fullRoute",
			},
			testdata.UpdateRequest,
			http.StatusInternalServerError,
			"",
			"",
			"calling Update router to port connection API:",
		},
		{
			"internal server error in case of paired-connection",
			nil,
			[]string{connectionID},
			map[string]string{
				"source-primary":      "null,2,10",
				"source-secondary":    "null,2",
				"source-route-filter": "noRoute,fullRoute",
			},
			testdata.UpdatePairedRequest,
			http.StatusInternalServerError,
			"",
			"",
			"calling Update router to port connection API:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.statusCode != 0 {
				th.SetupHTTP()
				defer th.TeardownHTTP()

				th.Mux.HandleFunc(
					fmt.Sprintf("/router-to-port-connections/%s", connectionID),
					func(w http.ResponseWriter, r *http.Request) {
						th.TestMethod(t, r, "PATCH")
						th.TestJSONRequest(t, r, tc.request)

						w.Header().Add("Content-Type", "application/json")
						w.WriteHeader(tc.statusCode)

						fmt.Fprintf(w, tc.response)
					},
				)
			}

			buf := new(bytes.Buffer)
			cmd := newUpdateCmd(
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

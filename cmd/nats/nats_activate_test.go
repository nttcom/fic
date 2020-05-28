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

func TestActivateCmd(t *testing.T) {
	const routerID = "F022000000335"
	const natID = "F052000000335"

	testCases := []struct {
		name        string
		clientError error
		args        []string
		flags       map[string]interface{}
		request     string
		statusCode  int
		response    string
		output      string
		errorMsg    string
	}{
		{
			"success",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-00000000000000000000000000000000000000000000000000000001,sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
			},
			testdata.ActivateRequest,
			http.StatusAccepted,
			testdata.ActivateResponse,
			testdata.ActivateOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
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
			[]string{natID, "too-many-arguments"},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"accepts 1 arg(s), received 2",
		},
		{
			"router flag is not specified",
			nil,
			[]string{natID},
			map[string]interface{}{
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"required flag(s) \"router\" not set",
		},
		{
			"one of user-ip-addresses is not cidr format",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0:30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"user-ip-addresses must be CIDR whose subnet mask is 30, e.g. 192.168.0.0/30: received 192.168.0.0:30",
		},
		{
			"ip address length is invalid",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"user-ip-address must have length of 4 or 8",
		},
		{
			"one of subnet mask of user-ip-addresseses is not 30",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/29", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"subnet mask of user-ip-addresses must be 30, e.g. 192.168.0.0/30: received 192.168.0.0/29",
		},
		{
			"one of element length of global-ip-address-sets is not 3",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"global-ip-address-sets must have format like <name>,<type>,<numberOfAddresses>, e.g. src-set-01,sourceNapt,5",
		},
		{
			"one of element length of global-ip-address-sets is not 3",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,AAA",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"numberOfAddresses of global-ip-address-sets must be numeric value",
		},
		{
			"one of type of global-ip-address-sets is not neither sourceNapt nor destinationNat",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceMapt,1",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"type of global-ip-address-sets must be either sourceNpat or destinationNat",
		},
		{
			"sourceNapt count exceeds 8",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					// src count = 9
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					// dest count = 30
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1", "dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1", "dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"maximum number of sourceNapt type is 8",
		},
		{
			"destinationNat count exceeds 30",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,1",
					// src count = 8
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					"src-set-01,sourceNapt,1",
					// dest count = 31
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"maximum number of destinationNat type is 30",
		},
		{
			"user-ip-addresses is not specified",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
			},
			testdata.ActivateRequestWithGlobalIPAddressSetsOnly,
			http.StatusAccepted,
			testdata.ActivateResponse,
			testdata.ActivateOutput,
			"",
		},
		{
			"global-ip-address-sets is not specified",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
			},
			"",
			0,
			"",
			"",
			"required flag(s) \"global-ip-address-sets\" not set",
		},
		{
			"resource not found",
			nil,
			[]string{""},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,1",
					"dst-set-01,destinationNat,1",
				},
			},
			testdata.ActivateRequest,
			http.StatusNotFound,
			"",
			"",
			"calling Activate nat API:",
		},
		{
			"name of global-ip-address-sets has length=0",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					",sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"name of global-ip-address-sets must be composed of alpha-numeric characters and & ( ) - _, and must have maximum length of 64 as well",
		},
		{
			"name of global-ip-address-sets length has length=65",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-000000000000000000000000000000000000000000000000000000001,sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"name of global-ip-address-sets must be composed of alpha-numeric characters and & ( ) - _, and must have maximum length of 64 as well",
		},
		{
			"name of global-ip-address-sets has invalid chars",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01@,sourceNapt,5",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"name of global-ip-address-sets must be composed of alpha-numeric characters and & ( ) - _, and must have maximum length of 64 as well",
		},
		{
			"sourceNapt rule has invalid numberOfAddresses 6",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,6",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"numberOfAddresses of global-ip-address-sets must be in the range of 1 to 5 in case sourceNapt",
		},
		{
			"sourceNapt rule has invalid numberOfAddresses 0",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,0",
					"dst-set-01,destinationNat,1",
				},
			},
			"",
			0,
			"",
			"",
			"numberOfAddresses of global-ip-address-sets must be in the range of 1 to 5 in case sourceNapt",
		},
		{
			"destinationNat rule has invalid numberOfAddresses 2",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"user-ip-addresses": []string{
					"192.168.0.0/30", "192.168.4.0/30", "192.168.8.0/30", "192.168.12.0/30",
					"192.168.16.0/30", "192.168.20.0/30", "192.168.24.0/30", "192.168.28.0/30",
				},
				"global-ip-address-sets": []string{
					"src-set-01,sourceNapt,5",
					"dst-set-01,destinationNat,2",
				},
			},
			"",
			0,
			"",
			"",
			"numberOfAddresses of global-ip-address-sets must be 1 in case destinationNat",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.statusCode != 0 {
				th.SetupHTTP()
				defer th.TeardownHTTP()

				th.Mux.HandleFunc(
					fmt.Sprintf("/routers/%s/nats/%s/activate", routerID, natID),
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
			cmd.SetArgs(tc.args)
			for key, value := range tc.flags {
				switch u := value.(type) {
				case []string:
					for _, v := range u {
						cmd.Flags().Set(key, v)
					}
				case string:
					cmd.Flags().Set(key, u)
				default:
					t.Fatalf("invalid type of flag is specified as: %s", key)
				}
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

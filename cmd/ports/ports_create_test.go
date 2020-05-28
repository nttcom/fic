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

func TestCreateCmd(t *testing.T) {
	const portName = "DummyPort"

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
			"success with number-of-vlans flag",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "16", "switch": "lxea02comnw1", "type": "1G", "vlan-ranges": ""},
			testdata.CreateNumOfVLANsRequest,
			http.StatusAccepted,
			testdata.CreateNumOfVLANsResponse,
			testdata.CreateNumOfVLANsOutput,
			"",
		},
		{
			"success with vlan-ranges flag",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "", "switch": "lxea02comnw1", "type": "1G", "vlan-ranges": "2145-2160"},
			testdata.CreateVLANRangesRequest,
			http.StatusAccepted,
			testdata.CreateVLANRangesResponse,
			testdata.CreateVLANRangesOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			[]string{portName},
			map[string]string{"number-of-vlans": "16", "switch": "lxea02comnw1", "type": "1G", "vlan-ranges": ""},
			"",
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"port already exists",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "", "switch": "lxea02comnw1", "type": "1G", "vlan-ranges": "2145-2160"},
			testdata.CreateVLANRangesRequest,
			http.StatusConflict,
			"",
			"",
			"calling Create port API:",
		},
		{
			"too many arguments",
			nil,
			[]string{portName, "too-many-arguments"},
			map[string]string{"number-of-vlans": "16", "switch": "lxea02comnw1", "type": "1G", "vlan-ranges": ""},
			"",
			0,
			"",
			"",
			"accepts 1 arg(s), received 2",
		},
		{
			"invalid port type",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "16", "switch": "lxea02comnw1", "type": "123456789G", "vlan-ranges": ""},
			"",
			0,
			"",
			"",
			"type must be 1G or 10G",
		},
		{
			"both of number-of-vlans and vlan-ranges are set",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "16", "switch": "lxea02comnw1", "type": "10G", "vlan-ranges": "2145-2160"},
			"",
			0,
			"",
			"",
			"one of number-of-vlans and vlan-ranges must be set",
		},
		{
			"number-of-vlans is less than 16",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "1", "switch": "lxea02comnw1", "type": "10G", "vlan-ranges": ""},
			"",
			0,
			"",
			"",
			"number-of-vlans must be greater than or equal to 16",
		},
		{
			"number-of-vlans is not multiple of 16",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "17", "switch": "lxea02comnw1", "type": "10G", "vlan-ranges": ""},
			"",
			0,
			"",
			"",
			"number-of-vlans must be multiple of 16",
		},
		{
			"length of vlan-ranges larger than 32",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "", "switch": "lxea02comnw1", "type": "10G", "vlan-ranges": "1-16,17-32,33-48,49-64,65-80,81-96,97-112,113-128,129-144,145-160,161-176,177-192,193-208,209-224,225-240,241-256,257-272,273-288,289-304,305-320,321-336,337-352,353-368,369-384,385-400,401-416,417-432,433-448,449-464,465-480,481-496,497-512,513-528"},
			"",
			0,
			"",
			"",
			"length of vlan-ranges must be less than or equal to 32: received 33",
		},
		{
			"number-of-vlans is invalid format",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "", "switch": "lxea02comnw1", "type": "10G", "vlan-ranges": "1-16-32"},
			"",
			0,
			"",
			"",
			"vlan-ranges element must be <start>-<end> format, e.g. 17-32 and 49-64: received 1-16-32",
		},
		{
			"start value of number-of-vlans is not integer",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "", "switch": "lxea02comnw1", "type": "10G", "vlan-ranges": "1.1-16"},
			"",
			0,
			"",
			"",
			"vlan-ranges element must be integer range, e.g. 17-32 and 49-64: received 1.1-16",
		},
		{
			"end value of number-of-vlans is not integer",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "", "switch": "lxea02comnw1", "type": "10G", "vlan-ranges": "1-16.9"},
			"",
			0,
			"",
			"",
			"vlan-ranges element must be integer range, e.g. 17-32 and 49-64: received 1-16.9",
		},
		{
			"vlan-ranges element is invalid range",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "", "switch": "lxea02comnw1", "type": "10G", "vlan-ranges": "1-15"},
			"",
			0,
			"",
			"",
			"vlan-ranges element must be from 16a+1 to 16a+16(a is natural number), e.g. 17-32 and 49-64: received 1-15",
		},
		{
			"no switch option",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "", "type": "10G", "vlan-ranges": "2145-2160"},
			"",
			0,
			"",
			"",
			"required flag(s) \"switch\" not set",
		},
		{
			"no type option",
			nil,
			[]string{portName},
			map[string]string{"number-of-vlans": "", "switch": "lxea02comnw1", "vlan-ranges": "2145-2160"},
			"",
			0,
			"",
			"",
			"required flag(s) \"type\" not set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.statusCode != 0 {
				th.SetupHTTP()
				defer th.TeardownHTTP()

				th.Mux.HandleFunc(
					fmt.Sprintf("/ports"),
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

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

var tooLongDestinationNATRules = "group_1,group_2," + strings.Repeat("dst-set-01:192.168.0.1/32,", 30) + "dst-set-01:192.168.0.1/32"

func TestUpdateCmd(t *testing.T) {
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
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			testdata.UpdateRequest,
			http.StatusAccepted,
			testdata.UpdateResponse,
			testdata.UpdateOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"router flag is not specified",
			errors.New("client auth error"),
			[]string{natID},
			map[string]interface{}{
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"required flag(s) \"router\" not set",
		},
		{
			"too many arguments",
			nil,
			[]string{natID, "too-many-arguments"},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"accepts 1 arg(s), received 2",
		},
		// source-napt-rules argument validations
		{
			"length of source-napt-rules is short",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"source-napt-rules must have format like from, to, sourceNaptRules.entries.0.then.0, ... ,sourceNaptRules.entries.0.then.N",
		},
		{
			"format of 1st element of source napt rules is not valid",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"first and second element of source-napt-rules must have format like group_[1|2|3|4]",
		},
		{
			"format of 2nd element of source napt rules is not valid",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"first and second element of source-napt-rules must have format like group_[1|2|3|4]",
		},
		{
			"length of source-napt-rules exceeds maximum length",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04,src-set-05,src-set-06,src-set-07,src-set-08,src-set-09",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"maximum length of entries of source-napt-rules is 8",
		},
		{
			"required argument source-napt-rules is not set",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"required flag(s) \"source-napt-rules\" not set",
		},
		// destination-nat-rules argument validations
		{
			"length of source-napt-rules is short",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"destination-nat-rules must have format like from, to, entries.0.match.destinationAddress:entries.0.then, ... ,entries.N.match.destinationAddress:entries.N.then",
		},
		{
			"format of 1st element of destination nat rules is not valid",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"first and second element of destination-nat-rules must have format like group_[1|2|3|4]",
		},
		{
			"format of 2nd element of destination nat rules is not valid",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"first and second element of destination-nat-rules must have format like group_[1|2|3|4]",
		},
		{
			"length of destination-nat-rules exceeds maximum length",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					tooLongDestinationNATRules,
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"maximum length of entries of destination-nat-rules is 30",
		},
		{
			"format of 3rd or element of is not correct",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03@192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"entries of destination-nat-rules must have format like entries.0.match.destinationAddress:entries.0.then, ... ,entries.N.match.destinationAddress:entries.N.then",
		},
		{
			"destination address is not cidr format",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3@32,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"destinationAddress must be CIDR whose subnet mask is 32, e.g. 192.168.0.1/32: received 192.168.0.3@32",
		},
		{
			"subnet mask length of destination address is not 32",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.1/31,dst-set-04:192.168.0.4/32",
				},
			},
			"",
			0,
			"",
			"",
			"subnet mask of destinationAddress must be 32, e.g. 192.168.0.1/32: received 192.168.0.1/31",
		},
		{
			"required argument destination-nat-rules is not set",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
			},
			"",
			0,
			"",
			"",
			"required flag(s) \"destination-nat-rules\" not set",
		},
		{
			"resource not found",
			nil,
			[]string{natID},
			map[string]interface{}{
				"router": routerID,
				"source-napt-rules": []string{
					"group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
					"group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08",
				},
				"destination-nat-rules": []string{
					"group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
					"group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
				},
			},
			testdata.UpdateRequest,
			http.StatusNotFound,
			"",
			"",
			"calling Update nat API:",
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
						th.TestMethod(t, r, "PUT")
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

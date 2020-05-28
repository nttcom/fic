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

func TestCreateCmd(t *testing.T) {
	const connectionName = "testConnection"

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
			"success in case of non paired-connection",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,101,10.0.0.2/30,1",
				"bandwidth":           "10M",
			},
			testdata.CreateRequest,
			http.StatusAccepted,
			testdata.CreateResponse,
			testdata.CreateOutput,
			"",
		},
		{
			"success in case of paired-connection",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,3300,10.0.0.6/30,65535",
				"bandwidth":             "10M",
			},
			testdata.CreatePairedRequest,
			http.StatusAccepted,
			testdata.CreatePairedResponse,
			testdata.CreatePairedOutput,
			"",
		},
		{
			"client auth error",
			errors.New("client auth error"),
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
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
			[]string{connectionName, "too-many-arguments"},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"accepts 1 arg(s), received 2",
		},
		{
			"too long connection name",
			nil,
			[]string{"testConnection_00000000000000000000000000000000000000000000000001"},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"name of router to port connection must be composed of alpha-numeric characters " +
				"and & ( ) - _, and must have maximum length of 64 as well: received testConnection_00000000000000000000000000000000000000000000000001",
		},
		{
			"invalid connection name",
			nil,
			[]string{"test@connection"},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"name of router to port connection must be composed of alpha-numeric characters " +
				"and & ( ) - _, and must have maximum length of 64 as well: received test@connection",
		},
		{
			"invalid group name",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_10",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"source-group must be one of [group_1 group_2 group_3 group_4 group_5 group_6 group_7 group_8]: received group_10",
		},
		{
			"invalid format of source-primary",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"source-primary must have format like <ipAddress>,<asPathPrepend.in>,<asPathPrepend.out>,<med.out>: received 10.0.0.1/30",
		},
		{
			"invalid format of destination-primary",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"destination-primary must have format like <portId>,<vlan>,<ipAddress>,<asn>: received F012000000220,161,10.0.0.2/30",
		},
		{
			"invalid format of source-secondary",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30",
				"source-route-filter":   "fullRout,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6/30,65000",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"source-secondary must have format like <ipAddress>,<asPathPrepend.in>,<asPathPrepend.out>: received 10.0.0.5/30",
		},
		{
			"invalid format of destination-secondary",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6/30",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"destination-secondary must have format like <portId>,<vlan>,<ipAddress>,<asn>: received F012000000219,190,10.0.0.6/30",
		},
		{
			"invalid ip address of source-primary",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      ",4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"in source-primary, ipAddress must be CIDR whose subnet mask is 30, e.g. 10.0.0.1/30: received ",
		},
		{
			"invalid subnet mask of source-primary",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/20,,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"in source-primary, subnet mask of ipAddress must be 30, e.g. 10.0.0.1/30: received 10.0.0.1/20",
		},
		{
			"invalid value of primary.asPathPrepend.in",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,9,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"in source-primary, asPathPrepend.in must be one of [null 1 2 3 4 5]: received 9",
		},
		{
			"invalid value of primary.asPathPrepend.out",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,10,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"in source-primary, asPathPrepend.out must be one of [null 1 2 3 4 5]: received 10",
		},
		{
			"invalid type of primary.med.out",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,null",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
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
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,50",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
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
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,20",
				"source-secondary":      "10.0.0.5/30,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6/30,65000",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in source-primary, med.out in paired connection must be one of [10 30]: received 20",
		},
		{
			"invalid ip address of source-secondary",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6/30,65000",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in source-secondary, ipAddress must be CIDR whose subnet mask is 30, e.g. 10.0.0.1/30: received 10.0.0.5",
		},
		{
			"invalid subnet mask of source-secondary",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/20,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6/30,65000",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in source-secondary, subnet mask of ipAddress must be 30, e.g. 10.0.0.1/30: received 10.0.0.5/20",
		},
		{
			"invalid value of secondary.asPathPrepend.in",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,10,2",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6/30,65000",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in source-secondary, asPathPrepend.in must be one of [null 1 2 3 4 5]: received 10",
		},
		{
			"invalid value of secondary.asPathPrepend.out",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,2,noValue",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6/30,65000",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in source-secondary, asPathPrepend.out must be one of [null 1 2 3 4 5]: received noValue",
		},
		{
			"invalid format of source-route-filter",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute,noRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"source-route-filter must have format like <routeFilter.in>,<routeFilter.out>: received fullRoute,fullRouteWithDefaultRoute,noRoute",
		},
		{
			"invalid value of source.routeFilter.in",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "0,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"routeFilter.in must be one of [fullRoute noRoute]: received 0",
		},
		{
			"invalid value of source.routeFilter.out",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "noRoute,invalidRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"routeFilter.out must be one of [fullRoute fullRouteWithDefaultRoute noRoute]: received invalidRoute",
		},
		{
			"vlan of destination-primary is not numeric value",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "noRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,vlan,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-primary, vlan must be numeric value: received vlan",
		},
		{
			"vlan of destination-primary is not range of 101 to 3300",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "noRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,100,10.0.0.2/30,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-primary, vlan must be range of 101 to 3300: received 100",
		},
		{
			"ip address of destination-primary is invalid",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "noRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-primary, ipAddress must be CIDR whose subnet mask is 30, e.g. 10.0.0.2/30: received 10.0.0.2",
		},
		{
			"ip address of destination-primary is invalid cidr",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "noRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/20,65000",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-primary, subnet mask of ipAddress must be 30, e.g. 10.0.0.2/30: received 10.0.0.2/20",
		},
		{
			"asn of destination-primary is not numeric value",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "noRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,asn",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-primary, asn must be numeric value: received asn",
		},
		{
			"asn of destination-primary is not range of 1 to 65535",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "noRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,0",
				"bandwidth":           "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-primary, asn must be range of 1 to 65535: received 0",
		},
		{
			"vlan of destination-secondary is not numeric value",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,vlan,10.0.0.6/30,65000",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-secondary, vlan must be numeric value: received vlan",
		},
		{
			"vlan of destination-secondary is not range of 101 to 3300",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,3301,10.0.0.6/30,65000",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-secondary, vlan must be range of 101 to 3300: received 3301",
		},
		{
			"ip address of destination-secondary is invalid",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6,65000",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-secondary, ipAddress must be CIDR whose subnet mask is 30, e.g. 10.0.0.2/30: received 10.0.0.6",
		},
		{
			"ip address of destination-secondary is invalid cidr",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6/20,65000",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-secondary, subnet mask of ipAddress must be 30, e.g. 10.0.0.2/30: received 10.0.0.6/20",
		},
		{
			"asn of destination-secondary is not numeric value",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6/30,asn",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-secondary, asn must be numeric value: received asn",
		},
		{
			"asn of destination-secondary is not range of 1 to 65535",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,190,10.0.0.6/30,65536",
				"bandwidth":             "10M",
			},
			"",
			0,
			"",
			"",
			"in destination-secondary, asn must be range of 1 to 65535: received 65536",
		},
		{
			"invalid bandwidth",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,161,10.0.0.2/30,65000",
				"bandwidth":           "1000M",
			},
			"",
			0,
			"",
			"",
			"bandwidth must be one of [10M 20M 30M 40M 50M 100M 200M 300M 400M 500M 1G 2G 3G 4G 5G 10G]: received 1000M",
		},
		{
			"internal server error in case of non paired-connection",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":       "F022000000168",
				"source-group":        "group_1",
				"source-primary":      "10.0.0.1/30,4,4,10",
				"source-route-filter": "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary": "F012000000220,101,10.0.0.2/30,1",
				"bandwidth":           "10M",
			},
			testdata.CreateRequest,
			http.StatusInternalServerError,
			"",
			"",
			"calling Create router to port connection API:",
		},
		{
			"internal server error in case of paired-connection",
			nil,
			[]string{connectionName},
			map[string]string{
				"source-router":         "F022000000168",
				"source-group":          "group_1",
				"source-primary":        "10.0.0.1/30,4,4,10",
				"source-secondary":      "10.0.0.5/30,2,1",
				"source-route-filter":   "fullRoute,fullRouteWithDefaultRoute",
				"destination-primary":   "F012000000220,161,10.0.0.2/30,65000",
				"destination-secondary": "F012000000219,3300,10.0.0.6/30,65535",
				"bandwidth":             "10M",
			},
			testdata.CreatePairedRequest,
			http.StatusInternalServerError,
			"",
			"",
			"calling Create router to port connection API:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.statusCode != 0 {
				th.SetupHTTP()
				defer th.TeardownHTTP()

				th.Mux.HandleFunc(
					fmt.Sprintf("/router-to-port-connections"),
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

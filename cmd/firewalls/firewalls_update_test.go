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

func TestUpdateCmd(t *testing.T) {
	const (
		firewallID = "F040123456789"
		routerID   = "F042000000118"
		rules      = "group_1,group_2,rule-01,group1_addset_1:group1_addset_2," +
			"group2_addset_1:group2_addset_2,app_set_1,permit"
		customApplications    = "google-drive-web,tcp,443"
		applicationSets       = "app_set_1,google-drive-web:pre-defined-ftp"
		routingGroupSettings1 = "group_1,group1_addset_1,172.18.1.0/24"
		routingGroupSettings2 = "group_2,group1_addset_1,192.168.1.0/24"
		routingGroupSettings3 = "group_3,group1_addset_1,"
	)

	duplicateProps := func(length int, valueFmt string) []string {
		var r []string
		for i := 0; i < length; i++ {
			r = append(r, fmt.Sprintf(valueFmt, i))
		}
		return r
	}

	testCases := []struct {
		name                       string
		clientError                error
		flags                      map[string]interface{}
		request                    string
		statusCode                 int
		response, output, errorMsg string
	}{
		{
			"Success",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			testdata.UpdateRequest,
			http.StatusAccepted,
			testdata.UpdateResponse,
			testdata.UpdateOutput,
			"",
		},
		{
			"Success: with DestinationPort in range",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": "google-drive-web,tcp,8080-8089",
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			testdata.UpdateRequestWithPortRange,
			http.StatusAccepted,
			testdata.UpdateResponseWithPortRange,
			testdata.UpdateOutput,
			"",
		},
		{
			"Success: Multiple Rules",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": []string{
					rules,
					"group_1,group_2,rule-01,group1_addset_1:group1_addset_2," +
						"group2_addset_1:group2_addset_2,google-drive-web,permit",
				},
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			testdata.UpdateRequestMultipleRules,
			http.StatusAccepted,
			testdata.UpdateResponseMultipleRules,
			testdata.UpdateOutput,
			"",
		},
		{
			"Success: Multiple CustomApplications",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules":  rules,
				"custom-applications": []string{
					customApplications,
					"my-web-app,tcp,443",
				},
				"application-sets": applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			testdata.UpdateRequestMultipleCustomApplications,
			http.StatusAccepted,
			testdata.UpdateResponseMultipleCustomApplications,
			testdata.UpdateOutput,
			"",
		},
		{
			"Success: Multiple ApplicationSets",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets": []string{
					applicationSets,
					"app_set_2,my-web-app:azure-virtual-instance",
				},
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			testdata.UpdateRequestMultipleApplicationSets,
			http.StatusAccepted,
			testdata.UpdateResponseMultipleApplicationSets,
			testdata.UpdateOutput,
			"",
		},
		{
			"Validation: Duplicated application names",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules":  rules,
				"custom-applications": []string{
					customApplications,
					customApplications,
				},
				"application-sets": applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"all of custom-application-names and application-set-names must be unique among them",
		},
		{
			"Validation: Invalid Rules format (fewer length)",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group1,group_2,rule-01,group1_addset_1:group1_addset_2," +
					"group2_addset_1:group2_addset_2,app_set_1",
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			rulesFormatError,
		},
		{
			"Validation: Invalid Rules format (wrong length)",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group_1,group_2,rule-01,group1_addset_1:group1_addset_2," +
					"group2_addset_1:group2_addset_2,google-drive-web,permit,group_1,group_2",
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			rulesFormatError,
		},
		{
			"Validation: Rules.From",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group1,group_2,rule-01,group1_addset_1:group1_addset_2," +
					"group2_addset_1:group2_addset_2,app_set_1,permit",
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"from in rules must be either one of following: [group_1 group_2 " +
				"group_3 group_4]: received 'group1'",
		},
		{
			"Validation: Rules.To",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group_1,グループ2,rule-01,group1_addset_1:group1_addset_2," +
					"group2_addset_1:group2_addset_2,app_set_1,permit",
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"to in rules must be either one of following: [group_1 group_2 " +
				"group_3 group_4]: received 'グループ2'",
		},
		{
			"Validation: Too many Rules.Entries",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group_1,group_2," +
					func(length int) string {
						e := "rule-01,group1_addset_1:group1_addset_2,group2_addset_1:" +
							"group2_addset_2,app_set_1,permit"
						return strings.Repeat(e+",", length-1) + e
					}(51),
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"length of entry in rules must be less than or equal to 50: received 51",
		},
		{
			"Validation: Rules.Entries.Name (Invalid characters)",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group_1,group_2,rule*01,group1_addset_1:group1_addset_2," +
					"group2_addset_1:group2_addset_2,app_set_1,permit",
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"only alphanumeric characters, hyphens, and underscores up to 20 " +
				"characters are allowed for name of entries in rule: received 'rule*01'",
		},
		{
			"Validation: Rules.Entries.Name (Too many length)",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group_1,group_2,rulerulerulerulerule01,group1_addset_1:group1_addset_2," +
					"group2_addset_1:group2_addset_2,app_set_1,permit",
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"only alphanumeric characters, hyphens, and underscores up to 20 " +
				"characters are allowed for name of entries in rule: received 'rulerulerulerulerule01'",
		},
		{
			"Validation: Rules.Entries.Match.SourceAddressSets",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group_1,group_2,rule-01," +
					func(length int) string {
						e := "group1_addset_1"
						return strings.Repeat(e+":", length-1) + e
					}(11) + ",group2_addset_1:group2_addset_2,app_set_1,permit",
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"length of source-address-sets of entries in rule must be less than or equal to 10: received 11",
		},
		{
			"Validation: Rules.Entries.Match.DestinationAddressSets",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group_1,group_2,rule-01,group1_addset_1:group1_addset_2," +
					func(length int) string {
						e := "group1_addset_1"
						return strings.Repeat(e+":", length-1) + e
					}(11) + ",app_set_1,permit",
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"length of destination-address-sets of entries in rule must be less than or equal to 10: received 11",
		},
		{
			"Validation: Rules.Entries.Match.Application",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group_1,group_2,rule-01,group1_addset_1:group1_addset_2," +
					"group2_addset_1:group2_addset_2,app_set_10,permit",
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"application of entries in rule must be either one of following: " +
				"custom-application, application, 'pre-defined-application', or 'any': received 'app_set_10'",
		},
		{
			"Validation: Rules.Entries.Action",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group_1,group_2,rule-01,group1_addset_1:group1_addset_2," +
					"group2_addset_1:group2_addset_2,app_set_1,foo",
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"action of entries in rule must be either one of following: [permit " +
				"deny deny-with-logging]: received 'foo'",
		},
		{
			"Validation: CustomApplications (Length)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": duplicateProps(21, "myapp-%v,tcp,443"),
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"length of custom-applications must be less than or equal to 20: received 21",
		},
		{
			"Validation: CustomApplications (Format)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": "google-drive-web,tcp",
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			customApplicationsFormatError,
		},
		{
			"Validation: CustomApplications.Name (Invalid characters)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": "app#1,tcp,443",
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"only alphanumeric characters, hyphens, underscores, ampersands, " +
				"and parenthesis up to 64 characters are allowed for name of custom-applications: received app#1",
		},
		{
			"Validation: CustomApplications.Name (Invalid keywords)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": "app-any-1,tcp,443",
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"only alphanumeric characters, hyphens, underscores, ampersands, " +
				"and parenthesis up to 64 characters are allowed for name of custom-applications: received app-any-1",
		},
		{
			"Validation: CustomApplications.Protocol",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": "google-drive-web,http,443",
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"protocol for custom-applications must be one of following:[tcp " +
				"udp]: received 'http'",
		},
		{
			"Validation: CustomApplications.DestinationPort",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": "google-drive-web,tcp,123-456-7890",
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"destination-port for custom-applications must be either a port number or " +
				"port number range (xxx-yyy): received '123-456-7890",
		},
		{
			"Validation: CustomApplications.DestinationPort",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": "google-drive-web,tcp,123456",
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"destination-port for custom-applications must be either a port number or " +
				"port number range (xxx-yyy): received '123456",
		},
		{
			"Validation: ApplicationSets (Length)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    duplicateProps(3, "app_set_%v,google-drive-web:pre-defined-ftp"),
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"length of application-sets must be up to 2: received 3",
		},
		{
			"Validation: ApplicationSets (Format)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    "app_set_1",
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			applicationSetsFormatError,
		},
		{
			"Validation: ApplicationSets.Name",
			nil,
			map[string]interface{}{
				"router": routerID,
				"rules": "group_1,group_2,rule-01,group1_addset_1:" +
					"group1_addset_2,group2_addset_1:group2_addset_2,pre-defined,permit",
				"custom-applications": customApplications,
				"application-sets":    "pre-defined,google-drive-web:pre-defined-ftp",
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"only alphanumeric characters, hyphens, underscores, ampersands, " +
				"and parenthesis up to 64 characters are allowed for name " +
				"of application-sets: received 'pre-defined'",
		},
		{
			"Validation: ApplicationSets.Applications",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    "app_set_1,app0:app1:app2:app3:app4:app5:app6:app7:app8:app9:app10",
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"length of applications in application-sets must be less than or equal to 10: received 11",
		},
		{
			"Validation: RoutingGroupSettings (Length)",
			nil,
			map[string]interface{}{
				"router":                 routerID,
				"rules":                  rules,
				"custom-applications":    customApplications,
				"application-sets":       applicationSets,
				"routing-group-settings": duplicateProps(41, "group_%v,group1_addset_1,172.18.1.0/24"),
			},
			"",
			0,
			"",
			"",
			"length of routing-group-settings must be up to 40: received 41",
		},
		{
			"Validation: RoutingGroupSettings (Format)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					"group_1,group1_addset_1",
				},
			},
			"",
			0,
			"",
			"",
			routingGroupSettingsFormatError,
		},
		{
			"Validation: RoutingGroupSettings.GroupName",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					"group1,group1_addset_1,172.18.1.0/24",
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"name in routing-group-settings must be either one of following: " +
				"[group_1 group_2 group_3 group_4]: received 'group1'",
		},
		{
			"Validation: RoutingGroupSettings.AddressSets",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					"group_1,as1,172.18.1.0/24,as1,172.18.1.0/24,as1,172.18.1.0/24,as1," +
						"172.18.1.0/24,as1,172.18.1.0/24,as1,172.18.1.0/24",
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"length of address-sets in routing-group-settings must be up to 5: received 6",
		},
		{
			"Validation: RoutingGroupSettings.AddressSets.Name (Invalid characters)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					"group_1,as#1,172.18.1.0/24",
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"only alphanumeric characters, hyphens, underscores, ampersands, " +
				"and parenthesis up to 64 characters are allowed for name of " +
				"address-sets in routing-group-settings: received 'as#1'",
		},
		{
			"Validation: RoutingGroupSettings.AddressSets.Name (Too many length)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					"group_1,abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde,172.18.1.0/24",
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"only alphanumeric characters, hyphens, underscores, ampersands, " +
				"and parenthesis up to 64 characters are allowed for name of " +
				"address-sets in routing-group-settings: received " +
				"'abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde'",
		},
		{
			"Validation: RoutingGroupSettings.AddressSets.Addresses (Too many length)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					"group_1,as1,172.18.1.0/24:" +
						"172.18.1.0/24:172.18.1.0/24:172.18.1.0/24:172.18.1.0/24:172.18.1.0/24:" +
						"172.18.1.0/24:172.18.1.0/24:172.18.1.0/24:172.18.1.0/24:172.18.1.0/24",
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"addresses of address-sets in routing-group-settings must be up to 10: received 11",
		},
		{
			"Validation: RoutingGroupSettings.AddressSets.Addresses (Invalid CIDR)",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					"group_1,group1_addset_1,172.18.1.0/2400",
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"addresses of address-sets in routing-group-settings, e.g. " +
				"192.168.1.0/24: received 172.18.1.0/2400",
		},
		{
			"Client with error",
			errors.New("client auth error"),
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			"",
			0,
			"",
			"",
			"creating FIC client:",
		},
		{
			"resource not found",
			nil,
			map[string]interface{}{
				"router":              routerID,
				"rules":               rules,
				"custom-applications": customApplications,
				"application-sets":    applicationSets,
				"routing-group-settings": []string{
					routingGroupSettings1,
					routingGroupSettings2,
					routingGroupSettings3,
				},
			},
			testdata.UpdateRequest,
			http.StatusNotFound,
			"",
			"",
			"calling Update firewalls API:",
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
			cmd.SetArgs([]string{firewallID})
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

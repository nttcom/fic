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
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/firewalls"
	"github.com/spf13/cobra"
)

const rulesFormatError = "rules must have format like " +
	"from,to,entries.0.name," +
	"entries.0.match.sourceAddressSet.0:...:" +
	"entries.0.match.sourceAddressSet.N," +
	"entries.0.match.destinationAddressSet.0:...:" +
	"entries.0.match.destinationAddressSet.N," +
	"entries.0.match.application," +
	"entries.0.action,...," +
	"entries.N.match.sourceAddressSet.0:...:" +
	"entries.N.match.sourceAddressSet.N," +
	"entries.N.match.destinationAddressSet.0:...:" +
	"entries.N.match.destinationAddressSet.N," +
	"entries.N.match.application," +
	"entries.N.action"
const customApplicationsFormatError = "custom-applications must have format like " +
	"name,protocol,destinationPort"
const applicationSetsFormatError = "application-sets must have format like " +
	"name,applications.0:...:applications.N"
const routingGroupSettingsFormatError = "routing-group-settings must have format like " +
	"groupName,addressSets.0.name,addressSets.0.addresses.0:...:addressSets.0.addresses.N,...," +
	"groupName,addressSets.N.name,addressSets.N.addresses.0:...:addressSets.N.addresses.N"

// newActivateCmd creates a new `fic firewalls activate` command
func newUpdateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var (
		routerID, firewallID                        string
		rawRules, rawCustomApplications             []string
		rawApplicationSets, rawRoutingGroupSettings []string
	)

	validGroupNames := []string{"group_1", "group_2", "group_3", "group_4"}
	validActions := []string{"permit", "deny", "deny-with-logging"}
	validProtocols := []string{"tcp", "udp"}

	matcherEntryNames := regexp.MustCompile("^[\\w\\d-_]{1,20}$")
	matcherKeywords := regexp.MustCompile("^[\\w\\d-_&()]{1,64}$")
	matcherPorts := regexp.MustCompile("^(\\d{1,5}|\\d{1,5}-\\d{1,5})$")

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update firewall",
		Example: "fic firewalls update F040123456789 --router F022000000335 " +
			"--rules group_1,group_2,rule-01,group1_addset_1:group1_addset_2," +
			"group2_addset_1:group2_addset_2,app_set_1,permit " +
			"--custom-applications google-drive-web,tcp,443 " +
			"--application-sets app_set_1,google-drive-web:pre-defined-ftp " +
			"--routing-group-settings group_1,group1_addset_1,172.18.1.0/24:192.168.1.0/24",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				rules                []firewalls.Rule
				customApplications   []firewalls.CustomApplication
				applicationSets      []firewalls.ApplicationSet
				routingGroupSettings []firewalls.RoutingGroupSetting
			)
			firewallID = args[0]

			for _, v := range rawRules {
				var rule firewalls.Rule
				rawRule := strings.Split(v, ",")
				if length := len(rawRule); length < 7 || (length-2)%5 != 0 {
					return fmt.Errorf(rulesFormatError)
				}

				rule.From = rawRule[0]
				if !utils.StringInSlice(rule.From, validGroupNames) {
					return fmt.Errorf("from in rules must be either one of following: "+
						"%s: received '%s'", validGroupNames, rule.From)
				}
				rule.To = rawRule[1]
				if !utils.StringInSlice(rule.To, validGroupNames) {
					return fmt.Errorf("to in rules must be either one of following: "+
						"%s: received '%s'", validGroupNames, rule.To)
				}
				lengthEntries := (len(rawRule) - 2) / 5
				if lengthEntries > 50 {
					return fmt.Errorf("length of entry in rules must be less than "+
						"or equal to 50: received %d", lengthEntries)
				}
				for j := 0; j < lengthEntries; j++ {
					var (
						name, application, action                 string
						sourceAddressSets, destinationAddressSets []string
					)
					name = rawRule[j*5+2]
					if !matcherEntryNames.MatchString(name) {
						return fmt.Errorf("only alphanumeric characters, hyphens, "+
							"and underscores up to 20 characters are allowed for name "+
							"of entries in rule: received '%s'", name)
					}
					sourceAddressSets = strings.Split(rawRule[j*5+3], ":")
					if length := len(sourceAddressSets); length > 10 {
						return fmt.Errorf("length of source-address-sets of entries in rule "+
							"must be less than or equal to 10: received %d", length)
					}
					destinationAddressSets = strings.Split(rawRule[j*5+4], ":")
					if length := len(destinationAddressSets); length > 10 {
						return fmt.Errorf("length of destination-address-sets of entries in "+
							"rule must be less than or equal to 10: received %d", length)
					}
					application = rawRule[j*5+5]
					action = rawRule[j*5+6]
					if !utils.StringInSlice(action, validActions) {
						return fmt.Errorf("action of entries in rule must be either one of "+
							"following: %s: received '%s'", validActions, action)
					}

					rule.Entries = append(rule.Entries, firewalls.Entry{
						Name: name,
						Match: firewalls.Match{
							SourceAddressSets:      sourceAddressSets,
							DestinationAddressSets: destinationAddressSets,
							Application:            application,
						},
						Action: action,
					})
				}
				rules = append(rules, rule)
			}

			if length := len(rawCustomApplications); length > 20 {
				return fmt.Errorf("length of custom-applications must be less than "+
					"or equal to 20: received %d", length)
			}
			for _, v := range rawCustomApplications {
				var customApplication firewalls.CustomApplication
				rawCustomApplication := strings.Split(v, ",")
				if len(rawCustomApplication) != 3 {
					return fmt.Errorf(customApplicationsFormatError)
				}
				customApplication.Name = rawCustomApplication[0]
				if !matcherKeywords.MatchString(customApplication.Name) ||
					strings.Contains(customApplication.Name, "any") ||
					strings.Contains(customApplication.Name, "pre-defined") {
					return fmt.Errorf("only alphanumeric characters, hyphens, underscores, "+
						"ampersands, and parenthesis up to 64 characters are allowed for name "+
						"of custom-applications: received %s", customApplication.Name)
				}
				customApplication.Protocol = rawCustomApplication[1]
				if !utils.StringInSlice(customApplication.Protocol, validProtocols) {
					return fmt.Errorf("protocol for custom-applications must be one of following:"+
						"%s: received '%s'", validProtocols, customApplication.Protocol)
				}
				customApplication.DestinationPort = rawCustomApplication[2]
				if !matcherPorts.MatchString(customApplication.DestinationPort) {
					return fmt.Errorf("destination-port for custom-applications must be either a port number "+
						"or port number range (xxx-yyy): received '%s'", customApplication.DestinationPort)
				}
				customApplications = append(customApplications, customApplication)
			}

			if length := len(rawApplicationSets); length > 2 {
				return fmt.Errorf("length of application-sets must be up to 2: received %d", length)
			}
			for _, v := range rawApplicationSets {
				var applicationSet firewalls.ApplicationSet
				rawApplicationSet := strings.Split(v, ",")
				if len(rawApplicationSet) != 2 {
					return fmt.Errorf(applicationSetsFormatError)
				}
				applicationSet.Name = rawApplicationSet[0]
				if !matcherKeywords.MatchString(applicationSet.Name) ||
					strings.Contains(applicationSet.Name, "any") ||
					strings.Contains(applicationSet.Name, "pre-defined") {
					return fmt.Errorf("only alphanumeric characters, hyphens, underscores, "+
						"ampersands, and parenthesis up to 64 characters are allowed for name "+
						"of application-sets: received '%s'", applicationSet.Name)
				}
				applications := strings.Split(rawApplicationSet[1], ":")
				if length := len(applications); length > 10 {
					return fmt.Errorf("length of applications in application-sets "+
						"must be less than or equal to 10: received %d", length)
				}
				for _, v := range applications {
					applicationSet.Applications = append(applicationSet.Applications, v)
				}
				applicationSets = append(applicationSets, applicationSet)
			}

			if length := len(rawRoutingGroupSettings); length > 40 {
				return fmt.Errorf("length of routing-group-settings must be up to 40: received %d", length)
			}
			for _, v := range rawRoutingGroupSettings {
				var routingGroupSetting firewalls.RoutingGroupSetting
				rawRoutingGroupSetting := strings.Split(v, ",")
				if length := len(rawRoutingGroupSetting); length < 3 || (length-1)%2 != 0 {
					return fmt.Errorf(routingGroupSettingsFormatError)
				}
				routingGroupSetting.GroupName = rawRoutingGroupSetting[0]
				if !utils.StringInSlice(routingGroupSetting.GroupName, validGroupNames) {
					return fmt.Errorf("name in routing-group-settings must be either one "+
						"of following: %s: received '%s'", validGroupNames, routingGroupSetting.GroupName)
				}
				lengthAddressSets := (len(rawRoutingGroupSetting) - 1) / 2
				if lengthAddressSets > 5 {
					return fmt.Errorf("length of address-sets in routing-group-settings must be up to 5: received %d", lengthAddressSets)
				}
				for i := 0; i < lengthAddressSets; i++ {
					var (
						name      string
						addresses []string
					)
					name = rawRoutingGroupSetting[i*2+1]
					if !matcherKeywords.MatchString(name) {
						return fmt.Errorf("only alphanumeric characters, hyphens, underscores, ampersands, "+
							"and parenthesis up to 64 characters are allowed for name of address-sets in "+
							"routing-group-settings: received '%s'", name)
					}
					addresses = strings.Split(rawRoutingGroupSetting[i*2+2], ":")
					if len(addresses) == 1 && addresses[0] == "" {
						addresses = []string{}
					}
					if length := len(addresses); length > 10 {
						return fmt.Errorf("addresses of address-sets in routing-group-settings must be up to 10: received %d", length)
					}
					for _, v := range addresses {
						if _, _, err := net.ParseCIDR(v); err != nil {
							return fmt.Errorf("addresses of address-sets in routing-group-settings, e.g. 192.168.1.0/24: received %s", v)
						}
					}
					routingGroupSetting.AddressSets = append(routingGroupSetting.AddressSets, firewalls.AddressSet{
						Name:      name,
						Addresses: addresses,
					})
				}
				routingGroupSettings = append(routingGroupSettings, routingGroupSetting)
			}

			var customApplicationNames []string
			var applicationSetsNames []string

			for _, v := range customApplications {
				customApplicationNames = append(customApplicationNames, v.Name)
			}
			for _, v := range applicationSets {
				applicationSetsNames = append(applicationSetsNames, v.Name)
			}
			if matched, ok := utils.FindDuplicatedStr(append(customApplicationNames, applicationSetsNames...)); ok {
				return fmt.Errorf("all of custom-application-names and application-set-names must be unique among them; "+
					"duplicated name is: %s", matched)
			}
			allAppNames := append(customApplicationNames, applicationSetsNames...)
			for _, v := range rules {
				for _, w := range v.Entries {
					if !utils.StringInSlice(w.Match.Application, append(allAppNames, "pre-defined", "any")) {
						return fmt.Errorf("application of entries in rule must be either one "+
							"of following: custom-application, application, 'pre-defined-application', "+
							"or 'any': received '%s'", w.Match.Application)
					}
				}
			}

			updateOpts := firewalls.UpdateOpts{
				Rules:                rules,
				CustomApplications:   customApplications,
				ApplicationSets:      applicationSets,
				RoutingGroupSettings: routingGroupSettings,
			}

			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			firewall, err := firewalls.Update(client, routerID, firewallID, updateOpts).Extract()
			if err != nil {
				return fmt.Errorf("calling Update firewalls API: %w", err)
			}

			t := utils.NewTabby(out)
			t.AddHeader("id", "tenantId", "redundant", "isActivated", "operationStatus", "operationId")
			t.AddLine(firewall.ID, firewall.TenantID, firewall.Redundant, firewall.IsActivated,
				firewall.OperationStatus, firewall.OperationID)
			t.Print()
			return nil
		},
	}
	cmd.Flags().StringVar(&routerID, "router", "", "(required) Router ID")
	cmd.Flags().StringArrayVar(
		&rawRules,
		"rules",
		[]string{},
		"(Required) List of Router Rules(<from>,<to>,<entries.name>,"+
			"<entries.match.sourceAddressSet>:<entries.match.sourceAddressSet>...,"+
			"<entries.match.destinationAddressSets>:<entries.match.destinationAddressSets>...,"+
			"<entries.match.application>,<entries.action>), "+
			"e.g. group_1,group_2,rule-01,group1_addset_1:group1_addset_2,"+
			"group2_addset_1:group2_addset_2,app_set_1,permit",
	)
	cmd.Flags().StringArrayVar(
		&rawCustomApplications,
		"custom-applications",
		[]string{},
		"(Required) Custom Applications(<name>,<protocol>,<destinationPort>), "+
			"e.g. google-drive-web,tcp,443",
	)
	cmd.Flags().StringArrayVar(
		&rawApplicationSets,
		"application-sets",
		[]string{},
		"(Required) Application Sets(<name>,<applications>:<applications>...), "+
			"e.g. app_set_1,google-drive-web:pre-defined-ftp",
	)
	cmd.Flags().StringArrayVar(
		&rawRoutingGroupSettings,
		"routing-group-settings",
		[]string{},
		"(Required) Routing Group Settings(<groupName>,<addressSets.name>,"+
			"<addressSets.addresses>:<addressSets.addresses>...), "+
			"e.g. group_1,group1_addset_1,172.18.1.0/24:192.168.1.0/24",
	)

	cmd.MarkFlagRequired("router")
	cmd.MarkFlagRequired("rules")
	cmd.MarkFlagRequired("custom-applications")
	cmd.MarkFlagRequired("application-sets")
	cmd.MarkFlagRequired("routing-group-settings")

	return cmd
}

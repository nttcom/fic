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

package testdata

const ShowResponse = `
{
	"firewall": {
	  "id": "F040123456789",
	  "tenantId": "080f290761484afabbec22938adc6a2e",
	  "redundant": false,
	  "isActivated": true,
	  "operationStatus": "Processing",
	  "customApplications": [
		{
		  "name": "google-drive-web",
		  "protocol": "tcp",
		  "destinationPort": "443"
		}
	  ],
	  "applicationSets": [
		{
		  "name": "app_set_1",
		  "applications": [
			"google-drive-web",
			"pre-defined-ftp"
		  ]
		}
	  ],
	  "routingGroupSettings": [
		{
		  "groupName": "group_1",
		  "addressSets": [
			{
			  "name": "group1_addset_1",
			  "addresses": [
				"172.18.1.0/24"
			  ]
			}
		  ]
		},
		{
		  "groupName": "group_2",
		  "addressSets": [
			{
			  "name": "group2_addset_1",
			  "addresses": [
				"192.168.1.0/24"
			  ]
			}
		  ]
		},
		{
		  "groupName": "group_3",
		  "addressSets": []
		},
		{
		  "groupName": "group_4",
		  "addressSets": []
		}
	  ],
	  "rules": [
		{
		  "from": "group_1",
		  "to": "group_2",
		  "entries": [
			{
			  "name": "rule-01",
			  "match": {
				"sourceAddressSets": [
				  "group1_addset_1"
				],
				"destinationAddressSets": [
				  "group2_addset_1"
				],
				"application": "app_set_1"
			  },
			  "action": "permit"
			}
		  ]
		}
	  ],
	  "userIpAddresses": [
		"192.168.0.0/30",
		"192.168.0.4/30",
		"192.168.0.8/30",
		"192.168.0.12/30"
	  ]
	}
  } 
`

const ShowOutput = `
id             tenantId                          redundant  isActivated  operationStatus
--             --------                          ---------  -----------  ---------------
F040123456789  080f290761484afabbec22938adc6a2e  false      true         Processing
`

const CustomApplicationsOutput = `
name              protocol  destinationPort
----              --------  ---------------
google-drive-web  tcp       443
`

const ApplicationSetsOutput = `
name       applications
----       ------------
app_set_1  [google-drive-web pre-defined-ftp]
`

const RoutingGroupSettingsOutput = `
groupName  addressSets
---------  -----------
group_1    [{group1_addset_1 [172.18.1.0/24]}]
group_2    [{group2_addset_1 [192.168.1.0/24]}]
group_3    []
group_4    []
`

const RulesOutput = `
from     to       entries
----     --       -------
group_1  group_2  [{rule-01 {[group1_addset_1] [group2_addset_1] app_set_1} permit}]
`

const UserIPAddressesOutput = `
userIPAddresses
---------------
192.168.0.0/30
192.168.0.4/30
192.168.0.8/30
192.168.0.12/30
`

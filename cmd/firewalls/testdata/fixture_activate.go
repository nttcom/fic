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

const ActivateRequest = `
{
    "firewall": {
        "userIpAddresses": [
            "192.168.0.0/30",
            "192.168.0.4/30",
            "192.168.0.8/30",
            "192.168.0.12/30"
        ]
    }
}
`

const ActivateResponse = `
{
    "firewall": {
        "id": "F040123456789",
        "tenantId": "080f290761484afabbec22938adc6a2e",
        "redundant": false,
        "isActivated": true,
        "operationStatus": "Processing",
        "customApplications": [],
        "applicationSets": [],
        "routingGroupSettings": [],
        "rules": [],
        "userIpAddresses": [
            "192.168.0.0/30",
            "192.168.0.4/30",
            "192.168.0.8/30",
            "192.168.0.12/30"
        ],
        "operationId": "4c7b0bfc17e84e9eae7b4d779f30dea1"
    }
}
`

const ActivateOutput = `
id             tenantId                          redundant  isActivated  operationStatus  operationId
--             --------                          ---------  -----------  ---------------  -----------
F040123456789  080f290761484afabbec22938adc6a2e  false      true         Processing       4c7b0bfc17e84e9eae7b4d779f30dea1
`

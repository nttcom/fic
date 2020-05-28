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

const DeactivateRequest = `{}`

const DeactivateResponse = `
{
    "firewall": {
        "id": "F040123456789",
        "tenantId": "080f290761484afabbec22938adc6a2e",
        "redundant": false,
        "isActivated": false,
        "operationStatus": "Processing",
        "customApplications": [],
        "applicationSets": [],
        "routingGroupSettings": [],
        "rules": [],
        "userIpAddresses": [],
        "operationId": "d40999c03a7642c6b86158889c7bebc9"
    }
}
`

const DeactivateOutput = `
id             tenantId                          redundant  isActivated  operationStatus  operationId
--             --------                          ---------  -----------  ---------------  -----------
F040123456789  080f290761484afabbec22938adc6a2e  false      false        Processing       d40999c03a7642c6b86158889c7bebc9
`

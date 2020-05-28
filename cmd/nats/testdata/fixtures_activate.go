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
    "nat": {
        "globalIpAddressSets": [
            {
                "name": "src-set-00000000000000000000000000000000000000000000000000000001",
                "numOfAddresses": 5,
                "type": "sourceNapt"
            },
            {
                "name": "dst-set-01",
                "numOfAddresses": 1,
                "type": "destinationNat"
            }
        ],
        "userIpAddresses": [
            "192.168.0.0/30",
            "192.168.4.0/30",
            "192.168.8.0/30",
            "192.168.12.0/30",
            "192.168.16.0/30",
            "192.168.20.0/30",
            "192.168.24.0/30",
            "192.168.28.0/30"
        ]
    }
}
`

const ActivateResponse = `
{
    "nat": {
        "destinationNatRules": [],
        "globalIpAddressSets": [
            {
                "addresses": [
                    "100.131.66.12",
                    "100.131.66.13",
                    "100.131.66.14",
                    "100.131.66.15",
                    "100.131.66.16"
                ],
                "id": "f9b94a8367f54b6384f02ab8c52d8430",
                "name": "src-set-01",
                "natComponentId": "F052000000335",
                "numOfAddresses": 5,
                "operationStatus": "Processing",
                "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
                "type": "sourceNapt"
            },
            {
                "addresses": [
                    "100.131.65.3"
                ],
                "id": "e2124ae289404249a6ca14dfea555373",
                "name": "dst-set-01",
                "natComponentId": "F052000000335",
                "numOfAddresses": 1,
                "operationStatus": "Processing",
                "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
                "type": "destinationNat"
            }
        ],
        "id": "F052000000335",
        "isActivated": true,
        "operationId": "4e266805a4564bd7a40e7d1dd7474b84",
        "operationStatus": "Processing",
        "redundant": true,
        "sourceNaptRules": [],
        "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
        "userIpAddresses": [
            "192.168.0.0/30",
            "192.168.4.0/30",
            "192.168.8.0/30",
            "192.168.12.0/30",
            "192.168.16.0/30",
            "192.168.20.0/30",
            "192.168.24.0/30",
            "192.168.28.0/30"
        ]
    }
}
`

const ActivateRequestWithGlobalIPAddressSetsOnly = `
{
    "nat": {
        "globalIpAddressSets": [
            {
                "name": "src-set-01",
                "numOfAddresses": 5,
                "type": "sourceNapt"
            },
            {
                "name": "dst-set-01",
                "numOfAddresses": 1,
                "type": "destinationNat"
            }
        ]
    }
}
`

const ActivateOutput = `
id             tenantID                          redundant  isActivated  operationID                       operationStatus
--             --------                          ---------  -----------  -----------                       ---------------
F052000000335  75ece6039bfd4d2f88c8b8d790482886  true       true         4e266805a4564bd7a40e7d1dd7474b84  Processing
`

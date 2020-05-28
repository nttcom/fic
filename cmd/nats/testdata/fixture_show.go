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
    "nat": {
        "destinationNatRules": [
            {
                "entries": [
                    {
                        "match": {
                            "destinationAddress": "dst-set-01"
                        },
                        "then": "192.168.0.1/32"
                    }
                ],
                "from": "group_1",
                "to": "group_2"
            }
        ],
        "id": "F052000000335",
        "isActivated": true,
        "operationStatus": "Completed",
        "redundant": true,
        "sourceNaptRules": [
            {
                "entries": [
                    {
                        "then": [
                            "src-set-01"
                        ]
                    }
                ],
                "from": [
                    "group_1"
                ],
                "to": "group_2"
            }
        ],
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
        ],
        "globalIpAddressSets": [
            {
                "addresses": [
                    "100.131.66.12",
                    "100.131.66.13",
                    "100.131.66.14",
                    "100.131.66.15",
                    "100.131.66.16"
                ],
                "id": "2b065939916f4ccfa1be3d532900c4de",
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
                "id": "1b5fe60eca174db6a4042da84c934d22",
                "name": "dst-set-01",
                "natComponentId": "F052000000335",
                "numOfAddresses": 1,
                "operationStatus": "Processing",
                "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
                "type": "destinationNat"
            }
        ]
    }
}
`

const ShowOutput = `
id             tenantID                          redundant  isActivated  operationStatus
--             --------                          ---------  -----------  ---------------
F052000000335  75ece6039bfd4d2f88c8b8d790482886  true       true         Completed
`

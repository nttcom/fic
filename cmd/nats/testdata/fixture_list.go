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

const ListResponse = `
{
    "nats": [
        {
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
            ]
        }
    ]
}
`

const ListOutput = `
id             tenantID                          redundant  isActivated  operationStatus
--             --------                          ---------  -----------  ---------------
F052000000335  75ece6039bfd4d2f88c8b8d790482886  true       true         Completed
`

// "isActivated": "true" should be "isActivated": true
const InvalidListResponse = `
{
    "nats": [
        {
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
            "isActivated": "true",
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
            ]
        }
    ]
}
`

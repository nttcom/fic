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

const UpdateRequest = `
{
    "nat": {
        "sourceNaptRules": [
            {
                "from": [
                    "group_1"
                ],
                "to": "group_2",
                "entries": [
                    {
                        "then": [
                            "src-set-01",
                            "src-set-02",
                            "src-set-03",
                            "src-set-04"
                        ]
                    }
                ]
            },
            {
                "from": [
                    "group_2"
                ],
                "to": "group_1",
                "entries": [
                    {
                        "then": [
                            "src-set-05",
                            "src-set-06",
                            "src-set-07",
                            "src-set-08"
                        ]
                    }
                ]
            }
        ],
        "destinationNatRules": [
            {
                "from": "group_1",
                "to": "group_2",
                "entries": [
                    {
                        "match": {
                            "destinationAddress": "dst-set-01"
                        },
                        "then": "192.168.0.1/32"
                    },
                    {
                        "match": {
                            "destinationAddress": "dst-set-02"
                        },
                        "then": "192.168.0.2/32"
                    }
                ]
            },
            {
                "from": "group_2",
                "to": "group_1",
                "entries": [
                    {
                        "match": {
                            "destinationAddress": "dst-set-03"
                        },
                        "then": "192.168.0.3/32"
                    },
                    {
                        "match": {
                            "destinationAddress": "dst-set-04"
                        },
                        "then": "192.168.0.4/32"
                    }
                ]
            }
        ]
    }
}
`

const UpdateResponse = `
{
    "nat": {
        "id": "F052000000335",
        "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
        "redundant": false,
        "isActivated": true,
        "operationStatus": "Processing",
        "sourceNaptRules": [
            {
                "from": [
                    "group_1"
                ],
                "to": "group_2",
                "entries": [
                    {
                        "then": [
                            "src-set-01",
                            "src-set-02",
                            "src-set-03",
                            "src-set-04"
                        ]
                    }
                ]
            },
            {
                "from": [
                    "group_2"
                ],
                "to": "group_1",
                "entries": [
                    {
                        "then": [
                            "src-set-05",
                            "src-set-06",
                            "src-set-07",
                            "src-set-08"
                        ]
                    }
                ]
            }
        ],
        "destinationNatRules": [
            {
                "from": "group_1",
                "to": "group_2",
                "entries": [
                    {
                        "match": {
                            "destinationAddress": "dst-set-01"
                        },
                        "then": "192.168.0.1/32"
                    },
                    {
                        "match": {
                            "destinationAddress": "dst-set-02"
                        },
                        "then": "192.168.0.2/32"
                    }
                ]
            },
            {
                "from": "group_2",
                "to": "group_1",
                "entries": [
                    {
                        "match": {
                            "destinationAddress": "dst-set-03"
                        },
                        "then": "192.168.0.3/32"
                    },
                    {
                        "match": {
                            "destinationAddress": "dst-set-04"
                        },
                        "then": "192.168.0.4/32"
                    }
                ]
            }
        ],
        "userIpAddresses": [
            "192.168.0.0/30",
            "192.168.0.4/30",
            "192.168.0.8/30",
            "192.168.0.12/30"
        ],
        "operationId": "f6a4b449b40d4660aa3aae7985dfe2a6"
    }
}
`

const UpdateOutput = `
id             tenantID                          redundant  isActivated  operationID                       operationStatus
--             --------                          ---------  -----------  -----------                       ---------------
F052000000335  75ece6039bfd4d2f88c8b8d790482886  false      true         f6a4b449b40d4660aa3aae7985dfe2a6  Processing
`

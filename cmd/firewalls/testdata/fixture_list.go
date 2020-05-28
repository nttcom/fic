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
    "firewalls": [
        {
            "id": "%s",
            "tenantId": "080f290761484afabbec22938adc6a2e",
            "redundant": false,
            "isActivated": true,
            "operationStatus": "Completed",
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
    ]
}
`

const ListOutput = `
id            tenantId                          redundant  isActivated  operationStatus
--            --------                          ---------  -----------  ---------------
%!s(MISSING)  080f290761484afabbec22938adc6a2e  true       Completed
`

// "redundant": "false" should be "redundant": false
const InvalidListResponse = `
{
    "firewalls": [
        {
            "id": "%s",
            "tenantId": "080f290761484afabbec22938adc6a2e",
            "redundant": "false",
            "isActivated": true,
            "operationStatus": "Completed",
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
    ]
}
`

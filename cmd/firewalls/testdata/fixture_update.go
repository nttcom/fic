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
  "firewall": {
    "applicationSets": [
      {
        "applications": [
          "google-drive-web",
          "pre-defined-ftp"
        ],
        "name": "app_set_1"
      }
    ],
    "customApplications": [
      {
        "destinationPort": "443",
        "name": "google-drive-web",
        "protocol": "tcp"
      }
    ],
    "routingGroupSettings": [
      {
        "addressSets": [
          {
            "addresses": [
              "172.18.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_1"
      },
      {
        "addressSets": [
          {
            "addresses": [
              "192.168.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_2"
      },
      {
        "addressSets": [
          {
            "addresses": [],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_3"
      }
    ],
    "rules": [
      {
        "entries": [
          {
            "action": "permit",
            "match": {
              "application": "app_set_1",
              "destinationAddressSets": [
                "group2_addset_1",
                "group2_addset_2"
              ],
              "sourceAddressSets": [
                "group1_addset_1",
                "group1_addset_2"
              ]
            },
            "name": "rule-01"
          }
        ],
        "from": "group_1",
        "to": "group_2"
      }
    ]
  }
}
`
const UpdateRequestWithPortRange = `
{
  "firewall": {
    "applicationSets": [
      {
        "applications": [
          "google-drive-web",
          "pre-defined-ftp"
        ],
        "name": "app_set_1"
      }
    ],
    "customApplications": [
      {
        "destinationPort": "8080-8089",
        "name": "google-drive-web",
        "protocol": "tcp"
      }
    ],
    "routingGroupSettings": [
      {
        "addressSets": [
          {
            "addresses": [
              "172.18.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_1"
      },
      {
        "addressSets": [
          {
            "addresses": [
              "192.168.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_2"
      },
      {
        "addressSets": [
          {
            "addresses": [],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_3"
      }
    ],
    "rules": [
      {
        "entries": [
          {
            "action": "permit",
            "match": {
              "application": "app_set_1",
              "destinationAddressSets": [
                "group2_addset_1",
                "group2_addset_2"
              ],
              "sourceAddressSets": [
                "group1_addset_1",
                "group1_addset_2"
              ]
            },
            "name": "rule-01"
          }
        ],
        "from": "group_1",
        "to": "group_2"
      }
    ]
  }
}
`
const UpdateRequestMultipleRules = `
{
  "firewall": {
    "applicationSets": [
      {
        "applications": [
          "google-drive-web",
          "pre-defined-ftp"
        ],
        "name": "app_set_1"
      }
    ],
    "customApplications": [
      {
        "destinationPort": "443",
        "name": "google-drive-web",
        "protocol": "tcp"
      }
    ],
    "routingGroupSettings": [
      {
        "addressSets": [
          {
            "addresses": [
              "172.18.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_1"
      },
      {
        "addressSets": [
          {
            "addresses": [
              "192.168.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_2"
      },
      {
        "addressSets": [
          {
            "addresses": [],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_3"
      }
    ],
    "rules": [
      {
        "entries": [
          {
            "action": "permit",
            "match": {
              "application": "app_set_1",
              "destinationAddressSets": [
                "group2_addset_1",
                "group2_addset_2"
              ],
              "sourceAddressSets": [
                "group1_addset_1",
                "group1_addset_2"
              ]
            },
            "name": "rule-01"
          }
        ],
        "from": "group_1",
        "to": "group_2"
      },
      {
        "entries": [
          {
            "action": "permit",
            "match": {
              "application": "google-drive-web",
              "destinationAddressSets": [
                "group2_addset_1",
                "group2_addset_2"
              ],
              "sourceAddressSets": [
                "group1_addset_1",
                "group1_addset_2"
              ]
            },
            "name": "rule-01"
          }
        ],
        "from": "group_1",
        "to": "group_2"
      }
    ]
  }
}
`

const UpdateRequestMultipleCustomApplications = `
{
  "firewall": {
    "applicationSets": [
      {
        "applications": [
          "google-drive-web",
          "pre-defined-ftp"
        ],
        "name": "app_set_1"
      }
    ],
    "customApplications": [
      {
        "destinationPort": "443",
        "name": "google-drive-web",
        "protocol": "tcp"
      },
      {
        "destinationPort": "443",
        "name": "my-web-app",
        "protocol": "tcp"
      }
    ],
    "routingGroupSettings": [
      {
        "addressSets": [
          {
            "addresses": [
              "172.18.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_1"
      },
      {
        "addressSets": [
          {
            "addresses": [
              "192.168.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_2"
      },
      {
        "addressSets": [
          {
            "addresses": [],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_3"
      }
    ],
    "rules": [
      {
        "entries": [
          {
            "action": "permit",
            "match": {
              "application": "app_set_1",
              "destinationAddressSets": [
                "group2_addset_1",
                "group2_addset_2"
              ],
              "sourceAddressSets": [
                "group1_addset_1",
                "group1_addset_2"
              ]
            },
            "name": "rule-01"
          }
        ],
        "from": "group_1",
        "to": "group_2"
      }
    ]
  }
}
`

const UpdateRequestMultipleApplicationSets = `
{
  "firewall": {
    "applicationSets": [
      {
        "applications": [
          "google-drive-web",
          "pre-defined-ftp"
        ],
        "name": "app_set_1"
      },
      {
        "applications": [
          "my-web-app",
          "azure-virtual-instance"
        ],
        "name": "app_set_2"
      }
    ],
    "customApplications": [
      {
        "destinationPort": "443",
        "name": "google-drive-web",
        "protocol": "tcp"
      }
    ],
    "routingGroupSettings": [
      {
        "addressSets": [
          {
            "addresses": [
              "172.18.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_1"
      },
      {
        "addressSets": [
          {
            "addresses": [
              "192.168.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_2"
      },
      {
        "addressSets": [
          {
            "addresses": [],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_3"
      }
    ],
    "rules": [
      {
        "entries": [
          {
            "action": "permit",
            "match": {
              "application": "app_set_1",
              "destinationAddressSets": [
                "group2_addset_1",
                "group2_addset_2"
              ],
              "sourceAddressSets": [
                "group1_addset_1",
                "group1_addset_2"
              ]
            },
            "name": "rule-01"
          }
        ],
        "from": "group_1",
        "to": "group_2"
      }
    ]
  }
}
`

const UpdateResponse = `
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
        ],
        "operationId": "2fb5baa5a5834e0b952d3f0d93c3e64a"
    }
}
`
const UpdateResponseWithPortRange = `
{
  "firewall": {
    "id": "F040123456789",
    "tenantId": "080f290761484afabbec22938adc6a2e",
    "redundant": false,
    "isActivated": true,
    "operationStatus": "Processing",
    "applicationSets": [
      {
        "applications": [
          "google-drive-web",
          "pre-defined-ftp"
        ],
        "name": "app_set_1"
      }
    ],
    "customApplications": [
      {
        "destinationPort": "8080-8089",
        "name": "google-drive-web",
        "protocol": "tcp"
      }
    ],
    "routingGroupSettings": [
      {
        "addressSets": [
          {
            "addresses": [
              "172.18.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_1"
      },
      {
        "addressSets": [
          {
            "addresses": [
              "192.168.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_2"
      },
      {
        "addressSets": [
          {
            "addresses": [],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_3"
      }
    ],
    "rules": [
      {
        "entries": [
          {
            "action": "permit",
            "match": {
              "application": "app_set_1",
              "destinationAddressSets": [
                "group2_addset_1",
                "group2_addset_2"
              ],
              "sourceAddressSets": [
                "group1_addset_1",
                "group1_addset_2"
              ]
            },
            "name": "rule-01"
          }
        ],
        "from": "group_1",
        "to": "group_2"
      }
    ],
    "operationId": "2fb5baa5a5834e0b952d3f0d93c3e64a"
  }
}
`
const UpdateResponseMultipleRules = `
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
          },
          {
            "entries": [
              {
                "action": "permit",
                "match": {
                  "application": "google-drive-web",
                  "destinationAddressSets": [
                    "group2_addset_1",
                    "group2_addset_2"
                  ],
                  "sourceAddressSets": [
                    "group1_addset_1",
                    "group1_addset_2"
                  ]
                },
                "name": "rule-01"
              }
            ],
            "from": "group_1",
            "to": "group_2"
          }          
      ],
      "userIpAddresses": [
          "192.168.0.0/30",
          "192.168.0.4/30",
          "192.168.0.8/30",
          "192.168.0.12/30"
      ],
      "operationId": "2fb5baa5a5834e0b952d3f0d93c3e64a"
  }
}
`

const UpdateResponseMultipleCustomApplications = `
{
  "firewall": {
    "id": "F040123456789",
    "tenantId": "080f290761484afabbec22938adc6a2e",
    "redundant": false,
    "isActivated": true,
    "operationStatus": "Processing",
    "applicationSets": [
      {
        "applications": [
          "google-drive-web",
          "pre-defined-ftp"
        ],
        "name": "app_set_1"
      }
    ],
    "customApplications": [
      {
        "destinationPort": "443",
        "name": "google-drive-web",
        "protocol": "tcp"
      },
      {
        "destinationPort": "443",
        "name": "my-web-app",
        "protocol": "tcp"
      }
    ],
    "routingGroupSettings": [
      {
        "addressSets": [
          {
            "addresses": [
              "172.18.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_1"
      },
      {
        "addressSets": [
          {
            "addresses": [
              "192.168.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_2"
      },
      {
        "addressSets": [
          {
            "addresses": [],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_3"
      }
    ],
    "rules": [
      {
        "entries": [
          {
            "action": "permit",
            "match": {
              "application": "app_set_1",
              "destinationAddressSets": [
                "group2_addset_1",
                "group2_addset_2"
              ],
              "sourceAddressSets": [
                "group1_addset_1",
                "group1_addset_2"
              ]
            },
            "name": "rule-01"
          }
        ],
        "from": "group_1",
        "to": "group_2"
      }
    ],
    "operationId": "2fb5baa5a5834e0b952d3f0d93c3e64a"
  }
}
`

const UpdateResponseMultipleApplicationSets = `
{
  "firewall": {
    "id": "F040123456789",
    "tenantId": "080f290761484afabbec22938adc6a2e",
    "redundant": false,
    "isActivated": true,
    "operationStatus": "Processing",
    "applicationSets": [
      {
        "applications": [
          "google-drive-web",
          "pre-defined-ftp"
        ],
        "name": "app_set_1"
      },
      {
        "applications": [
          "my-web-app",
          "azure-virtual-instance"
        ],
        "name": "app_set_2"
      }
    ],
    "customApplications": [
      {
        "destinationPort": "443",
        "name": "google-drive-web",
        "protocol": "tcp"
      }
    ],
    "routingGroupSettings": [
      {
        "addressSets": [
          {
            "addresses": [
              "172.18.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_1"
      },
      {
        "addressSets": [
          {
            "addresses": [
              "192.168.1.0/24"
            ],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_2"
      },
      {
        "addressSets": [
          {
            "addresses": [],
            "name": "group1_addset_1"
          }
        ],
        "groupName": "group_3"
      }
    ],
    "rules": [
      {
        "entries": [
          {
            "action": "permit",
            "match": {
              "application": "app_set_1",
              "destinationAddressSets": [
                "group2_addset_1",
                "group2_addset_2"
              ],
              "sourceAddressSets": [
                "group1_addset_1",
                "group1_addset_2"
              ]
            },
            "name": "rule-01"
          }
        ],
        "from": "group_1",
        "to": "group_2"
      }
    ],
    "operationId": "2fb5baa5a5834e0b952d3f0d93c3e64a"
  }
}
`

const UpdateOutput = `
id             tenantId                          redundant  isActivated  operationStatus  operationId
--             --------                          ---------  -----------  ---------------  -----------
F040123456789  080f290761484afabbec22938adc6a2e  false      true         Processing       2fb5baa5a5834e0b952d3f0d93c3e64a
`

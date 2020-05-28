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
    "connection": {
        "source": {
            "primary": {
                "asPathPrepend": {
                    "in": null,
                    "out": 2
                },
                "med": {
                    "out": 10
                }
            },
            "routeFilter": {
                "in": "noRoute",
                "out": "fullRoute"
            }
        }
    }
}
`

const UpdateResponse = `
{
    "connection": {
        "area": "JPWEST",
        "bandwidth": "10M",
        "destination": {
            "primary": {
                "asn": "65000",
                "ipAddress": "10.0.0.2/30",
                "portId": "F012000000220",
                "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
                "vlan": 161
            },
            "publicServiceKey": null
        },
        "id": "F032000001027",
        "interTenant": false,
        "name": "testConnection",
        "operationId": "4f5d6b7ba4a5433398e6712a3ea0c126",
        "operationStatus": "Processing",
        "redundant": false,
        "source": {
            "groupName": "group_1",
            "primary": {
                "asPathPrepend": {
                    "in": null,
                    "out": 2
                },
                "ipAddress": "10.0.0.1/30",
                "med": {
                    "out": 10
                }
            },
            "routeFilter": {
                "in": "noRoute",
                "out": "fullRoute"
            },
            "routerId": "F022000000168",
            "tenantId": "75ece6039bfd4d2f88c8b8d790482886"
        },
        "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
        "toPublicService": false
    }
}
`

const UpdateOutput = `
id             name            redundant  tenantId                          area    operationStatus                   bandwidth   operationId
--             ----            ---------  --------                          ----    ---------------                   ---------   -----------
F032000001027  testConnection  false      75ece6039bfd4d2f88c8b8d790482886  JPWEST  4f5d6b7ba4a5433398e6712a3ea0c126  Processing  10M
`

const UpdatePairedRequest = `
{
    "connection": {
        "source": {
            "primary": {
                "asPathPrepend": {
                    "in": null,
                    "out": 2
                },
                "med": {
                    "out": 10
                }
            },
            "routeFilter": {
                "in": "noRoute",
                "out": "fullRoute"
            },
            "secondary": {
                "asPathPrepend": {
                    "in": null,
                    "out": 2
                }
            }
        }
    }
}
`

const UpdatePairedResponse = `
{
    "connection": {
        "area": "JPWEST",
        "bandwidth": "10M",
        "destination": {
            "primary": {
                "asn": "65000",
                "ipAddress": "10.0.0.2/30",
                "portId": "F012000000220",
                "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
                "vlan": 161
            },
            "publicServiceKey": null,
            "secondary": {
                "asn": "65000",
                "ipAddress": "10.0.0.6/30",
                "portId": "F012000000219",
                "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
                "vlan": 190
            }
        },
        "id": "F032000001012",
        "interTenant": false,
        "name": "testConnection",
        "operationId": "9ccffe095d4b4913a029dce7d7799ca3",
        "operationStatus": "Processing",
        "redundant": true,
        "source": {
            "groupName": "group_1",
            "primary": {
                "asPathPrepend": {
                    "in": null,
                    "out": 2
                },
                "ipAddress": "10.0.0.1/30",
                "med": {
                    "out": 10
                }
            },
            "routeFilter": {
                "in": "noRoute",
                "out": "fullRoute"
            },
            "routerId": "F022000000168",
            "secondary": {
                "asPathPrepend": {
                    "in": null,
                    "out": 2
                },
                "ipAddress": "10.0.0.5/30",
                "med": {
                    "out": 20
                }
            },
            "tenantId": "75ece6039bfd4d2f88c8b8d790482886"
        },
        "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
        "toPublicService": false
    }
}
`

const UpdatePairedOutput = `
id             name            redundant  tenantId                          area    operationStatus                   bandwidth   operationId
--             ----            ---------  --------                          ----    ---------------                   ---------   -----------
F032000001012  testConnection  true       75ece6039bfd4d2f88c8b8d790482886  JPWEST  9ccffe095d4b4913a029dce7d7799ca3  Processing  10M
`

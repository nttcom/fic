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

const SourceResponse = `
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
        "operationStatus": "Completed",
        "redundant": false,
        "source": {
            "groupName": "group_1",
            "primary": {
                "asPathPrepend": {
                    "in": 4,
                    "out": 4
                },
                "ipAddress": "10.0.0.1/30",
                "med": {
                    "out": 10
                }
            },
            "routeFilter": {
                "in": "fullRoute",
                "out": "fullRouteWithDefaultRoute"
            },
            "routerId": "F022000000168",
            "tenantId": "75ece6039bfd4d2f88c8b8d790482886"
        },
        "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
        "toPublicService": false
    }
}
`

const SourceOutput = `
tenantId                          routerId       groupName  routeFilter.in  routeFilter.out            primary.ipAddress  primary.asPathPrepend.in  primary.asPathPrepend.out  primary.med.out
--------                          --------       ---------  --------------  ---------------            -----------------  ------------------------  -------------------------  ---------------
75ece6039bfd4d2f88c8b8d790482886  F022000000168  group_1    fullRoute       fullRouteWithDefaultRoute  10.0.0.1/30        4                         4                          10
`

const SourcePairedResponse = `
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
        "operationStatus": "Completed",
        "redundant": true,
        "source": {
            "groupName": "group_1",
            "primary": {
                "asPathPrepend": {
                    "in": 4,
                    "out": 4
                },
                "ipAddress": "10.0.0.1/30",
                "med": {
                    "out": 10
                }
            },
            "routeFilter": {
                "in": "fullRoute",
                "out": "fullRouteWithDefaultRoute"
            },
            "routerId": "F022000000168",
            "secondary": {
                "asPathPrepend": {
                    "in": 2,
                    "out": 1
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

const SourcePairedOutput = `
tenantId                          routerId       groupName  routeFilter.in  routeFilter.out            primary.ipAddress  primary.asPathPrepend.in  primary.asPathPrepend.out  primary.med.out  secondary.ipAddress  secondary.asPathPrepend.in  secondary.asPathPrepend.out  secondary.med.out
--------                          --------       ---------  --------------  ---------------            -----------------  ------------------------  -------------------------  ---------------  -------------------  --------------------------  ---------------------------  -----------------
75ece6039bfd4d2f88c8b8d790482886  F022000000168  group_1    fullRoute       fullRouteWithDefaultRoute  10.0.0.1/30        4                         4                          10               10.0.0.5/30          2                           1                            20
`

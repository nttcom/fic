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
    "globalIpAddressSets": [
        {
            "addresses": [
                "100.131.65.1"
            ],
            "id": "e5c8186e695d4cda9e68b801c3f6f67c",
            "name": "dst-set-01",
            "natComponentId": "F052000000013",
            "numOfAddresses": 1,
            "operationStatus": "Completed",
            "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
            "type": "destinationNat"
        },
        {
            "addresses": [
                "100.131.66.1",
                "100.131.66.2",
                "100.131.66.3",
                "100.131.66.4",
                "100.131.66.5"
            ],
            "id": "3e262888541d404f97546af412ed1dab",
            "name": "src-set-01",
            "natComponentId": "F052000000013",
            "numOfAddresses": 5,
            "operationStatus": "Completed",
            "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
            "type": "sourceNapt"
        }
    ]
}
`

const ListOutput = `
id                                name        type            natComponentID  operationStatus  tenantID                          numOfAddresses
--                                ----        ----            --------------  ---------------  --------                          --------------
e5c8186e695d4cda9e68b801c3f6f67c  dst-set-01  destinationNat  F052000000013   Completed        75ece6039bfd4d2f88c8b8d790482886  1
3e262888541d404f97546af412ed1dab  src-set-01  sourceNapt      F052000000013   Completed        75ece6039bfd4d2f88c8b8d790482886  5
`

// "numOfAddresses": "1" should be "numOfAddresses": 1
const InvalidListResponse = `
{
    "globalIpAddressSets": [
        {
            "addresses": [
                "100.131.65.1"
            ],
            "id": "e5c8186e695d4cda9e68b801c3f6f67c",
            "name": "dst-set-01",
            "natComponentId": "F052000000013",
            "numOfAddresses": "1",
            "operationStatus": "Completed",
            "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
            "type": "destinationNat"
        },
        {
            "addresses": [
                "100.131.66.1",
                "100.131.66.2",
                "100.131.66.3",
                "100.131.66.4",
                "100.131.66.5"
            ],
            "id": "3e262888541d404f97546af412ed1dab",
            "name": "src-set-01",
            "natComponentId": "F052000000013",
            "numOfAddresses": 5,
            "operationStatus": "Completed",
            "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
            "type": "sourceNapt"
        }
    ]
}
`

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

const CreateDstRequest = `
{
    "globalIpAddressSet": {
        "name": "testDst&()-_0123456789012345678901234567890123456789012345678901",
        "numOfAddresses": 1,
        "type": "destinationNat"
    }
}
`

const CreateDstResponse = `
{
    "globalIpAddressSet": {
        "addresses": [
            "100.131.65.2"
        ],
        "id": "bf839f90b4ea4b4192ea6c33ebb85dc5",
        "name": "testDst&()-_0123456789012345678901234567890123456789012345678901",
        "natComponentId": "F052000000013",
        "numOfAddresses": 1,
        "operationId": "78f4791268d7493f8eb48acf28b11d38",
        "operationStatus": "Completed",
        "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
        "type": "destinationNat"
    }
}
`

const CreateDstOutput = `
id                                name                                                              type            natComponentID  operationID                       operationStatus  tenantID                          numOfAddresses
--                                ----                                                              ----            --------------  -----------                       ---------------  --------                          --------------
bf839f90b4ea4b4192ea6c33ebb85dc5  testDst&()-_0123456789012345678901234567890123456789012345678901  destinationNat  F052000000013   78f4791268d7493f8eb48acf28b11d38  Completed        75ece6039bfd4d2f88c8b8d790482886  1
`

const CreateSrcRequest = `
{
    "globalIpAddressSet": {
        "name": "testSrc&()-_0123456789012345678901234567890123456789012345678901",
        "numOfAddresses": 5,
        "type": "sourceNapt"
    }
}
`

const CreateSrcResponse = `
{
    "globalIpAddressSet": {
        "addresses": [
            "100.131.66.6",
            "100.131.66.7",
            "100.131.66.8",
            "100.131.66.9",
            "100.131.66.10"
        ],
        "id": "e9548c69d26a49df96fd7680800fd907",
        "name": "testSrc&()-_0123456789012345678901234567890123456789012345678901",
        "natComponentId": "F052000000013",
        "numOfAddresses": 5,
        "operationId": "afa5c4badc92434e946d01d41c281f8a",
        "operationStatus": "Completed",
        "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
        "type": "sourceNapt"
    }
}
`

const CreateSrcOutput = `
id                                name                                                              type        natComponentID  operationID                       operationStatus  tenantID                          numOfAddresses
--                                ----                                                              ----        --------------  -----------                       ---------------  --------                          --------------
e9548c69d26a49df96fd7680800fd907  testSrc&()-_0123456789012345678901234567890123456789012345678901  sourceNapt  F052000000013   afa5c4badc92434e946d01d41c281f8a  Completed        75ece6039bfd4d2f88c8b8d790482886  5
`

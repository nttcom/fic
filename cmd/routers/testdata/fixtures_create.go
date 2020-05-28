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

const CreateRequest = `
{
    "router": {
        "name": "DummyRouter",
        "area": "JPEAST",
        "userIpAddress": "10.100.0.0/27",
        "redundant": true
    }
}
`

const CreateResponse = `
{
  "router": {
    "area": "JPWEST",
    "firewalls": [
      {
        "id": "F042000000118",
        "isActivated": false
      }
    ],
    "id": "F022000000118",
    "name": "cli-test",
    "nats": [
      {
        "id": "F052000000118",
        "isActivated": false
      }
    ],
    "operationId": "c8e1f724e9e04e899f7a80067c5eb64c",
    "operationStatus": "Processing",
    "redundant": true,
    "routingGroups": [
      {
        "name": "group_1"
      },
      {
        "name": "group_2"
      },
      {
        "name": "group_7"
      },
      {
        "name": "group_6"
      },
      {
        "name": "group_8"
      },
      {
        "name": "group_4"
      },
      {
        "name": "group_3"
      },
      {
        "name": "group_5"
      }
    ],
    "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
    "userIpAddress": "10.100.0.0/27"
  }
}
`

const CreateOutput = `
id             tenantID                          name      area    userIPAddress  redundant  operationStatus  firewalls                nats                     operationID
--             --------                          ----      ----    -------------  ---------  ---------------  ---------                ----                     -----------
F022000000118  75ece6039bfd4d2f88c8b8d790482886  cli-test  JPWEST  10.100.0.0/27  true       Processing       [{F042000000118 false}]  [{F052000000118 false}]  c8e1f724e9e04e899f7a80067c5eb64c
`

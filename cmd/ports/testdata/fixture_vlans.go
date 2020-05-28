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

const VLANsResponse = `
{
  "port": {
    "area": "JPEAST",
    "id": "F012000000172",
    "isActivated": true,
    "location": "NTTComTokyo(NW1)",
    "name": "shark-test",
    "operationStatus": "Completed",
    "portType": "1G",
    "switchName": "lxea02comnw1",
    "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
    "vlanRanges": [
      "2145-2160"
    ],
    "vlans": [
      {
        "status": "unused",
        "vid": 2145
      },
      {
        "status": "unused",
        "vid": 2146
      },
      {
        "status": "unused",
        "vid": 2147
      },
      {
        "status": "unused",
        "vid": 2148
      },
      {
        "status": "unused",
        "vid": 2149
      },
      {
        "status": "unused",
        "vid": 2150
      },
      {
        "status": "unused",
        "vid": 2151
      },
      {
        "status": "unused",
        "vid": 2152
      },
      {
        "status": "unused",
        "vid": 2153
      },
      {
        "status": "unused",
        "vid": 2154
      },
      {
        "status": "unused",
        "vid": 2155
      },
      {
        "status": "unused",
        "vid": 2156
      },
      {
        "status": "unused",
        "vid": 2157
      },
      {
        "status": "unused",
        "vid": 2158
      },
      {
        "status": "unused",
        "vid": 2159
      },
      {
        "status": "unused",
        "vid": 2160
      }
    ]
  }
}
`

const VLANsOutput = `
vid   status
---   ------
2145  unused
2146  unused
2147  unused
2148  unused
2149  unused
2150  unused
2151  unused
2152  unused
2153  unused
2154  unused
2155  unused
2156  unused
2157  unused
2158  unused
2159  unused
2160  unused
`

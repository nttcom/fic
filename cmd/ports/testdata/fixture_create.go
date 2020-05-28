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

const CreateNumOfVLANsRequest = `
{
    "port": {
        "name": "DummyPort",
        "switchName": "lxea02comnw1",
        "numOfVlans": 16,
        "portType": "1G"
    }
}
`

const CreateNumOfVLANsResponse = `
{
  "port": {
    "area": "JPEAST",
    "id": "F012000000168",
    "isActivated": false,
    "location": "NTTComTokyo(NW1)",
    "name": "shark-test",
    "operationId": "ed9906743905491e9413e7d5745a6005",
    "operationStatus": "Processing",
    "portType": "1G",
    "switchName": "lxea02comnw1",
    "tenantId": "75ece6039bfd4d2f88c8b8d790482886",
    "vlanRanges": [
      "1025-1040"
    ],
    "vlans": [
      {
        "status": "unused",
        "vid": 1025
      },
      {
        "status": "unused",
        "vid": 1026
      },
      {
        "status": "unused",
        "vid": 1027
      },
      {
        "status": "unused",
        "vid": 1028
      },
      {
        "status": "unused",
        "vid": 1029
      },
      {
        "status": "unused",
        "vid": 1030
      },
      {
        "status": "unused",
        "vid": 1031
      },
      {
        "status": "unused",
        "vid": 1032
      },
      {
        "status": "unused",
        "vid": 1033
      },
      {
        "status": "unused",
        "vid": 1034
      },
      {
        "status": "unused",
        "vid": 1035
      },
      {
        "status": "unused",
        "vid": 1036
      },
      {
        "status": "unused",
        "vid": 1037
      },
      {
        "status": "unused",
        "vid": 1038
      },
      {
        "status": "unused",
        "vid": 1039
      },
      {
        "status": "unused",
        "vid": 1040
      }
    ]
  }
}
`

const CreateNumOfVLANsOutput = `
id             name        operationStatus  isActivated  vlanRanges   tenantID                          switchName    portType  location          area    operationID
--             ----        ---------------  -----------  ----------   --------                          ----------    --------  --------          ----    -----------
F012000000168  shark-test  Processing       false        [1025-1040]  75ece6039bfd4d2f88c8b8d790482886  lxea02comnw1  1G        NTTComTokyo(NW1)  JPEAST  ed9906743905491e9413e7d5745a6005
`

const CreateVLANRangesRequest = `
{
    "port": {
        "name": "DummyPort",
        "switchName": "lxea02comnw1",
        "vlanRanges": [
            "2145-2160"
        ],
        "portType": "1G"
    }
}
`

const CreateVLANRangesResponse = `
{
  "port": {
    "area": "JPEAST",
    "id": "F012000000169",
    "isActivated": false,
    "location": "NTTComTokyo(NW1)",
    "name": "shark-test",
    "operationId": "83f9e4421f0c49db87de0c58f325a82a",
    "operationStatus": "Processing",
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

const CreateVLANRangesOutput = `
id             name        operationStatus  isActivated  vlanRanges   tenantID                          switchName    portType  location          area    operationID
--             ----        ---------------  -----------  ----------   --------                          ----------    --------  --------          ----    -----------
F012000000169  shark-test  Processing       false        [2145-2160]  75ece6039bfd4d2f88c8b8d790482886  lxea02comnw1  1G        NTTComTokyo(NW1)  JPEAST  83f9e4421f0c49db87de0c58f325a82a
`

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

const DeleteResponse = `
{
    "globalIpAddressSet": {
        "id": "bf839f90b4ea4b4192ea6c33ebb85dc5",
        "operationId": "76a2e7331a974e4597c5404899c2e6a6",
        "operationStatus": "Completed"
    }
}
`

const DeleteOutput = `
id                                operationID                       operationStatus
--                                -----------                       ---------------
bf839f90b4ea4b4192ea6c33ebb85dc5  76a2e7331a974e4597c5404899c2e6a6  Completed
`

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

const ShowResponse = `
{
    "operation": {
        "commitTime": "2020-02-21T07:28:15Z",
        "error": "",
        "id": "1ce786c33ff44d148e621330b7a7cd72",
        "receptionTime": "2020-02-21T07:27:41Z",
        "requestBody": "{\"port\":{\"name\":\"cli-test\",\"numOfVlans\":16,\"portType\":\"1G\",\"switchName\":\"lxea04comnw1\"}}",
        "requestType": "Create",
        "resourceId": "F012000000330",
        "resourceName": "cli-test",
        "resourceType": "Port",
        "status": "Completed",
        "tenantId": "75ece6039bfd4d2f88c8b8d790482886"
    }
}
`

const ShowOutput = `
id                                tenantID                          resourceID     resourceName  resourceType  requestType  status     receptionTime         commitTime            error  requestBody
--                                --------                          ----------     ------------  ------------  -----------  ------     -------------         ----------            -----  -----------
1ce786c33ff44d148e621330b7a7cd72  75ece6039bfd4d2f88c8b8d790482886  F012000000330  cli-test      Port          Create       Completed  2020-02-21T07:27:41Z  2020-02-21T07:28:15Z         {"port":{"name":"cli-test","numOfVlans":16,"portType":"1G","switchName":"lxea04comnw1"}}
`

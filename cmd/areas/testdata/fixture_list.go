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
  "areas": [
    {
      "description": "Japan East",
      "enabled": true,
      "id": "85d77ef691284a1a8dfcd1aaaf5490fc",
      "isPublic": true,
      "name": "JPEAST",
      "number": 1
    },
    {
      "description": "Japan West",
      "enabled": true,
      "id": "e76a147c54584d49a36fb4c83ba4a91b",
      "isPublic": true,
      "name": "JPWEST",
      "number": 2
    }
  ]
}
`

const ListOutput = `
id                                name    number  isPublic  enabled
--                                ----    ------  --------  -------
85d77ef691284a1a8dfcd1aaaf5490fc  JPEAST  1       true      true
e76a147c54584d49a36fb4c83ba4a91b  JPWEST  2       true      true
`

// "enabled": "true" should be "enabled": true
const InvalidListResponse = `
{
  "areas": [
    {
      "description": "Japan East",
      "enabled": "true",
      "id": "85d77ef691284a1a8dfcd1aaaf5490fc",
      "isPublic": true,
      "name": "JPEAST",
      "number": 1
    }
  ]
}
`

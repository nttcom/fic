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

package router_to_port_connections

import "strconv"

var validGroupNames = []string{"group_1", "group_2", "group_3", "group_4", "group_5", "group_6", "group_7", "group_8"}
var validPrepends = []string{"null", "1", "2", "3", "4", "5"}
var validRouteFilterIns = []string{"fullRoute", "noRoute"}
var validRouteFilterOuts = []string{"fullRoute", "fullRouteWithDefaultRoute", "noRoute"}
var validBandwidths = []string{"10M", "20M", "30M", "40M", "50M", "100M", "200M", "300M", "400M", "500M", "1G", "2G", "3G", "4G", "5G", "10G"}
var validPrimarySingleMEDs = []int{10, 20, 30, 40}
var validPrimaryPairedMEDs = []int{10, 30}

func convPrepend(rawPrepend string) *int {
	if rawPrepend == "null" {
		return nil
	}
	prepend, _ := strconv.Atoi(rawPrepend)
	return &prepend
}

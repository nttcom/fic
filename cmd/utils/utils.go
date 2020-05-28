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

package utils

// StringInSlice returns true if the provided string value
// matches the value of an element in the valid slice
func StringInSlice(value string, validSlice []string) bool {
	for _, validName := range validSlice {
		if value == validName {
			return true
		}
	}
	return false
}

// FindDuplicatedStr tests whether there are duplicated
// values in a string slice
func FindDuplicatedStr(s []string) (string, bool) {
	for i, v := range s {
		for j, w := range s {
			if i != j && v == w {
				return v, true
			}
		}
	}
	return "", false
}

func IntInSlice(value int, validSlice []int) bool {
	for _, validName := range validSlice {
		if value == validName {
			return true
		}
	}
	return false
}

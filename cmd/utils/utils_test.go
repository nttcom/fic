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

import (
	"testing"
)

func TestFindDuplicatedStr(t *testing.T) {
	testCases := []struct {
		name     string
		s        []string
		expected struct {
			matched string
			ok      bool
		}
	}{
		{
			"With duplicated elements in the same slice",
			[]string{"foo", "foo", "bar", "baz"},
			struct {
				matched string
				ok      bool
			}{
				"foo",
				true,
			},
		},
		{
			"Without duplicated elements",
			[]string{"foo", "bar", "baz", "qux"},
			struct {
				matched string
				ok      bool
			}{
				"",
				false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matched, ok := FindDuplicatedStr(tc.s)
			if matched != tc.expected.matched {
				t.Errorf("Expected %s, but got %s", tc.expected.matched, matched)
			}
			if ok != tc.expected.ok {
				t.Errorf("Expected %v, but got %v", tc.expected.ok, ok)
			}
		})
	}
}

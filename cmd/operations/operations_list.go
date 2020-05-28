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

package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/operations"
	"github.com/spf13/cobra"
)

// newListCmd creates a new `fic operations list` command
func newListCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List operations",
		Example: "fic operations list",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			pages, err := operations.List(client, nil).AllPages()
			if err != nil {
				var e *json.UnmarshalTypeError
				if errors.As(err, &e) {
					return fmt.Errorf("extracting operations from API response: %w", err)
				}
				return fmt.Errorf("calling List operations API: %w", err)
			}

			opes, _ := operations.ExtractOperations(pages)

			t := utils.NewTabby(out)
			t.AddHeader("id", "tenantID", "resourceID", "resourceName", "resourceType", "requestType", "status",
				"receptionTime", "commitTime", "error", "requestBody")
			for _, ope := range opes {
				t.AddLine(ope.ID, ope.TenantID, ope.ResourceID, ope.ResourceName, ope.ResourceType, ope.RequestType, ope.Status,
					ope.ReceptionTime, ope.CommitTime, ope.Error, ope.RequestBody)
			}
			t.Print()

			return nil
		},
	}

	return cmd
}

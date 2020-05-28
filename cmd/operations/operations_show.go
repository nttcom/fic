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
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/operations"
	"github.com/spf13/cobra"
)

// newShowCmd creates a new `fic operations show` command
func newShowCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show <id>",
		Short:   "Show operation",
		Example: "fic operations show 1ce786c33ff44d148e621330b7a7cd72",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			ope, err := operations.Get(client, args[0]).Extract()
			if err != nil {
				return fmt.Errorf("calling Show operation API: %w", err)
			}

			t := utils.NewTabby(out)
			t.AddHeader("id", "tenantID", "resourceID", "resourceName", "resourceType", "requestType", "status",
				"receptionTime", "commitTime", "error", "requestBody")
			t.AddLine(ope.ID, ope.TenantID, ope.ResourceID, ope.ResourceName, ope.ResourceType, ope.RequestType, ope.Status,
				ope.ReceptionTime, ope.CommitTime, ope.Error, ope.RequestBody)
			t.Print()

			return nil
		},
	}

	return cmd
}

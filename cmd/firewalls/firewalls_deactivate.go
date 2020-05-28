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

package firewalls

import (
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/firewalls"
	"github.com/spf13/cobra"
)

// newDeactivateCmd creates a new `fic firewalls deactivate` command
func newDeactivateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var routerID string

	cmd := &cobra.Command{
		Use:     "deactivate <id>",
		Short:   "Deactivate firewall",
		Example: "fic firewalls deactivate F040123456789 --router F022000000335",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			firewallID := args[0]

			firewall, err := firewalls.Deactivate(client, routerID, firewallID).Extract()
			if err != nil {
				return fmt.Errorf("calling Deactivate firewall API: %w", err)
			}

			t := utils.NewTabby(out)
			t.AddHeader("id", "tenantId", "redundant", "isActivated", "operationStatus", "operationId")
			t.AddLine(firewall.ID, firewall.TenantID, firewall.Redundant, firewall.IsActivated,
				firewall.OperationStatus, firewall.OperationID)
			t.Print()
			return nil
		},
	}
	cmd.Flags().StringVar(&routerID, "router", "", "(required) Router ID")
	cmd.MarkFlagRequired("router")

	return cmd
}

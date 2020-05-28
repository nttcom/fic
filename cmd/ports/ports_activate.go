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

package ports

import (
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/ports"
	"github.com/spf13/cobra"
)

// newsActivateCmd creates a new `fic ports activate` command
func newActivateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "activate <id>",
		Short:   "Activate port",
		Example: "fic ports activate F012000000168",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			port, err := ports.Activate(client, args[0]).Extract()
			if err != nil {
				return fmt.Errorf("calling Activate port API: %w", err)
			}

			t := utils.NewTabby(out)
			t.AddHeader("id", "name", "operationStatus", "isActivated", "vlanRanges", "tenantID", "switchName",
				"portType", "location", "area", "operationID")
			t.AddLine(port.ID, port.Name, port.OperationStatus, port.IsActivated, port.VLANRanges, port.TenantID, port.SwitchName,
				port.PortType, port.Location, port.Area, port.OperationID)
			t.Print()

			return nil
		},
	}

	return cmd
}

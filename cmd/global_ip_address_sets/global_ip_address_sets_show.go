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

package global_ip_address_sets

import (
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/nat_global_ip_address_sets"
	"github.com/spf13/cobra"
)

// newShowCmd creates a new `fic global-ip-address-sets show` command
func newShowCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var (
		routerID string
		natID    string
	)

	cmd := &cobra.Command{
		Use:     "show <id>",
		Short:   "Show global IP address set",
		Example: "fic global-ip-address-sets show bf839f90b4ea4b4192ea6c33ebb85dc5 --router F022000000013 --nat F052000000013",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			ip, err := nat_global_ip_address_sets.Get(client, routerID, natID, args[0]).Extract()
			if err != nil {
				return fmt.Errorf("calling Show global IP address set API: %w", err)
			}

			t := utils.NewTabby(out)
			t.AddHeader("id", "name", "type", "natComponentID", "operationStatus", "tenantID", "numOfAddresses")
			t.AddLine(ip.ID, ip.Name, ip.Type, ip.NATComponentID, ip.OperationStatus, ip.TenantID, ip.NumberOfAddresses)
			t.Print()

			return nil
		},
	}

	cmd.Flags().StringVar(&routerID, "router", "", "(Required) Router ID")
	cmd.Flags().StringVar(&natID, "nat", "", "(Required) NAT ID")

	cmd.MarkFlagRequired("router")
	cmd.MarkFlagRequired("nat")

	return cmd
}

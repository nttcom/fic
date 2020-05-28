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
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/nat_global_ip_address_sets"
	"github.com/spf13/cobra"
)

// newListCmd creates a new `fic global-ip-address-sets list` command
func newListCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var (
		routerID string
		natID    string
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List global IP address sets",
		Example: "fic global-ip-address-sets list --router F022000000013 --nat F052000000013",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			pages, err := nat_global_ip_address_sets.List(client, routerID, natID, nil).AllPages()
			if err != nil {
				var e *json.UnmarshalTypeError
				if errors.As(err, &e) {
					return fmt.Errorf("extracting global IP address sets from API response: %w", err)
				}
				return fmt.Errorf("calling List global IP address sets API: %w", err)
			}

			ips, _ := nat_global_ip_address_sets.ExtractGlobalIPAddressSets(pages)

			t := utils.NewTabby(out)
			t.AddHeader("id", "name", "type", "natComponentID", "operationStatus", "tenantID", "numOfAddresses")
			for _, ip := range ips {
				t.AddLine(ip.ID, ip.Name, ip.Type, ip.NATComponentID, ip.OperationStatus, ip.TenantID, ip.NumberOfAddresses)
			}
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

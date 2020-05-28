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
	"regexp"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/nat_global_ip_address_sets"
	"github.com/spf13/cobra"
)

// newCreateCmd creates a new `fic global-ip-address-sets create` command
func newCreateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var (
		routerID       string
		natID          string
		natType        string
		numOfAddresses int
	)

	nameRe := regexp.MustCompile(`^[\w&()-]{1,64}$`)

	cmd := &cobra.Command{
		Use:     "create <name>",
		Short:   "Create global IP address set",
		Example: "fic global-ip-address-sets create testSet --router F022000000013 --nat F052000000013 --type sourceNapt --num-of-addresses 5",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !nameRe.MatchString(args[0]) {
				return fmt.Errorf("name must consist of half-width alphanumeric characters and some symbols &()-_ less than 64: received %s", args[0])
			}

			if natType != "sourceNapt" && natType != "destinationNat" {
				return fmt.Errorf("type must be sourceNapt or destinationNat: received %s", natType)
			}

			if natType == "sourceNapt" && (numOfAddresses < 1 || numOfAddresses > 5) {
				return fmt.Errorf("num-of-addresses must be [1-5] with sourceNapt type: received %d", numOfAddresses)
			}

			if natType == "destinationNat" && numOfAddresses != 1 {
				return fmt.Errorf("num-of-addresses must be 1 with destinationNat type: received %d", numOfAddresses)
			}

			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			opts := nat_global_ip_address_sets.CreateOpts{
				Name:              args[0],
				Type:              natType,
				NumberOfAddresses: numOfAddresses,
			}

			ip, err := nat_global_ip_address_sets.Create(client, routerID, natID, opts).Extract()
			if err != nil {
				return fmt.Errorf("calling Create global IP address set API: %w", err)
			}

			t := utils.NewTabby(out)
			t.AddHeader("id", "name", "type", "natComponentID", "operationID", "operationStatus", "tenantID",
				"numOfAddresses")
			t.AddLine(ip.ID, ip.Name, ip.Type, ip.NATComponentID, ip.OperationID, ip.OperationStatus, ip.TenantID,
				ip.NumberOfAddresses)
			t.Print()

			return nil
		},
	}

	cmd.Flags().StringVar(&routerID, "router", "", "(Required) Router ID")
	cmd.Flags().StringVar(&natID, "nat", "", "(Required) NAT ID")
	cmd.Flags().StringVar(&natType, "type", "", "(Required) NAT type, sourceNapt or destinationNat")
	cmd.Flags().IntVar(&numOfAddresses, "num-of-addresses", 0, "(Required) Number of NAT addresses")

	cmd.MarkFlagRequired("router")
	cmd.MarkFlagRequired("nat")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("num-of-addresses")

	return cmd
}

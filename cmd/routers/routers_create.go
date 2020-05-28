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

package routers

import (
	"fmt"
	"io"
	"net"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers"
	"github.com/spf13/cobra"
)

// newCreateCmd creates a new `fic routers create` command
func newCreateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var (
		area          string
		userIPAddress string
		redundant     bool
	)

	cmd := &cobra.Command{
		Use:     "create <name>",
		Short:   "Create router",
		Example: "fic routers create testRouter --area JPEAST --user-ip-address 10.100.0.0/27 --redundant",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, ipNet, err := net.ParseCIDR(userIPAddress)
			if err != nil {
				return fmt.Errorf("user-ip-address must be CIDR whose subnet mask is 27, e.g. 10.100.0.0/27: received %s", userIPAddress)
			}
			if ipNet.Mask.String() != "ffffffe0" {
				return fmt.Errorf("subnet mask of user-ip-address must be 27, e.g. 10.100.0.0/27: received %s", userIPAddress)
			}

			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			opts := routers.CreateOpts{
				Name:          args[0],
				Area:          area,
				UserIPAddress: userIPAddress,
				Redundant:     &redundant,
			}

			router, err := routers.Create(client, opts).Extract()
			if err != nil {
				return fmt.Errorf("calling Create router API: %w", err)
			}

			t := utils.NewTabby(out)
			t.AddHeader("id", "tenantID", "name", "area", "userIPAddress", "redundant",
				"operationStatus", "firewalls", "nats", "operationID")
			t.AddLine(router.ID, router.TenantID, router.Name, router.Area, router.UserIPAddress, *router.Redundant,
				router.OperationStatus, router.Firewalls, router.NATs, router.OperationID)
			t.Print()

			return nil
		},
	}

	cmd.Flags().StringVar(&area, "area", "", "(Required) Area name")
	cmd.Flags().StringVar(&userIPAddress, "user-ip-address", "", "(Required) CIDR whose subnet mask is 27, e.g. 10.100.0.0/27")
	cmd.Flags().BoolVar(&redundant, "redundant", false, "Enable redundant configuration")

	cmd.MarkFlagRequired("area")
	cmd.MarkFlagRequired("user-ip-address")

	return cmd
}

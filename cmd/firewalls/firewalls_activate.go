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
	"net"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/firewalls"
	"github.com/spf13/cobra"
)

// newActivateCmd creates a new `fic firewalls activate` command
func newActivateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var routerID string
	var userIPAddresses []string

	cmd := &cobra.Command{
		Use:   "activate <id>",
		Short: "Activate firewall",
		Example: "fic firewalls activate F040123456789 --router F022000000335 " +
			"--user-ip-addresses 192.168.0.0/30,192.168.0.4/30,192.168.0.8/30,192.168.0.12/30",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			firewallID := args[0]

			if ipLength := len(userIPAddresses); !(ipLength == 0 || ipLength == 4 || ipLength == 8) {
				return fmt.Errorf("user-ip-address must have length of 4 or 8")
			}

			for _, userIPAddress := range userIPAddresses {
				_, ipNet, err := net.ParseCIDR(userIPAddress)
				if err != nil {
					return fmt.Errorf("user-ip-addresses must be CIDR whose subnet mask is 30, e.g. 192.168.0.0/30: received %s", userIPAddress)
				}
				submaskLength, _ := ipNet.Mask.Size()
				if submaskLength != 30 {
					return fmt.Errorf("subnet mask of user-ip-addresses must be 30, e.g. 192.168.0.0/30: received %s", userIPAddress)
				}
			}

			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			activateOpt := firewalls.ActivateOpts{
				UserIPAddresses: userIPAddresses,
			}

			firewall, err := firewalls.Activate(client, routerID, firewallID, activateOpt).Extract()
			if err != nil {
				return fmt.Errorf("calling Activate firewall API: %w", err)
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
	cmd.Flags().StringSliceVar(&userIPAddresses, "user-ip-addresses", []string{}, "Local IP addresses")
	cmd.MarkFlagRequired("router")

	return cmd
}

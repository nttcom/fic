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

package router_to_port_connections

import (
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	conn "github.com/nttcom/go-fic/fic/eri/v1/router_paired_to_port_connections"
	"github.com/spf13/cobra"
)

// newDestinationCmd creates a new `fic router-to-port-connections destination` command
func newDestinationCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destination <connection_id>",
		Short:   "Show destination belonging to connection",
		Example: "fic router-to-port-connections destination F030123456789",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			c, err := conn.Get(client, args[0]).Extract()
			if err != nil {
				return fmt.Errorf("calling Show router to port connection API: %w", err)
			}

			t := utils.NewTabby(out)
			if !c.Redundant {
				t.AddHeader("tenantId", "primary.portId", "primary.vlan", "primary.ipAddress", "primary.asn")
				t.AddLine(
					c.TenantID,
					c.Destination.Primary.PortID, c.Destination.Primary.VLAN, c.Destination.Primary.IPAddress, c.Destination.Primary.ASN,
				)
				t.Print()

				return nil
			}

			t.AddHeader(
				"tenantId",
				"primary.portId", "primary.vlan", "primary.ipAddress", "primary.asn",
				"secondary.portId", "secondary.vlan", "secondary.ipAddress", "secondary.asn",
			)
			t.AddLine(
				c.TenantID,
				c.Destination.Primary.PortID, c.Destination.Primary.VLAN, c.Destination.Primary.IPAddress, c.Destination.Primary.ASN,
				c.Destination.Secondary.PortID, c.Destination.Secondary.VLAN, c.Destination.Secondary.IPAddress, c.Destination.Secondary.ASN,
			)
			t.Print()

			return nil
		},
	}

	return cmd
}

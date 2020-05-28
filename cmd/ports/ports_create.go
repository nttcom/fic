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
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/ports"
	"github.com/spf13/cobra"
)

// newCreateCmd creates a new `fic ports create` command
func newCreateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var (
		switchName    string
		numberOfVLANs int
		vlanRanges    []string
		portType      string
	)

	cmd := &cobra.Command{
		Use:     "create <name>",
		Short:   "Create port",
		Example: "fic ports create testPort --switch lxea02comnw1 --number-of-vlans 16 --type 1G",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if portType != "1G" && portType != "10G" {
				return errors.New("type must be 1G or 10G")
			}

			if numberOfVLANs == 0 && len(vlanRanges) == 0 || numberOfVLANs != 0 && len(vlanRanges) != 0 {
				return errors.New("one of number-of-vlans and vlan-ranges must be set")
			}

			if numberOfVLANs != 0 {
				if numberOfVLANs < 16 {
					return errors.New("number-of-vlans must be greater than or equal to 16")
				}
				if numberOfVLANs > 512 {
					return errors.New("number-of-vlans must be less than or equal to 512")
				}
				if numberOfVLANs%16 != 0 {
					return errors.New("number-of-vlans must be multiple of 16")
				}
			}

			if len(vlanRanges) != 0 {
				if len(vlanRanges) > 32 {
					return fmt.Errorf("length of vlan-ranges must be less than or equal to 32: received %d", len(vlanRanges))
				}
				for _, vlanRange := range vlanRanges {
					vlans := strings.Split(vlanRange, "-")
					if len(vlans) != 2 {
						return fmt.Errorf("vlan-ranges element must be <start>-<end> format, e.g. 17-32 and 49-64: received %s", vlanRange)
					}

					start, err := strconv.Atoi(vlans[0])
					if err != nil {
						return fmt.Errorf("vlan-ranges element must be integer range, e.g. 17-32 and 49-64: received %s", vlanRange)
					}
					end, err := strconv.Atoi(vlans[1])
					if err != nil {
						return fmt.Errorf("vlan-ranges element must be integer range, e.g. 17-32 and 49-64: received %s", vlanRange)
					}

					if start <= 0 || end <= 0 || start%16 != 1 || end-start != 15 {
						return fmt.Errorf("vlan-ranges element must be from 16a+1 to 16a+16(a is natural number), e.g. 17-32 and 49-64: received %s", vlanRange)
					}
				}
			}

			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			opts := ports.CreateOpts{
				Name:          args[0],
				SwitchName:    switchName,
				NumberOfVLANs: numberOfVLANs,
				VLANRanges:    vlanRanges,
				PortType:      portType,
			}
			port, err := ports.Create(client, opts).Extract()
			if err != nil {
				return fmt.Errorf("calling Create port API: %w", err)
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

	cmd.Flags().StringVar(&switchName, "switch", "", "(Required) Switch name")
	cmd.Flags().IntVar(&numberOfVLANs, "number-of-vlans", 0, "Number of VLANs, multiple of 16")
	cmd.Flags().StringSliceVar(&vlanRanges, "vlan-ranges", []string{}, "List of VLAN range, e.g. 17-32 and 49-64")
	cmd.Flags().StringVar(&portType, "type", "", "(Required) Port type, 1G or 10G")

	cmd.MarkFlagRequired("switch")
	cmd.MarkFlagRequired("type")

	return cmd
}

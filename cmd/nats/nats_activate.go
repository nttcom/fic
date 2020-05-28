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

package nats

import (
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/nats"
	"github.com/spf13/cobra"
)

// newActivateCmd creates a new `fic nats activate` command
func newActivateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var (
		routerID               string
		userIPAddresses        []string
		rawGlobalIPAddressSets []string
	)
	regexpIPName := regexp.MustCompile(`^[\w&()-]{1,64}$`)
	cmd := &cobra.Command{
		Use:   "activate <id>",
		Short: "Activate nat",
		Example: "fic nats activate F052000000335 --router F022000000335 " +
			"--user-ip-addresses 192.168.0.0/30,192.168.0.4/30,192.168.0.8/30,192.168.0.12/30 " +
			"--global-ip-address-sets src-set-01,sourceNapt,5 --global-ip-address-sets dst-set-01,destinationNat,1",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var globalIPAddressSets []nats.GlobalIPAddressSet

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

			var naptTypeCount, destTypeCount int
			for _, rgip := range rawGlobalIPAddressSets {
				gip := strings.Split(strings.TrimSpace(rgip), ",")
				if len(gip) != 3 {
					return fmt.Errorf("global-ip-address-sets must have format like <name>,<type>,<numberOfAddresses>, e.g. src-set-01,sourceNapt,5")
				}

				if !regexpIPName.MatchString(gip[0]) {
					return fmt.Errorf("name of global-ip-address-sets must be composed of alpha-numeric characters and & ( ) - _, and must have maximum length of 64 as well")
				}

				numberOfAddresses, err := strconv.Atoi(gip[2])
				if err != nil {
					return fmt.Errorf("numberOfAddresses of global-ip-address-sets must be numeric value")
				}

				switch gip[1] {
				case "sourceNapt":
					naptTypeCount += 1
					if numberOfAddresses < 1 || numberOfAddresses > 5 {
						return fmt.Errorf("numberOfAddresses of global-ip-address-sets must be in the range of 1 to 5 in case sourceNapt")
					}
				case "destinationNat":
					destTypeCount += 1
					if numberOfAddresses != 1 {
						return fmt.Errorf("numberOfAddresses of global-ip-address-sets must be 1 in case destinationNat")
					}
				default:
					return fmt.Errorf("type of global-ip-address-sets must be either sourceNpat or destinationNat")
				}

				if naptTypeCount > 8 {
					return fmt.Errorf("maximum number of sourceNapt type is 8")
				}

				if destTypeCount > 30 {
					return fmt.Errorf("maximum number of destinationNat type is 30")
				}

				g := nats.GlobalIPAddressSet{
					Name:              gip[0],
					Type:              gip[1],
					NumberOfAddresses: numberOfAddresses,
				}
				globalIPAddressSets = append(globalIPAddressSets, g)
			}

			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			opts := nats.ActivateOpts{
				UserIPAddresses:     userIPAddresses,
				GlobalIPAddressSets: globalIPAddressSets,
			}

			n, err := nats.Activate(client, routerID, args[0], opts).Extract()
			if err != nil {
				return fmt.Errorf("calling Activate nat API: %w", err)
			}

			t := utils.NewTabby(out)
			t.AddHeader(
				"id",
				"tenantID",
				"redundant",
				"isActivated",
				"operationID",
				"operationStatus",
			)
			t.AddLine(
				n.ID,
				n.TenantID,
				n.Redundant,
				n.IsActivated,
				n.OperationID,
				n.OperationStatus,
			)
			t.Print()

			return nil
		},
	}

	cmd.Flags().StringVar(&routerID, "router", "", "(Required) Router ID")

	cmd.Flags().StringSliceVar(
		&userIPAddresses,
		"user-ip-addresses",
		[]string{},
		"List of CIDR of user ip addresses, e.g. 192.168.0.0/30 and 192.168.4.0/30")

	cmd.Flags().StringArrayVar(
		&rawGlobalIPAddressSets,
		"global-ip-address-sets",
		[]string{},
		"(Required) List of global ip address sets(<name>,<type>,<numOfAddresses>), e.g. src-set-01,sourceNapt,5 and dst-set-01,destinationNat,1")

	cmd.MarkFlagRequired("router")
	cmd.MarkFlagRequired("global-ip-address-sets")

	return cmd
}

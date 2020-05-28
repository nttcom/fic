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
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	pairedConn "github.com/nttcom/go-fic/fic/eri/v1/router_paired_to_port_connections"
	singleConn "github.com/nttcom/go-fic/fic/eri/v1/router_single_to_port_connections"
	"github.com/spf13/cobra"
)

// newCreateCmd creates a new `fic router-to-port-connections create` command
func newCreateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var (
		srcRouterID    string
		srcGroupName   string
		srcPrimary     string
		srcSecondary   string
		srcRouteFilter string
		destPrimary    string
		destSecondary  string
		bandwidth      string
	)

	r := regexp.MustCompile(`^[\w&()-]{1,64}$`)
	validateSrc := func(splitSrc []string, isPrimary, hasSecondary bool) error {
		_, ipNet, err := net.ParseCIDR(splitSrc[0])
		if err != nil {
			return fmt.Errorf("ipAddress must be CIDR whose subnet mask is 30, e.g. 10.0.0.1/30: received %s", splitSrc[0])
		}
		subNetMaskLength, _ := ipNet.Mask.Size()
		if subNetMaskLength != 30 {
			return fmt.Errorf("subnet mask of ipAddress must be 30, e.g. 10.0.0.1/30: received %s", splitSrc[0])
		}

		if !utils.StringInSlice(splitSrc[1], validPrepends) {
			return fmt.Errorf("asPathPrepend.in must be one of %s: received %s", validPrepends, splitSrc[1])
		}

		if !utils.StringInSlice(splitSrc[2], validPrepends) {
			return fmt.Errorf("asPathPrepend.out must be one of %s: received %s", validPrepends, splitSrc[2])
		}

		if isPrimary {
			med, err := strconv.Atoi(splitSrc[3])
			if err != nil {
				return fmt.Errorf("med.out must be numeric value: received %s", splitSrc[3])
			}

			if hasSecondary {
				if !utils.IntInSlice(med, validPrimaryPairedMEDs) {
					return fmt.Errorf("med.out in paired connection must be one of %v: received %d", validPrimaryPairedMEDs, med)
				}
			} else {
				if !utils.IntInSlice(med, validPrimarySingleMEDs) {
					return fmt.Errorf("med.out in single connection must be one of %v: received %d", validPrimarySingleMEDs, med)
				}
			}
		}

		return nil
	}
	validateDest := func(splitDest []string) error {
		vlan, err := strconv.Atoi(splitDest[1])
		if err != nil {
			return fmt.Errorf("vlan must be numeric value: received %s", splitDest[1])
		}
		if vlan < 101 || vlan > 3300 {
			return fmt.Errorf("vlan must be range of 101 to 3300: received %s", splitDest[1])
		}

		_, ipNet, err := net.ParseCIDR(splitDest[2])
		if err != nil {
			return fmt.Errorf("ipAddress must be CIDR whose subnet mask is 30, e.g. 10.0.0.2/30: received %s", splitDest[2])
		}
		subNetMaskLength, _ := ipNet.Mask.Size()
		if subNetMaskLength != 30 {
			return fmt.Errorf("subnet mask of ipAddress must be 30, e.g. 10.0.0.2/30: received %s", splitDest[2])
		}

		asn, err := strconv.Atoi(splitDest[3])
		if err != nil {
			return fmt.Errorf("asn must be numeric value: received %s", splitDest[3])
		}
		if asn < 1 || asn > 65535 {
			return fmt.Errorf("asn must be range of 1 to 65535: received %s", splitDest[3])
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create router to port connection",
		Example: "# In case of non paired-connection \n" +
			"fic router-to-port-connections create YourConnectionName " +
			"--source-router F020123456789 " +
			"--source-group group_1 " +
			"--source-primary 10.0.0.1/30,4,4,10 " +
			"--source-route-filter fullRoute,fullRouteWithDefaultRoute " +
			"--destination-primary F010123456789,101,10.0.0.2/30,65000 " +
			"--bandwidth 10M \n\n" +
			"# In case of paired-connection \n" +
			"fic router-to-port-connections create YourConnectionName " +
			"--source-router F020123456789 " +
			"--source-group group_1 " +
			"--source-primary 10.0.0.1/30,4,4,10 " +
			"--source-secondary 10.0.0.5/30,2,1 " +
			"--source-route-filter fullRoute,fullRouteWithDefaultRoute " +
			"--destination-primary F010123456789,101,10.0.0.2/30,65000 " +
			"--destination-secondary F019876543210,102,10.0.0.6/30,65000 " +
			"--bandwidth 10M",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !r.MatchString(args[0]) {
				return fmt.Errorf("name of router to port connection must be composed of alpha-numeric "+
					"characters and & ( ) - _, and must have maximum length of 64 as well: received %s", args[0])
			}

			if !utils.StringInSlice(srcGroupName, validGroupNames) {
				return fmt.Errorf("source-group must be one of %s: received %s", validGroupNames, srcGroupName)
			}

			splitSrcPrimary := strings.Split(strings.TrimSpace(srcPrimary), ",")
			splitSrcSecondary := strings.Split(strings.TrimSpace(srcSecondary), ",")
			splitSrcRouteFilter := strings.Split(strings.TrimSpace(srcRouteFilter), ",")
			splitDestPrimary := strings.Split(strings.TrimSpace(destPrimary), ",")
			splitDestSecondary := strings.Split(strings.TrimSpace(destSecondary), ",")

			if len(splitSrcPrimary) != 4 {
				return fmt.Errorf("source-primary must have format like "+
					"<ipAddress>,<asPathPrepend.in>,<asPathPrepend.out>,<med.out>: received %s", srcPrimary)
			}
			if err := validateSrc(splitSrcPrimary, true, srcSecondary != ""); err != nil {
				return fmt.Errorf("in source-primary, %w", err)
			}

			if len(splitDestPrimary) != 4 {
				return fmt.Errorf("destination-primary must have format like "+
					"<portId>,<vlan>,<ipAddress>,<asn>: received %s", destPrimary)
			}
			if err := validateDest(splitDestPrimary); err != nil {
				return fmt.Errorf("in destination-primary, %w", err)
			}

			if srcSecondary != "" {
				if len(splitSrcSecondary) != 3 {
					return fmt.Errorf("source-secondary must have format like "+
						"<ipAddress>,<asPathPrepend.in>,<asPathPrepend.out>: received %s", srcSecondary)
				}
				if err := validateSrc(splitSrcSecondary, false, true); err != nil {
					return fmt.Errorf("in source-secondary, %w", err)
				}
			}

			redundant := false
			if destSecondary != "" {
				if len(splitDestSecondary) != 4 {
					return fmt.Errorf("destination-secondary must have format like "+
						"<portId>,<vlan>,<ipAddress>,<asn>: received %s", destSecondary)
				}
				if err := validateDest(splitDestSecondary); err != nil {
					return fmt.Errorf("in destination-secondary, %w", err)
				}
				redundant = true
			}

			if len(splitSrcRouteFilter) != 2 {
				return fmt.Errorf("source-route-filter must have format like "+
					"<routeFilter.in>,<routeFilter.out>: received %s", srcRouteFilter)
			}
			if !utils.StringInSlice(splitSrcRouteFilter[0], validRouteFilterIns) {
				return fmt.Errorf("routeFilter.in must be one of %s: received %s", validRouteFilterIns, splitSrcRouteFilter[0])
			}
			if !utils.StringInSlice(splitSrcRouteFilter[1], validRouteFilterOuts) {
				return fmt.Errorf("routeFilter.out must be one of %s: received %s", validRouteFilterOuts, splitSrcRouteFilter[1])
			}

			if !utils.StringInSlice(bandwidth, validBandwidths) {
				return fmt.Errorf("bandwidth must be one of %s: received %s", validBandwidths, bandwidth)
			}

			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			prependIn, prependOut := interface{}(convPrepend(splitSrcPrimary[1])), interface{}(convPrepend(splitSrcPrimary[2]))
			med, _ := strconv.Atoi(splitSrcPrimary[3])
			vlan, _ := strconv.Atoi(splitDestPrimary[1])

			t := utils.NewTabby(out)
			t.AddHeader("id", "name", "redundant", "tenantId", "area", "operationStatus", "bandwidth", "operationId")

			if !redundant {
				opts := singleConn.CreateOpts{
					Name: args[0],
					Source: singleConn.Source{
						RouterID:  srcRouterID,
						GroupName: srcGroupName,
						RouteFilter: singleConn.RouteFilter{
							In:  splitSrcRouteFilter[0],
							Out: splitSrcRouteFilter[1],
						},
						Primary: singleConn.SourceHAInfo{
							IPAddress: splitSrcPrimary[0],
							ASPathPrepend: singleConn.ASPathPrepend{
								In:  &prependIn,
								Out: &prependOut,
							},
							MED: &singleConn.MED{
								Out: med,
							},
						},
					},
					Destination: singleConn.Destination{
						Primary: singleConn.DestinationHAInfo{
							PortID:    splitDestPrimary[0],
							VLAN:      vlan,
							IPAddress: splitDestPrimary[2],
							ASN:       splitDestPrimary[3],
						},
					},
					Bandwidth: bandwidth,
				}

				c, err := singleConn.Create(client, opts).Extract()
				if err != nil {
					return fmt.Errorf("calling Create router to port connection API: %w", err)
				}

				t.AddLine(c.ID, c.Name, c.Redundant, c.TenantID, c.Area, c.OperationID, c.OperationStatus, c.Bandwidth)
				t.Print()

				return nil
			}

			secondPrependIn, secondPrependOut := interface{}(convPrepend(splitSrcSecondary[1])), interface{}(convPrepend(splitSrcSecondary[2]))
			secondVlan, _ := strconv.Atoi(splitDestSecondary[1])

			opts := pairedConn.CreateOpts{
				Name: args[0],
				Source: pairedConn.Source{
					RouterID:  srcRouterID,
					GroupName: srcGroupName,
					RouteFilter: pairedConn.RouteFilter{
						In:  splitSrcRouteFilter[0],
						Out: splitSrcRouteFilter[1],
					},
					Primary: pairedConn.SourceHAInfo{
						IPAddress: splitSrcPrimary[0],
						ASPathPrepend: pairedConn.ASPathPrepend{
							In:  &prependIn,
							Out: &prependOut,
						},
						MED: &pairedConn.MED{
							Out: med,
						},
					},
					Secondary: pairedConn.SourceHAInfo{
						IPAddress: splitSrcSecondary[0],
						ASPathPrepend: pairedConn.ASPathPrepend{
							In:  &secondPrependIn,
							Out: &secondPrependOut,
						},
					},
				},
				Destination: pairedConn.Destination{
					Primary: pairedConn.DestinationHAInfo{
						PortID:    splitDestPrimary[0],
						VLAN:      vlan,
						IPAddress: splitDestPrimary[2],
						ASN:       splitDestPrimary[3],
					},
					Secondary: pairedConn.DestinationHAInfo{
						PortID:    splitDestSecondary[0],
						VLAN:      secondVlan,
						IPAddress: splitDestSecondary[2],
						ASN:       splitDestSecondary[3],
					},
				},
				Bandwidth: bandwidth,
			}

			c, err := pairedConn.Create(client, opts).Extract()
			if err != nil {
				return fmt.Errorf("calling Create router to port connection API: %w", err)
			}

			t.AddLine(c.ID, c.Name, c.Redundant, c.TenantID, c.Area, c.OperationID, c.OperationStatus, c.Bandwidth)
			t.Print()

			return nil
		},
	}

	cmd.Flags().StringVar(&srcRouterID, "source-router", "", "(Required) Router ID belonging to source")
	cmd.Flags().StringVar(&srcGroupName, "source-group", "", "(Required) Group Name belonging to source")
	cmd.Flags().StringVar(
		&srcPrimary,
		"source-primary",
		"",
		"(Required) Source Primary Info specified in the format <ipAddress>,<asPathPrepend.in>,<asPathPrepend.out>,<med.out>")
	cmd.Flags().StringVar(
		&srcSecondary,
		"source-secondary",
		"",
		"Source Secondary Info specified in the format <ipAddress>,<asPathPrepend.in>,<asPathPrepend.out>")
	cmd.Flags().StringVar(
		&srcRouteFilter,
		"source-route-filter",
		"",
		"(Required) Set of BGP Filter Ingress and Egress specified in the format <routeFilter.in>,<routeFilter.out>")
	cmd.Flags().StringVar(
		&destPrimary,
		"destination-primary",
		"",
		"(Required) Destination Primary Info specified in the format <portId>,<vlan>,<ipAddress>,<asn>")
	cmd.Flags().StringVar(
		&destSecondary,
		"destination-secondary",
		"",
		"Destination Secondary Info specified in the format <portId>,<vlan>,<ipAddress>,<asn>")
	cmd.Flags().StringVar(&bandwidth, "bandwidth", "", "(Required) Bandwidth of router to port connection")

	cmd.MarkFlagRequired("source-router")
	cmd.MarkFlagRequired("source-group")
	cmd.MarkFlagRequired("source-primary")
	cmd.MarkFlagRequired("source-rout-filter")
	cmd.MarkFlagRequired("destination-primary")
	cmd.MarkFlagRequired("bandwidth")

	return cmd
}

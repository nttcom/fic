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
	"strconv"
	"strings"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	pairedConn "github.com/nttcom/go-fic/fic/eri/v1/router_paired_to_port_connections"
	singleConn "github.com/nttcom/go-fic/fic/eri/v1/router_single_to_port_connections"
	"github.com/spf13/cobra"
)

// newUpdateCmd creates a new `fic router-to-port-connections update` command
func newUpdateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var (
		srcPrimary     string
		srcSecondary   string
		srcRouteFilter string
	)

	validateSrc := func(splitSrc []string, isPrimary, hasSecondary bool) error {
		if !utils.StringInSlice(splitSrc[0], validPrepends) {
			return fmt.Errorf("asPathPrepend.in must be one of %s: received %s", validPrepends, splitSrc[0])
		}

		if !utils.StringInSlice(splitSrc[1], validPrepends) {
			return fmt.Errorf("asPathPrepend.out must be one of %s: received %s", validPrepends, splitSrc[1])
		}

		if isPrimary {
			med, err := strconv.Atoi(splitSrc[2])
			if err != nil {
				return fmt.Errorf("med.out must be numeric value: received %s", splitSrc[2])
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

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update router to port connection",
		Example: "# In case of non paired-connection \n" +
			"fic router-to-port-connections update F030123456789 " +
			"--source-primary null,2,10 " +
			"--source-route-filter noRoute,fullRoute \n\n" +
			"# In case of paired-connection \n" +
			"fic router-to-port-connections update F030123456789 " +
			"--source-primary null,2,10 " +
			"--source-secondary null,2 " +
			"--source-route-filter noRoute,fullRoute",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			splitSrcPrimary := strings.Split(strings.TrimSpace(srcPrimary), ",")
			splitSrcSecondary := strings.Split(strings.TrimSpace(srcSecondary), ",")
			splitSrcRouteFilter := strings.Split(strings.TrimSpace(srcRouteFilter), ",")

			if len(splitSrcPrimary) != 3 {
				return fmt.Errorf("source-primary must have format like "+
					"<asPathPrepend.in>,<asPathPrepend.out>,<med.out>: received %s", srcPrimary)
			}
			if err := validateSrc(splitSrcPrimary, true, srcSecondary != ""); err != nil {
				return fmt.Errorf("in source-primary, %w", err)
			}

			redundant := false
			if srcSecondary != "" {
				if len(splitSrcSecondary) != 2 {
					return fmt.Errorf("source-secondary must have format like "+
						"<asPathPrepend.in>,<asPathPrepend.out>: received %s", srcSecondary)
				}
				if err := validateSrc(splitSrcSecondary, false, true); err != nil {
					return fmt.Errorf("in source-secondary, %w", err)
				}
				redundant = true
			}

			if len(splitSrcRouteFilter) != 2 {
				return fmt.Errorf("source-route-filter must have format like "+
					"<routeFilter.in>:<routeFilter.out>: received %s", srcRouteFilter)
			}
			if !utils.StringInSlice(splitSrcRouteFilter[0], validRouteFilterIns) {
				return fmt.Errorf("routeFilter.in must be one of %s: received %s", validRouteFilterIns, splitSrcRouteFilter[0])
			}
			if !utils.StringInSlice(splitSrcRouteFilter[1], validRouteFilterOuts) {
				return fmt.Errorf("routeFilter.out must be one of %s: received %s", validRouteFilterOuts, splitSrcRouteFilter[1])
			}

			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			prependIn, prependOut := interface{}(convPrepend(splitSrcPrimary[0])), interface{}(convPrepend(splitSrcPrimary[1]))
			med, _ := strconv.Atoi(splitSrcPrimary[2])

			t := utils.NewTabby(out)
			t.AddHeader("id", "name", "redundant", "tenantId", "area", "operationStatus", "bandwidth", "operationId")

			if !redundant {
				opts := singleConn.UpdateOpts{
					Source: singleConn.SourceForUpdate{
						RouteFilter: singleConn.RouteFilter{
							In:  splitSrcRouteFilter[0],
							Out: splitSrcRouteFilter[1],
						},
						Primary: singleConn.SourceHAInfoForUpdate{
							ASPathPrepend: singleConn.ASPathPrepend{
								In:  &prependIn,
								Out: &prependOut,
							},
							MED: &singleConn.MED{
								Out: med,
							},
						},
					},
				}

				c, err := singleConn.Update(client, args[0], opts).Extract()
				if err != nil {
					return fmt.Errorf("calling Update router to port connection API: %w", err)
				}
				t.AddLine(c.ID, c.Name, c.Redundant, c.TenantID, c.Area, c.OperationID, c.OperationStatus, c.Bandwidth)
				t.Print()

				return nil
			}

			secondPrependIn, secondPrependOut := interface{}(convPrepend(splitSrcSecondary[0])), interface{}(convPrepend(splitSrcSecondary[1]))

			opts := pairedConn.UpdateOpts{
				Source: pairedConn.SourceForUpdate{
					RouteFilter: pairedConn.RouteFilter{
						In:  splitSrcRouteFilter[0],
						Out: splitSrcRouteFilter[1],
					},
					Primary: pairedConn.SourceHAInfoForUpdate{
						ASPathPrepend: pairedConn.ASPathPrepend{
							In:  &prependIn,
							Out: &prependOut,
						},
						MED: &pairedConn.MED{
							Out: med,
						},
					},
					Secondary: pairedConn.SourceHAInfoForUpdate{
						ASPathPrepend: pairedConn.ASPathPrepend{
							In:  &secondPrependIn,
							Out: &secondPrependOut,
						},
					},
				},
			}

			c, err := pairedConn.Update(client, args[0], opts).Extract()
			if err != nil {
				return fmt.Errorf("calling Update router to port connection API: %w", err)
			}
			t.AddLine(c.ID, c.Name, c.Redundant, c.TenantID, c.Area, c.OperationID, c.OperationStatus, c.Bandwidth)
			t.Print()

			return nil
		},
	}

	cmd.Flags().StringVar(
		&srcPrimary,
		"source-primary",
		"",
		"(Required) Source Primary Info specified in the format <asPathPrepend.in>,<asPathPrepend.out>,<med.out>")
	cmd.Flags().StringVar(
		&srcSecondary,
		"source-secondary",
		"",
		"Source Secondary Info specified in the format <asPathPrepend.in>,<asPathPrepend.out>")
	cmd.Flags().StringVar(
		&srcRouteFilter,
		"source-route-filter",
		"",
		"(Required) Set of BGP Filter Ingress and BGP Filter Egress specified in the format <routeFilter.in>,<routeFilter.out>")

	cmd.MarkFlagRequired("source-primary")
	cmd.MarkFlagRequired("source-route-filter")

	return cmd
}

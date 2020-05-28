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
	"strings"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/nats"
	"github.com/spf13/cobra"
)

// newUpdateCmd creates a new `fic nats update` command
func newUpdateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var (
		routerID               string
		rawSourceNAPTRules     []string
		rawDestinationNATRules []string
	)

	validGroupNames := []string{"group_1", "group_2", "group_3", "group_4"}
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update nat",
		Example: "fic nats update F052000000335 --router F022000000335 " +
			"--source-napt-rules group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04 " +
			"--source-napt-rules group_2,group_1,src-set-05,src-set-06,src-set-07,src-set-08 " +
			"--destination-nat-rules group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32 " +
			"--destination-nat-rules group_2,group_1,dst-set-03:192.168.0.3/32,dst-set-04:192.168.0.4/32",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var sourceNAPTRules []nats.SourceNAPTRule
			var destinationNATRules []nats.DestinationNATRule

			for _, rawSourceNAPTRule := range rawSourceNAPTRules {
				srcNAPT := strings.Split(strings.TrimSpace(rawSourceNAPTRule), ",")
				if len(srcNAPT) < 3 {
					return fmt.Errorf("source-napt-rules must have format like " +
						"from, to, sourceNaptRules.entries.0.then.0, ... ,sourceNaptRules.entries.0.then.N")
				}

				for _, v := range srcNAPT[:2] {
					if !utils.StringInSlice(v, validGroupNames) {
						return fmt.Errorf("first and second element of source-napt-rules must have format like group_[1|2|3|4]")
					}
				}

				if len(srcNAPT) > 10 {
					return fmt.Errorf("maximum length of entries of source-napt-rules is 8")
				}

				src := nats.SourceNAPTRule{
					From: []string{srcNAPT[0]},
					To:   srcNAPT[1],
					Entries: []nats.EntryInSourceNAPTRule{
						{
							Then: srcNAPT[2:],
						},
					},
				}
				sourceNAPTRules = append(sourceNAPTRules, src)
			}

			for _, rawDestinationNATRule := range rawDestinationNATRules {
				destNAT := strings.Split(strings.TrimSpace(rawDestinationNATRule), ",")

				if len(destNAT) < 3 {
					return fmt.Errorf("destination-nat-rules must have format like " +
						"from, to, entries.0.match.destinationAddress:entries.0.then, ... ,entries.N.match.destinationAddress:entries.N.then")
				}

				for _, v := range destNAT[:2] {
					if !utils.StringInSlice(v, validGroupNames) {
						return fmt.Errorf("first and second element of destination-nat-rules must have format like group_[1|2|3|4]")
					}
				}

				if len(destNAT) > 32 {
					return fmt.Errorf("maximum length of entries of destination-nat-rules is 30")
				}

				var destinationNATEntries []nats.EntryInDestinationNATRule
				for _, entryPart := range destNAT[2:] {

					entry := strings.Split(entryPart, ":")

					if len(entry) != 2 {
						return fmt.Errorf("entries of destination-nat-rules must have format like " +
							"entries.0.match.destinationAddress:entries.0.then, ... ,entries.N.match.destinationAddress:entries.N.then")
					}

					_, ipNet, err := net.ParseCIDR(entry[1])
					if err != nil {
						return fmt.Errorf("destinationAddress must be CIDR whose subnet mask is 32, e.g. 192.168.0.1/32: received %s", entry[1])
					}
					submaskLength, _ := ipNet.Mask.Size()
					if submaskLength != 32 {
						return fmt.Errorf("subnet mask of destinationAddress must be 32, e.g. 192.168.0.1/32: received %s", entry[1])
					}

					destinationNATEntry := nats.EntryInDestinationNATRule{
						Match: nats.Match{
							DestinationAddress: entry[0],
						},
						Then: entry[1],
					}
					destinationNATEntries = append(destinationNATEntries, destinationNATEntry)
				}

				dst := nats.DestinationNATRule{
					From:    destNAT[0],
					To:      destNAT[1],
					Entries: destinationNATEntries,
				}
				destinationNATRules = append(destinationNATRules, dst)
			}

			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			opts := nats.UpdateOpts{
				SourceNAPTRules:     sourceNAPTRules,
				DestinationNATRules: destinationNATRules,
			}

			n, err := nats.Update(client, routerID, args[0], opts).Extract()
			if err != nil {
				return fmt.Errorf("calling Update nat API: %w", err)
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

	cmd.Flags().StringArrayVar(
		&rawSourceNAPTRules,
		"source-napt-rules",
		[]string{},
		"(Required) List of source NAPT rules(<from>,<to>,<entries.then>,<entries.then>,...), e.g. group_1,group_2,src-set-01,src-set-02,src-set-03,src-set-04",
	)

	cmd.Flags().StringArrayVar(
		&rawDestinationNATRules,
		"destination-nat-rules",
		[]string{},
		"(Required) List of destination NAT rules(<from>,<to>,<entries.match>:<entries.then>,<entries.match>:<entries.then>,...), e.g. group_1,group_2,dst-set-01:192.168.0.1/32,dst-set-02:192.168.0.2/32",
	)

	cmd.MarkFlagRequired("router")
	cmd.MarkFlagRequired("source-napt-rules")
	cmd.MarkFlagRequired("destination-nat-rules")

	return cmd
}

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

package switches

import (
	"errors"
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/switches"
	"github.com/spf13/cobra"
)

// newVLANRangesCmd creates a new `fic switches vlan-ranges` command
func newVLANRangesCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vlan-ranges <switch_id>",
		Short:   "List VLAN ranges belonging to switch",
		Example: "fic switches vlan-ranges 6c1b07d60d844b58bb2277190dc5dd91",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			pages, err := switches.List(client, nil).AllPages()
			if err != nil {
				return fmt.Errorf("calling List switches API: %w", err)
			}

			sws, _ := switches.ExtractSwitches(pages)

			for _, sw := range sws {
				if sw.ID == args[0] {
					t := utils.NewTabby(out)
					t.AddHeader("vlanRange", "available")
					for _, vlanRange := range sw.VLANRanges {
						t.AddLine(vlanRange.Range, vlanRange.Available)
					}
					t.Print()

					return nil
				}
			}

			return errors.New("switch not found")
		},
	}

	return cmd
}

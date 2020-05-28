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

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/firewalls"
	"github.com/spf13/cobra"
)

// newCustomApplicationCmd creates a new `fic firewalls custom-application` command
func newCustomApplicationCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var routerID string

	cmd := &cobra.Command{
		Use:     "custom-applications <id>",
		Short:   "Show Custom Applications of Firewall",
		Example: "fic firewalls custom-applications F040123456789 --router F022000000335",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			firewallID := args[0]

			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			fw, err := firewalls.Get(client, routerID, firewallID).Extract()
			if err != nil {
				return fmt.Errorf("calling Get firewall API: %w", err)
			}

			t := utils.NewTabby(out)
			t.AddHeader("name", "protocol", "destinationPort")
			for _, v := range fw.CustomApplications {
				t.AddLine(v.Name, v.Protocol, v.DestinationPort)
			}
			t.Print()
			return nil
		},
	}
	cmd.Flags().StringVar(&routerID, "router", "", "(required) Router ID")
	cmd.MarkFlagRequired("router")

	return cmd
}

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
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/firewalls"
	"github.com/spf13/cobra"
)

// newListCmd creates a new `fic firewalls list` command
func newListCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var routerID string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List firewalls",
		Example: "fic firewalls list --router F022000000335",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			pages, err := firewalls.List(client, routerID, nil).AllPages()
			if err != nil {
				var e *json.UnmarshalTypeError
				if errors.As(err, &e) {
					return fmt.Errorf("extracting firewalls from API response: %w", err)
				}
				return fmt.Errorf("calling List firewall API: %w", err)
			}

			allFirewalls, _ := firewalls.ExtractFirewalls(pages)

			t := utils.NewTabby(out)
			t.AddHeader("id", "tenantId", "redundant", "isActivated", "operationStatus")
			for _, firewall := range allFirewalls {
				t.AddLine(firewall.ID, firewall.TenantID, firewall.IsActivated, firewall.OperationStatus)
			}
			t.Print()

			return nil
		},
	}

	cmd.Flags().StringVar(&routerID, "router", "", "(required) Router ID")
	cmd.MarkFlagRequired("router")

	return cmd
}

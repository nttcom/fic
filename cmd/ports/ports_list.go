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
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/ports"
	"github.com/spf13/cobra"
)

// newListCmd creates a new `fic ports list` command
func newListCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List ports",
		Example: "fic ports list",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			pages, err := ports.List(client, nil).AllPages()
			if err != nil {
				var e *json.UnmarshalTypeError
				if errors.As(err, &e) {
					return fmt.Errorf("extracting ports from API response: %w", err)
				}
				return fmt.Errorf("calling List ports API: %w", err)
			}

			ps, _ := ports.ExtractPorts(pages)

			t := utils.NewTabby(out)
			t.AddHeader("id", "name", "operationStatus", "isActivated", "vlanRanges", "tenantID", "switchName",
				"portType", "location", "area")
			for _, p := range ps {
				t.AddLine(p.ID, p.Name, p.OperationStatus, p.IsActivated, p.VLANRanges, p.TenantID, p.SwitchName,
					p.PortType, p.Location, p.Area)
			}
			t.Print()

			return nil
		},
	}

	return cmd
}

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

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers/components/nats"
	"github.com/spf13/cobra"
)

// newDeactivateCmd creates a new `fic nats deactivate` command
func newDeactivateCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	var routerID string

	cmd := &cobra.Command{
		Use:     "deactivate <id>",
		Short:   "Deactivate nat",
		Example: "fic nats deactivate F052000000335 --router F022000000335",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			n, err := nats.Deactivate(client, routerID, args[0]).Extract()
			if err != nil {
				return fmt.Errorf("calling Deactivate nat API: %w", err)
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
	cmd.MarkFlagRequired("router")

	return cmd
}

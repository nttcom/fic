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

package routers

import (
	"fmt"
	"io"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers"
	"github.com/spf13/cobra"
)

// newRoutingGroupsCmd creates a new `fic routers routing-groups` command
func newRoutingGroupsCmd(clientFn func() (*fic.ServiceClient, error), out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "routing-groups <router_id>",
		Short:   "List routing groups belonging to router",
		Example: "fic routers routing-groups F022000000118",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			router, err := routers.Get(client, args[0]).Extract()
			if err != nil {
				return fmt.Errorf("calling Show router API: %w", err)
			}

			t := utils.NewTabby(out)
			t.AddHeader("name")
			for _, rg := range router.RoutingGroups {
				t.AddLine(rg.Name)
			}
			t.Print()

			return nil
		},
	}

	return cmd
}

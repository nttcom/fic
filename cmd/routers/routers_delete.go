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

	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/eri/v1/routers"
	"github.com/spf13/cobra"
)

// newsDeleteCmd creates a new `fic routers delete` command
func newDeleteCmd(clientFn func() (*fic.ServiceClient, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <id>",
		Short:   "Delete router",
		Example: "fic routers delete F022000000118",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFn()
			if err != nil {
				return fmt.Errorf("creating FIC client: %w", err)
			}

			if err := routers.Delete(client, args[0]).ExtractErr(); err != nil {
				return fmt.Errorf("calling Delete router API: %w", err)
			}

			return nil
		},
	}

	return cmd
}

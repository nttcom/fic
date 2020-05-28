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

package global_ip_address_sets

import (
	"os"

	"github.com/nttcom/fic/cmd/utils"
	"github.com/spf13/cobra"
)

// NewCmd creates a new `fic global-ip-address-sets` command
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "global-ip-address-sets",
		Short: "Global IP address set management",
	}

	cmd.AddCommand(
		newCreateCmd(utils.NewClient, os.Stdout),
		newDeleteCmd(utils.NewClient, os.Stdout),
		newShowCmd(utils.NewClient, os.Stdout),
		newListCmd(utils.NewClient, os.Stdout),
		newAddressesCmd(utils.NewClient, os.Stdout),
	)

	return cmd
}
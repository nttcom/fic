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

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the `completion` command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates completion scripts",
}

// bashCompletionCmd represents the `completion bash` command
var bashCompletionCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RootCmd.GenBashCompletion(os.Stdout)
	},
}

// zshCompletionCmd represents the `completion zsh` command
var zshCompletionCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generates zsh completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RootCmd.GenZshCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(bashCompletionCmd)
	completionCmd.AddCommand(zshCompletionCmd)

	RootCmd.AddCommand(completionCmd)
}

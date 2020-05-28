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
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/nttcom/fic/cmd/areas"
	"github.com/nttcom/fic/cmd/firewalls"
	"github.com/nttcom/fic/cmd/global_ip_address_sets"
	"github.com/nttcom/fic/cmd/nats"
	"github.com/nttcom/fic/cmd/operations"
	"github.com/nttcom/fic/cmd/ports"
	"github.com/nttcom/fic/cmd/router_to_port_connections"
	"github.com/nttcom/fic/cmd/routers"
	"github.com/nttcom/fic/cmd/switches"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	debug   bool
)

// RootCmd represents the base `fic` command
var RootCmd = &cobra.Command{
	Use:           "fic",
	Short:         "Command line interface for Flexible InterConnect",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, disableLog)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fic.yaml)")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")

	RootCmd.AddCommand(
		areas.NewCmd(),
		firewalls.NewCmd(),
		global_ip_address_sets.NewCmd(),
		nats.NewCmd(),
		operations.NewCmd(),
		ports.NewCmd(),
		router_to_port_connections.NewCmd(),
		routers.NewCmd(),
		switches.NewCmd(),
	)

	viper.SetDefault("auth_url", "https://api.ntt.com/keystone/v3")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".fic")
	}

	viper.SetEnvPrefix("fic")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

// disableLog discards go-fic logs without debug mode.
func disableLog() {
	if !debug {
		log.SetOutput(ioutil.Discard)
	}
}

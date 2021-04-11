// Copyright 2021 Charles Francoise
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "scrt",
	Short:   "A secret manager for the command-line",
	Version: "0.0.0",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet("password") {
			return fmt.Errorf("missing password")
		}

		cmd.SilenceUsage = true

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(unsetCmd)

	rootCmd.PersistentFlags().StringP("password", "p", "", "master password to unlock the store")
	err := viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	if err != nil {
		panic(err)
	}

	viper.SetEnvPrefix("scrt")
	viper.AutomaticEnv()
}

// Execute executes the root cobra command
func Execute() error {
	return rootCmd.Execute()
}

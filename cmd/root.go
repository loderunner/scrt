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
	"os"
	"path/filepath"

	"github.com/loderunner/scrt/backend"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "none"
)

var rootCmd = &cobra.Command{
	Use:   "scrt",
	Short: "A secret manager for the command-line",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Read configuration from .scrt file if exists, recursively searching
		// for .scrt file in parent directories until root is reached
		viper.SetConfigName(".scrt")
		viper.SetConfigType("yaml")
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		viper.AddConfigPath(dir)
		for dir != "/" {
			dir = filepath.Dir(dir)
			viper.AddConfigPath(dir)
		}

		err = viper.ReadInConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return err
			}
		}

		// Validate configuration
		if !viper.IsSet(configKeyPassword) {
			return fmt.Errorf("missing password")
		}
		if !viper.IsSet(configKeyStorage) {
			return fmt.Errorf("missing storage type")
		}
		if !viper.IsSet(configKeyLocation) {
			return fmt.Errorf("missing store location")
		}

		storage := viper.GetString(configKeyStorage)
		factory, ok := backend.Backends[storage]
		if !ok {
			return fmt.Errorf("unknown storage type: %s", storage)
		}

		// Add backend flags to command's flagset, bind to config and re-parse
		cmd.FParseErrWhitelist.UnknownFlags = false
		cmd.Flags().AddFlagSet(factory.Flags())
		err = viper.BindPFlags(cmd.Flags())
		if err != nil {
			return err
		}
		err = cmd.ParseFlags(os.Args[1:])
		if err != nil {
			return cmd.FlagErrorFunc()(cmd, err)
		}

		// Silence usage on error, since errors are runtime, not config, from
		// this point onwards
		cmd.SilenceUsage = true

		return nil
	},
}

func addCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
	cmd.FParseErrWhitelist.UnknownFlags = true
}

func init() {
	if version == "dev" {
		rootCmd.Version = fmt.Sprintf("%s-%s", version, commit)
	} else {
		rootCmd.Version = version
	}

	addCommand(initCmd)
	addCommand(setCmd)
	addCommand(getCmd)
	addCommand(unsetCmd)

	rootCmd.PersistentFlags().StringP("password", "p", "", "master password to unlock the store")
	err := viper.BindPFlag(configKeyPassword, rootCmd.PersistentFlags().Lookup("password"))
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().String("storage", "", "storage type")
	err = viper.BindPFlag(configKeyStorage, rootCmd.PersistentFlags().Lookup("storage"))
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().String("location", "", "store location")
	err = viper.BindPFlag(configKeyLocation, rootCmd.PersistentFlags().Lookup("location"))
	if err != nil {
		panic(err)
	}

	viper.SetEnvPrefix("scrt")
	viper.AutomaticEnv()

	cobra.OnInitialize()
}

// Execute executes the root cobra command
func Execute() error {
	return rootCmd.Execute()
}

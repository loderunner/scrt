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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/loderunner/scrt/backend"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "scrt",
	Short:   "A secret manager for the command-line",
	Version: "0.0.0",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Read configuration from .scrt file if exists, recursively searching
		// for .scrt file in parent directories until root is reached
		viper.SetConfigType("yaml")
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		for {
			configPath := filepath.Join(dir, ".scrt")
			fileInfo, err := os.Stat(configPath)
			if err == nil {
				// .scrt exists at path

				if fileInfo.IsDir() {
					return fmt.Errorf("%s is a directory", configPath)
				}

				f, err := os.Open(configPath)
				if err != nil {
					return err
				}
				defer f.Close()

				err = viper.ReadConfig(f)
				if err != nil {
					return fmt.Errorf("could not read %s: %w", configPath, err)
				}

				// config loaded, break out of loop
				break
			}

			// .scrt does not exist at path

			if !errors.Is(err, os.ErrNotExist) {
				return err
			}

			if dir == "/" {
				// reached root path, no .scrt to load, break out of loop
				break
			}

			dir = filepath.Dir(dir)
		}

		// Validate configuration
		if !viper.IsSet(configKeyPassword) {
			return fmt.Errorf("missing password")
		}
		if !viper.IsSet(configKeyStorage) {
			return fmt.Errorf("missing store type")
		}
		if !viper.IsSet(configKeyLocation) {
			return fmt.Errorf("missing store location")
		}

		storage := viper.GetString(configKeyStorage)
		if _, ok := backend.Backends[storage]; !ok {
			return fmt.Errorf("unknown storage type: %s", storage)
		}

		// Silence usage on error, since errors are runtime, not config, from
		// this point onwards
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
	err := viper.BindPFlag(configKeyPassword, rootCmd.PersistentFlags().Lookup("password"))
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().String("storage", "", "storage type")
	err = viper.BindPFlag(configKeyStorage, rootCmd.PersistentFlags().Lookup("storage"))
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().String("location", "", "storage location")
	err = viper.BindPFlag(configKeyLocation, rootCmd.PersistentFlags().Lookup("location"))
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

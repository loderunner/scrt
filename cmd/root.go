// Copyright 2021-2023 Charles Francoise
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
	"reflect"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/loderunner/scrt/backend"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var verbose bool

type fielder struct {
	fields map[string]interface{}
}

func (f fielder) Fields() log.Fields {
	return f.fields
}

// RootCmd is the root command for scrt.
var RootCmd = &cobra.Command{
	Use:   "scrt",
	Short: "A secret manager for the command-line",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Short circuit for storage command
		if cmd == storageCmd {
			return nil
		}

		err := readConfig(cmd)
		if err != nil {
			return err
		}

		// Set logger in context
		if verbose {
			logger = &log.Logger{Handler: cli.Default}
		} else {
			logger = &log.Logger{Handler: discard.Default}
		}
		cmdContext = log.NewContext(
			cmdContext,
			logger,
		)

		// Validate configuration
		if !viper.IsSet(configKeyStorage) {
			for k := range backend.Backends {
				if viper.InConfig(k) {
					viper.Set(configKeyStorage, k)
					break
				}
			}
			if !viper.IsSet(configKeyStorage) {
				return fmt.Errorf("missing storage type")
			}
		}
		if !viper.IsSet(configKeyPassword) {
			return fmt.Errorf("missing password")
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

		// Log configuration
		if viper.ConfigFileUsed() != "" {
			logger.
				WithField("path", viper.ConfigFileUsed()).
				Infof("read configuration file")
		}
		settings := make(map[string]interface{})
		for k, v := range viper.AllSettings() {
			// Do not log password or unset settings
			if k != configKeyPassword && !reflect.ValueOf(v).IsZero() {
				settings[k] = v
			}
		}
		logger.
			WithFields(fielder{fields: settings}).
			Info("using configuration")

		// Silence usage on error, since errors are runtime, not config, from
		// this point onwards
		cmd.SilenceUsage = true

		return nil
	},
}

func readConfig(cmd *cobra.Command) error {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Read configuration from .scrt file if exists, recursively searching
		// for .scrt file in parent directories until root is reached
		viper.SetConfigName(".scrt")
		viper.SetConfigType("yaml")
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		viper.AddConfigPath(dir)
		for parentDir := filepath.Dir(dir); dir != parentDir; parentDir = filepath.Dir(dir) {
			dir = parentDir
			viper.AddConfigPath(dir)
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return err
		}
	}

	return nil
}

func addCommand(cmd *cobra.Command) {
	RootCmd.AddCommand(cmd)
	cmd.FParseErrWhitelist.UnknownFlags = true
}

func init() {
	cobra.EnableCommandSorting = false

	addCommand(initCmd)
	addCommand(setCmd)
	addCommand(getCmd)
	addCommand(listCmd)
	addCommand(unsetCmd)
	addCommand(storageCmd)
	addCommand(exportCmd)

	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "configuration file")
	RootCmd.PersistentFlags().StringP("password", "p", "", "master password to unlock the store")
	err := viper.BindPFlag(configKeyPassword, RootCmd.PersistentFlags().Lookup("password"))
	if err != nil {
		panic(err)
	}
	RootCmd.PersistentFlags().String("storage", "", "storage type")
	err = viper.BindPFlag(configKeyStorage, RootCmd.PersistentFlags().Lookup("storage"))
	if err != nil {
		panic(err)
	}
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	viper.SetEnvPrefix("scrt")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	cobra.OnInitialize()
}

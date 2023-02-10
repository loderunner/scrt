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
	"fmt"
	"os"
	"strings"

	"github.com/loderunner/scrt/backend"
	"github.com/loderunner/scrt/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var OutFile string

var exportCmd = &cobra.Command{
	Use:   "export [--file|-f <file>]",
	Short: "Export all the keys in a store to an environment file",
	Args: func(cmd *cobra.Command, args []string) error {
		errIsEmpty := cobra.ExactArgs(0)(cmd, args)
		errIsTwo := cobra.ExactArgs(2)(cmd, args)

		if errIsEmpty != nil && errIsTwo != nil {
			return fmt.Errorf("export requires 0 or 2 arguments")
		}

		if len(args) == 2 {
			if args[0] != "-f" && args[0] != "--file" {
				return fmt.Errorf("export requires -f or --file as the first argument")
			}

			OutFile = args[1]
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		storage := viper.GetString(configKeyStorage)

		b, err := backend.Backends[storage].NewContext(cmdContext, viper.AllSettings())
		if err != nil {
			return err
		}

		// print(OutFile)

		exists, err := b.ExistsContext(cmdContext)
		if err != nil {
			return fmt.Errorf("could not check store existence: %w", err)
		}
		if !exists {
			return fmt.Errorf("store does not exist")
		}

		data, err := b.LoadContext(cmdContext)
		if err != nil {
			return fmt.Errorf("could not load data from store: %w", err)
		}

		password := []byte(viper.GetString(configKeyPassword))
		s, err := store.ReadStoreContext(cmdContext, password, data)
		if err != nil {
			return fmt.Errorf("could not read store from data: %w", err)
		}

		keys := s.ListContext(cmdContext)

		var sb strings.Builder

		for _, key := range keys {
			sb.WriteString(key)
			sb.WriteString("=")

			val, err := s.GetContext(cmdContext, key)

			if err != nil {
				return fmt.Errorf("could not get value for key %s: %w", key, err)
			}

			sb.WriteString(string(val))
			sb.WriteString("\n")
		}

		err = os.WriteFile(OutFile, []byte(sb.String()), 0644)

		if err != nil {
			return fmt.Errorf("could not write file %s: %w", OutFile, err)
		}

		return nil
	},
}

func init() {
	exportCmd.Flags().StringVarP(&OutFile, "file", "f", ".env", "file to write to")
}

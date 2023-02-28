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
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/loderunner/scrt/backend"
	"github.com/loderunner/scrt/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var exportCmd = &cobra.Command{
	Use:   "export [flags]",
	Short: "Export all the keys in a store",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var out string
		var format string

		out, err = cmd.Flags().GetString("output")

		if err != nil {
			return fmt.Errorf("could not read options: %w", err)
		}

		format, err = cmd.Flags().GetString("format")

		if err != nil {
			return fmt.Errorf("could not read options: %w", err)
		}

		if format == "" {
			return fmt.Errorf("missing export format")
		}

		if format != "dotenv" && format != "json" && format != "yaml" {
			return fmt.Errorf("invalid export format")
		}

		storage := viper.GetString(configKeyStorage)

		b, err := backend.Backends[storage].NewContext(cmdContext, viper.AllSettings())
		if err != nil {
			return err
		}

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

		exportData := s.Export()

		var marshalBytes []byte
		var marshalErr error

		if format == "json" {
			marshalBytes, marshalErr = json.Marshal(exportData)
		} else if format == "yaml" {
			marshalBytes, marshalErr = yaml.Marshal(exportData)
		} else if format == "dotenv" {
			bytes, err := godotenv.Marshal(exportData)

			if err != nil {
				marshalErr = err
			} else {
				marshalBytes = []byte(bytes)
			}
		}

		if marshalErr != nil {
			return fmt.Errorf("could not marshal data: %w", marshalErr)
		}

		if out == "" {
			fmt.Println(string(marshalBytes))
			return nil
		}

		err = os.WriteFile(out, marshalBytes, 0644)

		if err != nil {
			return fmt.Errorf("could not write file %s: %w", out, err)
		}

		return nil
	},
}

func init() {
	exportCmd.Flags().StringP("output", "o", "", "export to file (defaults to stdout)")

	exportCmd.Flags().StringP("format", "f", "", "export file format (json,yaml,dotenv)")

}

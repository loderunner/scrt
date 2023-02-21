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

var exportCmd = &cobra.Command{
	Use:   "export [flags]",
	Short: "Export all the keys in a store",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var out string
		var format string

		out, err = cmd.Flags().GetString("out")

		if err != nil {
			return fmt.Errorf("could not get out flag: %w", err)
		}

		format, err = cmd.Flags().GetString("format")

		if err != nil {
			return fmt.Errorf("could not get format flag: %w", err)
		}

		if out == "" {
			out = "STD_OUT"
		}

		if format == "" {
			return fmt.Errorf("format is required")
		}

		if format != "dotenv" && format != "json" && format != "yaml" {
			return fmt.Errorf("invalid format %s", format)
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

		keys := s.ListContext(cmdContext)

		var sb strings.Builder

		if format == "json" {
			sb.WriteString("{\n")

			for i, key := range keys {
				if i > 0 {
					sb.WriteString(",\n")
				}

				sb.WriteString(`"` + key + `": `)

				val, err := s.GetContext(cmdContext, key)

				if err != nil {
					return fmt.Errorf("could not get value for key %s: %w", key, err)
				}

				sb.WriteString(`"` + string(val) + `"`)
			}

			sb.WriteString("\n}")
		} else if format == "yaml" {

			sb.WriteString("---\n")

			for _, key := range keys {
				sb.WriteString(key)
				sb.WriteString(": ")

				val, err := s.GetContext(cmdContext, key)

				if err != nil {
					return fmt.Errorf("could not get value for key %s: %w", key, err)
				}

				sb.WriteString(string(val))
				sb.WriteString("\n")
			}

		} else if format == "dotenv" {
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
		}

		if out == "STD_OUT" {
			fmt.Println(sb.String())
			return nil
		}

		err = os.WriteFile(out, []byte(sb.String()), 0644)

		if err != nil {
			return fmt.Errorf("could not write file %s: %w", out, err)
		}

		return nil
	},
}

func init() {
	exportCmd.Flags().StringP("out", "o", "STD_OUT", "export to file (defaults to stdout)")

	exportCmd.Flags().StringP("format", "f", "", "export file format (json,yaml,dotenv)")

}

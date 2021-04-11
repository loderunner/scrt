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
	"gitlab.com/loderunner/ask"

	"github.com/loderunner/scrt/backend"
	"github.com/loderunner/scrt/store"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new store",
	RunE: func(cmd *cobra.Command, args []string) error {
		storage := viper.GetString(configKeyStorage)
		location := viper.GetString(configKeyLocation)

		b := backend.Backends[storage](location)

		if b.Exists() {
			var overwrite bool
			var err error
			if cmd.Flags().Changed("overwrite") {
				overwrite, err = cmd.Flags().GetBool("overwrite")
			} else {
				overwrite, err = ask.Boolf("%s store already exists at %s. Do you want to overwrite it?", storage, location).
					Default(false).
					Ask()
			}
			if err != nil {
				return fmt.Errorf("could not read options: %w", err)
			}
			if !overwrite {
				return fmt.Errorf("aborted")
			}
		}

		s := store.NewStore()
		password := []byte(viper.GetString(configKeyPassword))

		data, err := store.WriteStore(password, s)
		if err != nil {
			return fmt.Errorf("could not write store to data: %w", err)
		}

		err = b.Save(data)
		if err != nil {
			return fmt.Errorf("could not save data to %s: %w", location, err)
		}

		fmt.Printf("%s store initialized at %s\n", storage, location)

		return nil
	},
}

func init() {
	initCmd.Flags().Bool("overwrite", false, "overwrite store if it exists")
}

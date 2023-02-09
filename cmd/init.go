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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/loderunner/scrt/backend"
	"github.com/loderunner/scrt/store"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new store",
	RunE: func(cmd *cobra.Command, args []string) error {
		storage := viper.GetString(configKeyStorage)

		b, err := backend.Backends[storage].NewContext(cmdContext, viper.AllSettings())
		if err != nil {
			return err
		}

		exists, err := b.ExistsContext(cmdContext)
		if err != nil {
			return fmt.Errorf("could not check store existence: %w", err)
		}
		if exists {
			overwrite, err := cmd.Flags().GetBool("overwrite")
			if err != nil {
				return fmt.Errorf("could not read options: %w", err)
			}
			if !overwrite {
				return fmt.Errorf("store already exists, use --overwrite to force init")
			}
			logger.Info("overwriting existing store")
		}

		s := store.NewStoreContext(cmdContext)
		password := []byte(viper.GetString(configKeyPassword))

		data, err := store.WriteStoreContext(cmdContext, password, s)
		if err != nil {
			return fmt.Errorf("could not write store to data: %w", err)
		}

		err = b.SaveContext(cmdContext, data)
		if err != nil {
			return fmt.Errorf("could not save data to store: %w", err)
		}

		fmt.Println("store initialized")

		return nil
	},
}

func init() {
	initCmd.Flags().Bool("overwrite", false, "overwrite store if it exists")
}

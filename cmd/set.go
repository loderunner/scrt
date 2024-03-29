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
	"io"
	"os"

	"github.com/loderunner/scrt/backend"
	"github.com/loderunner/scrt/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setCmd = &cobra.Command{
	Use:   "set [flags] key [value]",
	Short: "Associate a key to a value in a store",
	Long: `Associate a key to a value in a store. If value is omitted from the command
line, it will be read from standard input.`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.MinimumNArgs(1)(cmd, args)
		if err != nil {
			return err
		}
		err = cobra.MaximumNArgs(2)(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		storage := viper.GetString(configKeyStorage)
		key := args[0]

		var val []byte
		var err error
		if len(args) == 1 {
			logger.Info("reading value from standard input")
			val, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("could not read from standard input %w", err)
			}
		} else {
			val = []byte(args[1])
		}

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

		var overwrite bool
		if cmd.Flags().Changed("overwrite") {
			overwrite, err = cmd.Flags().GetBool("overwrite")
			if err != nil {
				return fmt.Errorf("could not read options: %w", err)
			}
		}

		if s.HasContext(cmdContext, key) {
			if !overwrite {
				return fmt.Errorf("value exists for key \"%s\", use --overwrite to force", key)
			}
			logger.WithField("key", key).Info("overwriting existing value")
		}

		err = s.SetContext(cmdContext, key, val)
		if err != nil {
			return fmt.Errorf("could not set value: %w", err)
		}

		data, err = store.WriteStoreContext(cmdContext, password, s)
		if err != nil {
			return fmt.Errorf("could not write store to data: %w", err)
		}

		err = b.SaveContext(cmdContext, data)
		if err != nil {
			return fmt.Errorf("could not save data to store: %w", err)
		}

		return nil
	},
}

func init() {
	setCmd.Flags().Bool("overwrite", false, "overwrite value if it exists")
}

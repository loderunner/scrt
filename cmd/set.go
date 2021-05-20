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
		location := viper.GetString(configKeyLocation)
		key := args[0]

		var val []byte
		var err error
		if len(args) == 1 {
			val, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("could not read from stdin %w", err)
			}
		} else {
			val = []byte(args[1])
		}

		b, err := backend.Backends[storage].New(location, viper.AllSettings())
		if err != nil {
			return err
		}

		exists, err := b.Exists()
		if err != nil {
			return fmt.Errorf("could not check store existence: %w", err)
		}
		if !exists {
			return fmt.Errorf("%s store at %s does not exist", storage, location)
		}

		data, err := b.Load()
		if err != nil {
			return fmt.Errorf("could not load data from %s: %w", location, err)
		}

		password := []byte(viper.GetString(configKeyPassword))
		s, err := store.ReadStore(password, data)
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

		if s.Has(key) && !overwrite {
			return fmt.Errorf("value exists for key \"%s\", use --overwrite to force", key)
		}

		err = s.Set(key, val)
		if err != nil {
			return fmt.Errorf("could not set value: %w", err)
		}

		data, err = store.WriteStore(password, s)
		if err != nil {
			return fmt.Errorf("could not write store to data: %w", err)
		}

		err = b.Save(data)
		if err != nil {
			return fmt.Errorf("could not save data to %s: %w", location, err)
		}

		return nil
	},
}

func init() {
	setCmd.Flags().Bool("overwrite", false, "overwrite value if it exists")
}

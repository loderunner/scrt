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
	"math"

	"github.com/loderunner/scrt/backend"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func padRight(s string, pad string, length int) string {
	for len(s) < length {
		s += pad
	}
	return s[:length]
}

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "List storage types and options",
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(0)(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		maxNameLength := 11
		for name := range backend.Backends {
			if len(name) > maxNameLength {
				maxNameLength = len(name)
			}
		}
		padLength := int(4 * math.Ceil(float64(maxNameLength+1)/4))
		for _, name := range backend.BackendNameList {
			factory := backend.Backends[name]
			fmt.Printf("%s:\n", factory.Name())
			fmt.Printf(
				"  %s%s\n",
				padRight(name, " ", padLength),
				factory.Description(),
			)

			flags := factory.Flags()
			flagCount := 0
			flags.VisitAll(func(_ *pflag.Flag) { flagCount++ })
			if flagCount == 0 {
				fmt.Println()
			} else {
				fmt.Printf("Flags:\n")
				fmt.Printf("%s\n", flags.FlagUsagesWrapped(80))
			}
		}
		return nil
	},
}

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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/loderunner/scrt/backend"
)

func TestRootCmd(t *testing.T) {
	viper.Reset()

	fs := afero.NewMemMapFs()
	viper.SetFs(fs)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	err := RootCmd.PersistentPreRunE(RootCmd, []string{})
	if err == nil {
		t.Fatal("expected error")
	}

	viper.Set(configKeyStorage, "mock")
	err = RootCmd.PersistentPreRunE(RootCmd, []string{})
	if err == nil {
		t.Fatal("expected error")
	}

	viper.Set(configKeyPassword, "")
	err = RootCmd.PersistentPreRunE(RootCmd, []string{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompletionDoesNotRequireStorage(t *testing.T) {
	viper.Reset()

	// Ensure the generated completion command is attached to RootCmd.
	RootCmd.InitDefaultCompletionCmd()

	// Find the generated completion command and one of its subcommands
	// (e.g. bash). Generating shell completion must not require a configured
	// storage backend or password.
	var completionCmd *cobra.Command
	for _, c := range RootCmd.Commands() {
		if c.Name() == "completion" {
			completionCmd = c
			break
		}
	}
	if completionCmd == nil {
		t.Fatal("completion command not found")
	}

	// The completion command itself
	if err := RootCmd.PersistentPreRunE(completionCmd, []string{}); err != nil {
		t.Fatalf("completion command should short-circuit: %v", err)
	}

	// A completion subcommand (bash) must also short-circuit via its parent.
	var bashCmd *cobra.Command
	for _, c := range completionCmd.Commands() {
		if c.Name() == "bash" {
			bashCmd = c
			break
		}
	}
	if bashCmd == nil {
		t.Fatal("completion bash subcommand not found")
	}
	if err := RootCmd.PersistentPreRunE(bashCmd, []string{}); err != nil {
		t.Fatalf("completion bash subcommand should short-circuit: %v", err)
	}
}

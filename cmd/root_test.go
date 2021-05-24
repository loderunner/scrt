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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
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

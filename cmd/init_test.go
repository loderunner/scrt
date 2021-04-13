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

//go:generate mockgen -destination mock_backend.go -package cmd "github.com/loderunner/scrt/backend" Backend

package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"

	"github.com/loderunner/scrt/backend"
)

func TestInitCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = func(name string) (backend.Backend, error) {
		return mockBackend, nil
	}

	viper.Set(configKeyPassword, "toto")
	viper.Set(configKeyStorage, "mock")
	viper.Set(configKeyLocation, "location")

	mockBackend.EXPECT().Exists().Return(false)
	mockBackend.EXPECT().Save(gomock.Any())

	args := []string{"path"}
	err := initCmd.RunE(initCmd, args)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInitOverWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = func(name string) (backend.Backend, error) {
		return mockBackend, nil
	}

	viper.Set(configKeyPassword, "toto")
	viper.Set(configKeyStorage, "mock")
	viper.Set(configKeyLocation, "location")

	mockBackend.EXPECT().Exists().Return(true)
	mockBackend.EXPECT().Save(gomock.Any())

	err := initCmd.Flags().Set("overwrite", "true")
	if err != nil {
		t.Fatal(err)
	}

	args := []string{"path"}
	err = initCmd.RunE(initCmd, args)
	if err != nil {
		t.Fatal(err)
	}
}

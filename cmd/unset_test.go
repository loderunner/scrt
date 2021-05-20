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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"

	"github.com/loderunner/scrt/backend"
	"github.com/loderunner/scrt/store"
)

func TestUnsetCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	password := "toto"

	viper.Reset()
	viper.Set(configKeyPassword, password)
	viper.Set(configKeyStorage, "mock")
	viper.Set(configKeyLocation, "location")

	s := store.NewStore()
	err := s.Set("hello", []byte("world"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := store.WriteStore([]byte(password), s)
	if err != nil {
		t.Fatal(err)
	}

	mockBackend.EXPECT().Exists().Return(true)
	mockBackend.EXPECT().Load().Return(data, nil)
	mockBackend.EXPECT().Save(gomock.Any())

	args := []string{"hello"}
	err = unsetCmd.Args(unsetCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = unsetCmd.RunE(unsetCmd, args)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnsetCmdNoValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	password := "toto"

	viper.Reset()
	viper.Set(configKeyPassword, password)
	viper.Set(configKeyStorage, "mock")
	viper.Set(configKeyLocation, "location")

	s := store.NewStore()
	data, err := store.WriteStore([]byte(password), s)
	if err != nil {
		t.Fatal(err)
	}

	mockBackend.EXPECT().Exists().Return(true)
	mockBackend.EXPECT().Load().Return(data, nil)
	mockBackend.EXPECT().Save(gomock.Any())

	args := []string{"hello"}
	err = unsetCmd.Args(unsetCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = unsetCmd.RunE(unsetCmd, args)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnsetCmdNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	viper.Reset()
	viper.Set(configKeyPassword, "toto")
	viper.Set(configKeyStorage, "mock")
	viper.Set(configKeyLocation, "location")

	mockBackend.EXPECT().Exists().Return(false)

	args := []string{"hello"}
	err := unsetCmd.Args(unsetCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = unsetCmd.RunE(unsetCmd, args)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUnsetCmdFailedLoad(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	viper.Reset()
	viper.Set(configKeyPassword, "toto")
	viper.Set(configKeyStorage, "mock")
	viper.Set(configKeyLocation, "location")

	mockBackend.EXPECT().Exists().Return(true)
	mockBackend.EXPECT().Load().Return(nil, fmt.Errorf("error"))

	args := []string{"hello"}
	err := unsetCmd.Args(unsetCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = unsetCmd.RunE(unsetCmd, args)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUnsetCmdFailedSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	password := "toto"

	viper.Reset()
	viper.Set(configKeyPassword, password)
	viper.Set(configKeyStorage, "mock")
	viper.Set(configKeyLocation, "location")

	s := store.NewStore()
	data, err := store.WriteStore([]byte(password), s)
	if err != nil {
		t.Fatal(err)
	}

	mockBackend.EXPECT().Exists().Return(true)
	mockBackend.EXPECT().Load().Return(data, nil)
	mockBackend.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("error"))

	args := []string{"hello"}
	err = unsetCmd.Args(unsetCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = unsetCmd.RunE(unsetCmd, args)
	if err == nil {
		t.Fatal("expected error")
	}
}

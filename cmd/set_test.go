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

func TestSetCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	password := "toto"

	viper.Reset()
	viper.Set(configKeyPassword, password)
	viper.Set(configKeyStorage, "mock")

	s := store.NewStore()
	data, err := store.WriteStore([]byte(password), s)
	if err != nil {
		t.Fatal(err)
	}

	mockBackend.EXPECT().Exists().Return(true, nil)
	mockBackend.EXPECT().Load().Return(data, nil)
	mockBackend.EXPECT().Save(gomock.Any())

	args := []string{"hello", "world"}
	err = setCmd.Args(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = setCmd.RunE(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetCmdStdin(t *testing.T) {
	// Hijack standard input
	hijack()
	defer restore()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	password := "toto"

	viper.Reset()
	viper.Set(configKeyPassword, password)
	viper.Set(configKeyStorage, "mock")

	s := store.NewStore()
	data, err := store.WriteStore([]byte(password), s)
	if err != nil {
		t.Fatal(err)
	}

	mockBackend.EXPECT().Exists().Return(true, nil)
	mockBackend.EXPECT().Load().Return(data, nil)
	mockBackend.EXPECT().Save(gomock.Any())

	_, err = hijackStdin.WriteString("world")
	if err != nil {
		t.Fatal(err)
	}
	hijackStdin.Close()

	args := []string{"hello"}
	err = setCmd.Args(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = setCmd.RunE(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetCmdNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	viper.Reset()
	viper.Set(configKeyPassword, "toto")
	viper.Set(configKeyStorage, "mock")

	mockBackend.EXPECT().Exists().Return(false, nil)

	args := []string{"hello", "world"}
	err := setCmd.Args(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = setCmd.RunE(setCmd, args)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSetCmdFailedLoad(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	viper.Reset()
	viper.Set(configKeyPassword, "toto")
	viper.Set(configKeyStorage, "mock")

	mockBackend.EXPECT().Exists().Return(true, nil)
	mockBackend.EXPECT().Load().Return(nil, fmt.Errorf("error"))

	args := []string{"hello", "world"}
	err := setCmd.Args(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = setCmd.RunE(setCmd, args)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSetCmdFailedInvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	viper.Reset()
	viper.Set(configKeyPassword, "toto")
	viper.Set(configKeyStorage, "mock")

	data := []byte("toto")

	mockBackend.EXPECT().Exists().Return(true, nil)
	mockBackend.EXPECT().Load().Return(data, nil)

	args := []string{"hello", "world"}
	err := setCmd.Args(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = setCmd.RunE(setCmd, args)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSetCmdFailedNoOverwrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	password := "toto"

	viper.Reset()
	viper.Set(configKeyPassword, password)
	viper.Set(configKeyStorage, "mock")

	s := store.NewStore()
	err := s.Set("hello", []byte("world"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := store.WriteStore([]byte(password), s)
	if err != nil {
		t.Fatal(err)
	}

	mockBackend.EXPECT().Exists().Return(true, nil)
	mockBackend.EXPECT().Load().Return(data, nil)

	args := []string{"hello", "world"}
	err = setCmd.Args(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = setCmd.RunE(setCmd, args)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSetCmdOverwrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	password := "toto"

	viper.Reset()
	viper.Set(configKeyPassword, password)
	viper.Set(configKeyStorage, "mock")

	s := store.NewStore()
	err := s.Set("hello", []byte("world"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := store.WriteStore([]byte(password), s)
	if err != nil {
		t.Fatal(err)
	}

	mockBackend.EXPECT().Exists().Return(true, nil)
	mockBackend.EXPECT().Load().Return(data, nil)
	mockBackend.EXPECT().Save(gomock.Any())

	err = setCmd.Flags().Set("overwrite", "true")
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"hello", "world"}
	err = setCmd.Args(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = setCmd.RunE(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetCmdFailedSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := NewMockBackend(ctrl)
	backend.Backends["mock"] = newMockFactory(mockBackend)

	password := "toto"

	viper.Reset()
	viper.Set(configKeyPassword, password)
	viper.Set(configKeyStorage, "mock")

	s := store.NewStore()
	data, err := store.WriteStore([]byte(password), s)
	if err != nil {
		t.Fatal(err)
	}

	mockBackend.EXPECT().Exists().Return(true, nil)
	mockBackend.EXPECT().Load().Return(data, nil)
	mockBackend.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("error"))

	args := []string{"hello", "world"}
	err = setCmd.Args(setCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = setCmd.RunE(setCmd, args)
	if err == nil {
		t.Fatal("expected error")
	}
}

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

//go:generate mockgen -destination mock_backend.go -package cmd "github.com/loderunner/scrt/backend" Backend

package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"

	"github.com/loderunner/scrt/backend"
	"github.com/loderunner/scrt/store"
)

func TestExportCmd(t *testing.T) {
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
	outputFile := ".env"

	testVal := []byte("world")
	err := s.Set("hello", testVal)
	if err != nil {
		t.Fatal(err)
	}

	testVal = []byte("world2")
	err = s.Set("hello2", testVal)
	if err != nil {
		t.Fatal(err)
	}

	data, err := store.WriteStore([]byte(password), s)
	if err != nil {
		t.Fatal(err)
	}

	mockBackend.EXPECT().ExistsContext(ctxMatcher).Return(true, nil)
	mockBackend.EXPECT().LoadContext(ctxMatcher).Return(data, nil)

	args := []string{}
	err = exportCmd.Args(getCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = exportCmd.RunE(getCmd, args)
	if err != nil {
		t.Fatal(err)
	}

	// check that the file was created
	_, err = os.Stat(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	// check that the file contains the expected content
	content, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "hello=world\nhello2=world2\n"

	err = os.Remove(outputFile)

	if string(content) != expected {
		if err != nil {
			t.Fatalf("expected '%s', got '%s'. Also, failed to delete out file '%s'", expected, string(content), outputFile)
		}
		t.Fatalf("expected '%s', got '%s'", expected, string(content))
	}

	// cleanup
	if err != nil {
		t.Fatal(err)
	}
}

func TestExportCmdWithFile(t *testing.T) {
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
	outputFile := "test.env"

	testVal := []byte("world")
	err := s.Set("hello", testVal)
	if err != nil {
		t.Fatal(err)
	}

	testVal = []byte("world2")
	err = s.Set("hello2", testVal)
	if err != nil {
		t.Fatal(err)
	}

	data, err := store.WriteStore([]byte(password), s)
	if err != nil {
		t.Fatal(err)
	}

	mockBackend.EXPECT().ExistsContext(ctxMatcher).Return(true, nil)
	mockBackend.EXPECT().LoadContext(ctxMatcher).Return(data, nil)

	args := []string{"--file", outputFile}
	err = exportCmd.Args(getCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	err = exportCmd.RunE(getCmd, args)
	if err != nil {
		t.Fatal(err)
	}

	// check that the file was created
	_, err = os.Stat(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	// check that the file contains the expected content
	content, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "hello=world\nhello2=world2\n"

	err = os.Remove(outputFile)

	if string(content) != expected {
		if err != nil {
			t.Fatalf("expected '%s', got '%s'. Also, failed to delete out file '%s'", expected, string(content), outputFile)
		}
		t.Fatalf("expected '%s', got '%s'", expected, string(content))
	}

	// cleanup
	if err != nil {
		t.Fatal(err)
	}
}

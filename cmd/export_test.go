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
	"sort"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"

	"github.com/loderunner/scrt/backend"
	"github.com/loderunner/scrt/store"
)

func TestExportCmdEnvFile(t *testing.T) {
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

	exportCmd.Flags().Set("output", outputFile)
	exportCmd.Flags().Set("format", "dotenv")

	err = exportCmd.RunE(exportCmd, []string{})
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

	expected := "hello=\"world\"\nhello2=\"world2\""

	err = os.Remove(outputFile)

	if !isSimilar(string(content), expected) {
		if err != nil {
			t.Fatalf("expected '%s', got '%s'. Also, failed to delete out file '%s'", expected, string(content), outputFile)
		}
		t.Fatalf("expected '%s', got '%s'", expected, string(content))
	}

	if err != nil {
		t.Fatal(err)
	}
}

func TestExportCmdJsonFile(t *testing.T) {
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
	outputFile := "test.json"

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

	exportCmd.Flags().Set("output", outputFile)
	exportCmd.Flags().Set("format", "json")

	err = exportCmd.RunE(exportCmd, []string{})
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

	expected := "{\"hello\":\"world\",\"hello2\":\"world2\"}"

	err = os.Remove(outputFile)

	if !isSimilar(string(content), expected) {
		if err != nil {
			t.Fatalf("expected '%s', got '%s'. Also, failed to delete out file '%s'", expected, string(content), outputFile)
		}
		t.Fatalf("expected '%s', got '%s'", expected, string(content))
	}

	if err != nil {
		t.Fatal(err)
	}
}

func TestExportCmdYamlFile(t *testing.T) {
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
	outputFile := "test.yaml"

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

	exportCmd.Flags().Set("output", outputFile)
	exportCmd.Flags().Set("format", "yaml")

	err = exportCmd.RunE(exportCmd, []string{})
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

	expected := "hello: world\nhello2: world2\n"

	err = os.Remove(outputFile)

	if !isSimilar(string(content), expected) {
		if err != nil {
			t.Fatalf("expected '%s', got '%s'. Also, failed to delete out file '%s'", expected, string(content), outputFile)
		}
		t.Fatalf("expected '%s', got '%s'", expected, string(content))
	}

	if err != nil {
		t.Fatal(err)
	}
}

func TestExportCmdDefault(t *testing.T) {
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

	// default
	exportCmd.Flags().Set("format", "dotenv")

	err = exportCmd.RunE(exportCmd, []string{})
	if err != nil {
		t.Fatal(err)
	}

	// hijackStdout is os.File
	os.Stdout.Close()
	content, err := ioutil.ReadAll(hijackStdout)

	if err != nil {
		t.Fatal(err)
	}

	expected := "hello=\"world\"\nhello2=\"world2\""

	if !isSimilar(string(content), expected) {
		t.Fatalf("expected '%s', got '%s'", expected, string(content))
	}

	// cleanup
	if err != nil {
		t.Fatal(err)
	}
}

func TestExportCmdInvalidFormat(t *testing.T) {
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

	exportCmd.Flags().Set("format", "fake_format")

	err := exportCmd.RunE(exportCmd, []string{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExportCmdMissingFormat(t *testing.T) {
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

	err := exportCmd.RunE(exportCmd, []string{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExportCmdViperSet(t *testing.T) {
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
	viper.Set(configKeyExportFormat, "yaml")
	viper.Set(configKeyExportOutput, "test.yaml")

	s := store.NewStore()
	outputFile := "test.yaml"

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

	err = exportCmd.RunE(exportCmd, []string{})
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

	expected := "hello: world\nhello2: world2\n"

	err = os.Remove(outputFile)

	if !isSimilar(string(content), expected) {
		if err != nil {
			t.Fatalf("expected '%s', got '%s'. Also, failed to delete out file '%s'", expected, string(content), outputFile)
		}
		t.Fatalf("expected '%s', got '%s'", expected, string(content))
	}

	if err != nil {
		t.Fatal(err)
	}
}

func isSimilar(a, b string) bool {
	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")

	trimLines(&aLines)
	trimLines(&bLines)

	// remove empty lines
	for i := 0; i < len(aLines); i++ {
		if aLines[i] == "" {
			aLines = append(aLines[:i], aLines[i+1:]...)
			i--
		}
	}

	for i := 0; i < len(bLines); i++ {
		if bLines[i] == "" {
			bLines = append(bLines[:i], bLines[i+1:]...)
			i--
		}
	}

	if len(aLines) != len(bLines) {
		return false
	}

	sort.Strings(aLines)
	sort.Strings(bLines)

	for i, line := range aLines {
		if line != bLines[i] {
			return false
		}
	}

	return true
}

func trimLines(lines *[]string) []string {
	for i, line := range *lines {
		(*lines)[i] = strings.TrimSpace(line)
	}
	return *lines
}

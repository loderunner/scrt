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
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/memory"
	"github.com/golang/mock/gomock"
	"github.com/loderunner/scrt/backend"
	"github.com/spf13/pflag"
)

var osStdin, osStdout *os.File
var hijackStdin, hijackStdout *os.File

// hijack stdin and stdout for testing.
func hijack() {
	var err error
	osStdin = os.Stdin
	os.Stdin, hijackStdin, err = os.Pipe()
	if err != nil {
		panic(err)
	}
	osStdout = os.Stdout
	hijackStdout, os.Stdout, err = os.Pipe()
	if err != nil {
		panic(err)
	}
}

// restore stdin and stdout. Usually deferred right after hijacking.
func restore() {
	os.Stdin = osStdin
	os.Stdout = osStdout
	hijackStdin.Close()
	hijackStdout.Close()
}

var ctxMatcher = gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem())

func TestMain(m *testing.M) {
	logger = &log.Logger{Handler: memory.New()}
	os.Exit(m.Run())
}

type mockFactory struct {
	b backend.Backend
}

func (f mockFactory) New(conf map[string]interface{}) (backend.Backend, error) {
	return f.NewContext(context.Background(), conf)
}

func (f mockFactory) NewContext(ctx context.Context, conf map[string]interface{}) (backend.Backend, error) {
	return f.b, nil
}

func (f mockFactory) Name() string {
	return "Mock"
}

func (f mockFactory) Description() string {
	return "mock storage"
}

func (f mockFactory) Flags() *pflag.FlagSet {
	return pflag.NewFlagSet("mock", pflag.ContinueOnError)
}

func newMockFactory(b backend.Backend) mockFactory {
	return mockFactory{b: b}
}

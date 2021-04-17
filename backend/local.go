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

package backend

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

var localFlagSet *pflag.FlagSet

func init() {
	localFlagSet = pflag.NewFlagSet("local", pflag.ContinueOnError)
}

type local struct {
	path string
	fs   afero.Fs
}

type localFactory struct{}

func (f localFactory) New(path string, conf map[string]interface{}) (Backend, error) {
	return newLocal(path, conf)
}

func (f localFactory) Flags() *pflag.FlagSet {
	return localFlagSet
}

func init() {
	Backends["local"] = localFactory{}
}

func newLocal(path string, conf map[string]interface{}) (Backend, error) {
	fs := afero.NewOsFs()
	_, err := fs.Stat(path)
	if err != nil && !errors.Is(err, afero.ErrFileNotFound) {
		return nil, fmt.Errorf("invalid location: \"path\"")
	}
	return local{path: path, fs: fs}, nil
}

func (l local) Valid() bool {
	_, err := os.Stat(l.path)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return true
	}
	return false
}

func (l local) Exists() bool {
	exists, _ := afero.Exists(l.fs, l.path)
	return exists
}

func (l local) Save(data []byte) error {
	return afero.WriteFile(l.fs, l.path, data, 0600)
}

func (l local) Load() ([]byte, error) {
	return afero.ReadFile(l.fs, l.path)
}

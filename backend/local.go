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

package backend

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

var localFlagSet *pflag.FlagSet

func init() {
	localFlagSet = pflag.NewFlagSet("local", pflag.ContinueOnError)
	localFlagSet.String("local-path", "", "path to the store in the local filesystem (required)")
}

type local struct {
	path string
	fs   afero.Fs
}

type localFactory struct{}

func (f localFactory) New(conf map[string]interface{}) (Backend, error) {
	return f.NewContext(context.Background(), conf)
}
func (f localFactory) NewContext(ctx context.Context, conf map[string]interface{}) (Backend, error) {
	return newLocal(ctx, conf)
}

func (f localFactory) Name() string {
	return "Local"
}

func (f localFactory) Description() string {
	return "store secrets to local filesystem"
}

func (f localFactory) Flags() *pflag.FlagSet {
	return localFlagSet
}

func init() {
	Backends["local"] = localFactory{}
}

func newLocal(ctx context.Context, conf map[string]interface{}) (Backend, error) {
	logger := getLogger(ctx)

	opt := readOpt("local", "path", conf)
	if opt == nil || opt == "" {
		return nil, fmt.Errorf("missing path")
	}
	path, ok := opt.(string)
	if !ok {
		return nil, fmt.Errorf("path is not a string: (%T)%s", path, path)
	}

	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	logger.WithField("path", path).Infof("using local path")

	fs := afero.NewOsFs()
	_, err = fs.Stat(path)
	if err != nil && !errors.Is(err, afero.ErrFileNotFound) {
		return nil, fmt.Errorf("invalid location: %s", path)
	}

	return local{path: path, fs: fs}, nil
}

func (l local) Exists() (bool, error) {
	return l.ExistsContext(context.Background())
}

func (l local) ExistsContext(ctx context.Context) (bool, error) {
	logger := getLogger(ctx)
	logger.WithField("path", l.path).Info("checking store existence")

	exists, err := afero.Exists(l.fs, l.path)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (l local) Save(data []byte) error {
	return l.SaveContext(context.Background(), data)
}

func (l local) SaveContext(ctx context.Context, data []byte) error {
	logger := getLogger(ctx)
	logger.WithField("path", l.path).Info("writing encrypted data to local storage")
	return afero.WriteFile(l.fs, l.path, data, 0600)
}

func (l local) Load() ([]byte, error) {
	return l.LoadContext(context.Background())
}

func (l local) LoadContext(ctx context.Context) ([]byte, error) {
	logger := getLogger(ctx)
	logger.WithField("path", l.path).Info("reading encrypted data from local storage")
	return afero.ReadFile(l.fs, l.path)
}

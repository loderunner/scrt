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
	"context"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/spf13/pflag"
)

// Backends associates backend type names to constructor functions
var Backends = map[string]Factory{}

// BackendNameList is an ordered list of backend names for listing in help
// message
var BackendNameList = []string{"local", "s3", "git"}

// Backend implements the common backend operations
type Backend interface {
	// Exists returns true if a store exists in the backend, false otherwise
	Exists() (bool, error)
	// Save persists encrypted data to the backend
	Save(data []byte) error
	// Load reads encrypted data from the backend
	Load() ([]byte, error)

	ExistsContext(ctx context.Context) (bool, error)
	SaveContext(ctx context.Context, data []byte) error
	LoadContext(ctx context.Context) ([]byte, error)
}

// Factory can instantiate a new Backend with New, and other static
// backend-related functions.
type Factory interface {
	// New returns an initialized backend from the given configuration
	New(conf map[string]interface{}) (Backend, error)
	NewContext(ctx context.Context, conf map[string]interface{}) (Backend, error)

	// Name returns a human-readable name for the backend
	Name() string
	// Description returns a short, human-readable description of the backend
	Description() string
	// Flags returns a pflag FlagSet containing the options related to the
	// backend
	Flags() *pflag.FlagSet
}

func readOpt(prefix, name string, conf map[string]interface{}) interface{} {
	var backendOpts map[string]interface{}
	l, ok := conf[prefix]
	if ok {
		backendOpts, _ = l.(map[string]interface{})
	}
	opt, ok := conf[prefix+"-"+name]
	if opt == "" || !ok {
		if backendOpts != nil {
			opt, ok = backendOpts[name]
		}
	}
	if opt == "" || !ok {
		return nil
	}
	return opt
}

func getLogger(ctx context.Context) log.Interface {
	logger := log.FromContext(ctx)
	if logger == log.Log {
		logger = &log.Logger{Handler: discard.Default}
	}
	return logger
}

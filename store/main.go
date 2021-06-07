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

package store

import (
	"fmt"

	"github.com/apex/log"
)

// Store defines a key-value storage in scrt.
type Store struct {
	data map[string][]byte
}

const saltLength = 16

// NewStore initializes a new Store.
func NewStore() Store {
	log.Info("creating new store")
	return Store{
		data: make(map[string][]byte),
	}
}

// Has returns true if a value is associated to key in the Store.
func (s Store) Has(key string) bool {
	log.WithField("key", key).Info("checking key existence")
	_, ok := s.data[key]
	return ok
}

// List returns all the keys is the Store
func (s Store) List() []string {
	log.Info("listing keys")
	keys := make([]string, len(s.data))
	i := 0
	for k := range s.data {
		keys[i] = k
		i++
	}
	return keys
}

// Get returns the value associated to key in the Store, or an error if none is
// associated.
func (s Store) Get(key string) ([]byte, error) {
	log.WithField("key", key).Info("retrieving value for key")
	if val, ok := s.data[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("no value for \"%s\"", key)
}

// Set associates the value to key in the Store, or an error if val is
// invalid.
func (s Store) Set(key string, val []byte) error {
	log.WithField("key", key).Info("setting value for key")
	if val == nil {
		return fmt.Errorf("cannot set value")
	}
	s.data[key] = val
	return nil
}

// Unset removes any value associated to key in the Store.
func (s Store) Unset(key string) {
	log.WithField("key", key).Info("unsetting value for key")
	delete(s.data, key)
}

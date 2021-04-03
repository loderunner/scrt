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

// Backends associates backend type names to constructor functions
var Backends = map[string]constructor{}

// Backend implements the common backend operations
type Backend interface {
	Exists() bool
	Save(data []byte) error
	Load() ([]byte, error)
}

type constructor func(name string) Backend

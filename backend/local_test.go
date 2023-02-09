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
	"os"
	"reflect"
	"testing"

	"github.com/loderunner/scrt/store"
	"github.com/spf13/afero"
)

func TestLocalExists(t *testing.T) {
	path := "/tmp/store.scrt"
	fs := afero.NewMemMapFs()
	f, err := fs.OpenFile(path, os.O_CREATE, 0600)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	b := local{path: "/tmp/nonexistent.scrt", fs: fs}
	exists, err := b.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Error("expected store not to exist")
	}

	b.path = path
	exists, err = b.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Error("expected store to exist")
	}
}

func TestLocalSaveLoad(t *testing.T) {
	path := "/tmp/store.scrt"
	fs := afero.NewMemMapFs()

	s := store.NewStore()
	data, _ := store.WriteStore([]byte("password"), s)

	b := local{path: path, fs: fs}
	err := b.Save(data)
	if err != nil {
		t.Fatal(err)
	}

	exists, err := b.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected store to exist")
	}

	got, err := b.Load()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(data, got) {
		t.Fatalf("expected %#v, got %#v", data, got)
	}
}

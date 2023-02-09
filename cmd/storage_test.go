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

package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestStorageCmd(t *testing.T) {
	hijack()
	defer restore()

	args := []string{"hello", "world"}
	err := storageCmd.Args(storageCmd, args)
	if err == nil {
		t.Fatal("expected error")
	}

	args = []string{}
	err = storageCmd.RunE(storageCmd, args)
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout.Close()
	data, err := ioutil.ReadAll(hijackStdout)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("no output")
	}
}

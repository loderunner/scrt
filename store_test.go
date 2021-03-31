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

package main

import (
	"crypto/rand"
	"math/big"
	"reflect"
	"testing"
)

func TestStore(t *testing.T) {
	expect := NewStore()

	passwordLength, err := rand.Int(rand.Reader, big.NewInt(256))
	if err != nil {
		t.Fatal(err)
	}
	password := make([]byte, passwordLength.Int64())
	n, err := rand.Read(password)
	if err != nil {
		t.Fatal(err)
	}
	if n != int(passwordLength.Int64()) {
		t.Fatalf("unexpected password length: %d", n)
	}

	data, err := WriteStore(password, expect)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ReadStore(password, data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expect.data, got.data) {
		t.Fatalf("expected %#v, got %#v", expect.data, got.data)
	}
}

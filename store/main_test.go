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
	"crypto/rand"
	"math/big"
	"reflect"
	"testing"
)

const testKey = "hello"
const testBadKey = "toto"

var testVal = []byte("world")
var testBinaryVal = make([]byte, 256)

func makePassword(t *testing.T) []byte {
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
	return password
}

func testWriteReadStore(t *testing.T, store Store) Store {
	password := makePassword(t)

	data, err := WriteStore(password, store)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ReadStore(password, data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(store.data, got.data) {
		t.Fatalf("expected %#v, got %#v", store.data, got.data)
	}

	return got
}

func testSetGet(t *testing.T, val []byte) {
	s := NewStore()

	err := s.Set(testKey, val)
	if err != nil {
		t.Fatal(err)
	}

	res, err := s.Get(testKey)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(res, val) {
		t.Fatalf("expected %#v, got %#v", val, res)
	}

	_, err = s.Get(testBadKey)
	if err == nil {
		t.Fatalf("expected s.Get(%#v) to return error", testBadKey)
	}
}

func TestEmptyStore(t *testing.T) {
	testWriteReadStore(t, NewStore())
}

func TestBasicStore(t *testing.T) {
	store := NewStore()

	err := store.Set(testKey, testVal)
	if err != nil {
		t.Fatal(err)
	}

	testWriteReadStore(t, store)
}

func TestTwice(t *testing.T) {
	store := NewStore()

	store = testWriteReadStore(t, store)
	testWriteReadStore(t, store)
}

func TestWriteNilData(t *testing.T) {
	store := Store{}
	password := makePassword(t)

	_, err := WriteStore(password, store)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestReadInvalidPassword(t *testing.T) {
	store := NewStore()

	store.data[testKey] = testVal

	password := makePassword(t)

	data, err := WriteStore(password, store)
	if err != nil {
		t.Fatal(err)
	}

	password = []byte("toto")
	_, err = ReadStore(password, data)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestReadInvalidData(t *testing.T) {
	// Random data for invalid store
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	password := makePassword(t)

	_, err = ReadStore(password, data)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestHasFail(t *testing.T) {
	s := NewStore()

	if s.Has(testKey) {
		t.Fatalf("expected s.Has(%#v) to return false", testKey)
	}
}

func TestSetGet(t *testing.T) {
	testSetGet(t, testVal)
}

func TestSetGetBinary(t *testing.T) {

	n, err := rand.Read(testBinaryVal)
	if n != 256 {
		t.Fatalf("expected 256 bytes, got %d", n)
	}
	if err != nil {
		t.Fatal(err)
	}
	testSetGet(t, testBinaryVal)
}

func TestSetHas(t *testing.T) {
	s := NewStore()

	err := s.Set(testKey, testVal)
	if err != nil {
		t.Fatal(err)
	}

	if !s.Has(testKey) {
		t.Fatalf("expected s.Has(%#v) to return true", testKey)
	}

	if s.Has(testBadKey) {
		t.Fatalf("expected s.Has(%#v) to return false", testKey)
	}
}

func TestGetFail(t *testing.T) {
	s := NewStore()

	_, err := s.Get(testKey)
	if err == nil {
		t.Fatalf("expected s.Get(%#v) to return error", testKey)
	}
}

func TestSetUnsetHasGet(t *testing.T) {
	s := NewStore()

	err := s.Set(testKey, testVal)
	if err != nil {
		t.Fatal(err)
	}
	s.Unset(testKey)

	if s.Has(testKey) {
		t.Fatalf("expected s.Has(%#v) to return false", testKey)
	}

	_, err = s.Get(testKey)
	if err == nil {
		t.Fatalf("expected s.Get(%#v) to return error", testKey)
	}
}

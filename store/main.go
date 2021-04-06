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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Store defines a key-value storage in scrt.
type Store struct {
	data map[string][]byte
	salt []byte
}

const saltLength = 16

// NewStore initializes a new Store.
func NewStore() Store {
	return Store{
		data: make(map[string][]byte),
	}
}

// ReadStore reads a scrt Store from raw data. ReadStore uses password to
// decrypt data and returns the Store, or an error if Store data could not be
// decrypted of parsed. A json.Unmarshal error can mean either that the wrong
// password was supplied, or that the Store is corrupted.
func ReadStore(password []byte, data []byte) (Store, error) {
	if len(data) < saltLength+aes.BlockSize {
		return Store{}, fmt.Errorf("invalid length")
	}
	salt := data[:saltLength]
	iv := data[saltLength : saltLength+aes.BlockSize]
	ciphertext := data[saltLength+aes.BlockSize:]

	key := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)

	block, err := aes.NewCipher(key)
	if err != nil {
		return Store{}, err
	}

	stream := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	store := Store{
		salt: salt,
	}
	err = json.Unmarshal(plaintext, &store.data)
	if err != nil {
		return Store{}, err
	}

	return store, nil
}

// WriteStore writes a Store as raw data to be saved. WriteStore uses password
// encrypt the Store and returns the encrypted data, or an error if the Store
// could not be encoded to JSON or could not be encrypted.
func WriteStore(password []byte, store Store) ([]byte, error) {
	if store.data == nil {
		return nil, fmt.Errorf("store data is nil")
	}

	plaintext, err := json.Marshal(store.data)
	if err != nil {
		return nil, err
	}

	salt := make([]byte, saltLength)
	n, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	if n != saltLength {
		return nil, fmt.Errorf("unexpected salt length: %d", n)
	}

	key := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)

	iv := make([]byte, aes.BlockSize)
	n, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}
	if n != aes.BlockSize {
		return nil, fmt.Errorf("unexpected IV length: %d", n)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	return append(salt, append(iv, ciphertext...)...), nil
}

// Has returns true if a value is associated to key in the Store.
func (s Store) Has(key string) bool {
	_, ok := s.data[key]
	return ok
}

// Get returns the value associated to key in the Store, or an error if none is
// associated.
func (s Store) Get(key string) ([]byte, error) {
	if val, ok := s.data[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("no value for \"%s\"", key)
}

// Set associates the value to key in the Store, or an error if val is
// invalid.
func (s Store) Set(key string, val []byte) error {
	if val == nil {
		return fmt.Errorf("cannot set value")
	}
	s.data[key] = val
	return nil
}

// Unset removes any value associated to key in the Store.
func (s Store) Unset(key string) {
	delete(s.data, key)
}

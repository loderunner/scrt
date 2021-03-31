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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type Store struct {
	data map[string][]byte
	salt []byte
}

const saltLength = 16

func NewStore() Store {
	return Store{
		data: make(map[string][]byte),
	}
}

func ReadStore(password []byte, data []byte) (Store, error) {
	salt := data[:saltLength]
	iv := data[saltLength : saltLength+aes.BlockSize]
	ciphertext := data[saltLength+aes.BlockSize:]

	key := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)

	block, err := aes.NewCipher(key)
	if err != nil {
		return Store{}, fmt.Errorf("failed to read store: %w", err)
	}

	stream := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	store := Store{
		salt: salt,
	}
	err = json.Unmarshal(plaintext, &store.data)
	if err != nil {
		return Store{}, fmt.Errorf("failed to read store: %w", err)
	}

	return store, nil
}

func WriteStore(password []byte, store Store) ([]byte, error) {
	plaintext, err := json.Marshal(store.data)
	if err != nil {
		return nil, fmt.Errorf("failed to write store: %w", err)
	}

	salt := make([]byte, saltLength)
	n, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to write store: %w", err)
	}
	if n != saltLength {
		return nil, fmt.Errorf("failed to write store: unexpected salt length: %d", n)
	}

	key := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)

	iv := make([]byte, aes.BlockSize)
	n, err = rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("failed to write store: %w", err)
	}
	if n != aes.BlockSize {
		return nil, fmt.Errorf("failed to write store: unexpected IV length: %d", n)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to write store: %w", err)
	}

	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	return append(salt, append(iv, ciphertext...)...), nil
}

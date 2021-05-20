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
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/loderunner/scrt/store"
)

type mockS3Client struct {
	s3iface.S3API
	data []byte
}

func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	if *input.Key != "/store.scrt" || m.data == nil {
		return nil, awserr.New(s3.ErrCodeNoSuchKey, "no such key", nil)
	}
	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader(m.data)),
	}, nil
}

func (m *mockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	var err error
	m.data, err = ioutil.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}
	return &s3.PutObjectOutput{}, nil
}

func TestS3Exists(t *testing.T) {
	b := s3Backend{bucket: "test-bucket", key: "/nonexistent.scrt", client: &mockS3Client{}}

	s := store.NewStore()
	data, _ := store.WriteStore([]byte("password"), s)

	err := b.Save(data)
	if err != nil {
		t.Fatal(err)
	}

	exists, err := b.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Error("expected store not to exist")
	}

	b.key = "/store.scrt"
	exists, err = b.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Error("expected store to exist")
	}
}

func TestS3SaveLoad(t *testing.T) {

	s := store.NewStore()
	data, _ := store.WriteStore([]byte("password"), s)

	b := s3Backend{bucket: "test-bucket", key: "/store.scrt", client: &mockS3Client{}}
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

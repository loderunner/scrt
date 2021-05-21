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
	"testing"
)

func TestS3Factory(t *testing.T) {
	f := s3Factory{}

	testGenericFactory(t, f)

	_, err := f.New("", map[string]interface{}{})
	if err == nil {
		t.Error("expected error")
	}

	_, err = f.New("", map[string]interface{}{"s3-bucket-name": "toto"})
	if err == nil {
		t.Errorf("expected error")
	}

	_, err = f.New("", map[string]interface{}{"s3-key": "toto"})
	if err == nil {
		t.Errorf("expected error")
	}

	_, err = f.New("", map[string]interface{}{
		"s3-bucket-name": "scrt-bucket",
		"s3-key":         "/store.scrt",
	})
	if err != nil {
		t.Error(err)
	}

	_, err = f.New("", map[string]interface{}{
		"s3": map[string]interface{}{
			"bucket-name": "scrt-bucket",
			"key":         "/store.scrt",
		},
	})
	if err != nil {
		t.Error(err)
	}

	_, err = f.New("", map[string]interface{}{
		"s3-bucket-name":  "scrt-bucket",
		"s3-key":          "/store.scrt",
		"s3-endpoint-url": "http://localhost:123456",
		"s3-region":       "fr-paris-75019",
	})
	if err != nil {
		t.Error(err)
	}

	_, err = f.New("", map[string]interface{}{
		"s3": map[string]interface{}{
			"bucket-name":  "scrt-bucket",
			"key":          "/store.scrt",
			"endpoint-url": "http://localhost:123456",
			"region":       "fr-paris-75019",
		},
	})
	if err != nil {
		t.Error(err)
	}

	_, err = f.New("", map[string]interface{}{
		"s3-bucket-name": []int{},
		"s3-key":         "/store.scrt",
	})
	if err == nil {
		t.Errorf("expected error")
	}

	_, err = f.New("", map[string]interface{}{
		"s3-bucket-name": "scrt-bucket",
		"s3-key":         []int{},
	})
	if err == nil {
		t.Errorf("expected error")
	}

	_, err = f.New("", map[string]interface{}{
		"s3-bucket-name":  "scrt-bucket",
		"s3-key":          "/store.scrt",
		"s3-endpoint-url": []int{},
	})
	if err == nil {
		t.Errorf("expected error")
	}

	_, err = f.New("", map[string]interface{}{
		"s3-bucket-name": "scrt-bucket",
		"s3-key":         "/store.scrt",
		"s3-region":      []int{},
	})
	if err == nil {
		t.Errorf("expected error")
	}
}

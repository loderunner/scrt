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

	_, err = f.New("toto", map[string]interface{}{})
	if err == nil {
		t.Errorf("expected error")
	}

	_, err = f.New("s3://test-bucket/store.scrt", map[string]interface{}{})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = f.New("s3://test-bucket/store.scrt", map[string]interface{}{
		"s3-endpoint-url": "http://localhost:123456",
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = f.New("s3://test-bucket/store.scrt", map[string]interface{}{
		"s3-endpoint-url": []int{},
	})
	if err == nil {
		t.Errorf("expected error")
	}
}

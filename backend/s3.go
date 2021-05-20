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
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/spf13/pflag"
)

var s3FlagSet *pflag.FlagSet

func init() {
	s3FlagSet = pflag.NewFlagSet("s3", pflag.ContinueOnError)
	s3FlagSet.String("s3-region", "", "region of the S3 storage")
	s3FlagSet.String("s3-endpoint-url", "", "override default S3 endpoint URL")
}

type s3Backend struct {
	bucket, key string
	client      s3iface.S3API
}

type s3Factory struct{}

func (f s3Factory) New(location string, conf map[string]interface{}) (Backend, error) {
	return newS3(location, conf)
}

func (f s3Factory) Name() string {
	return "S3"
}

func (f s3Factory) Description() string {
	return "store secrets to AWS S3 or S3-compatible object storage"
}

func (f s3Factory) Flags() *pflag.FlagSet {
	return s3FlagSet
}

func init() {
	Backends["s3"] = s3Factory{}
}

func newS3(location string, conf map[string]interface{}) (Backend, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	cfgs := []*aws.Config{}

	if url, ok := conf["s3-endpoint-url"]; ok {
		endpoint, ok := url.(string)
		if !ok {
			stringer, ok := url.(fmt.Stringer)
			if !ok {
				return nil, fmt.Errorf("S3 endpoint URL could not be converted to string: %v", url)
			}
			endpoint = stringer.String()
		}
		if endpoint != "" {
			cfg := aws.NewConfig().WithEndpoint(endpoint).WithS3ForcePathStyle(true)
			cfgs = append(cfgs, cfg)
		}
	}

	if region, ok := conf["s3-region"]; ok {
		r, ok := region.(string)
		if !ok {
			stringer, ok := region.(fmt.Stringer)
			if !ok {
				return nil, fmt.Errorf("S3 region could not be converted to string: %v", region)
			}
			r = stringer.String()
		}
		if r != "" {
			cfg := aws.NewConfig().WithRegion(r)
			cfgs = append(cfgs, cfg)
		}
	}

	client := s3.New(sess, cfgs...)

	s3URL, err := url.Parse(location)
	if err != nil ||
		s3URL.Scheme != "s3" ||
		s3URL.Path == "" ||
		s3URL.Path == "/" ||
		s3URL.User != nil ||
		s3URL.RawQuery != "" ||
		s3URL.Fragment != "" {
		return nil, fmt.Errorf("S3 URI is not written in the form: s3://mybucket/mykey")
	}

	return s3Backend{
		bucket: s3URL.Host,
		key:    s3URL.Path,
		client: client,
	}, nil
}

func (s s3Backend) Exists() (bool, error) {
	req := (&s3.GetObjectInput{}).
		SetBucket(s.bucket).
		SetKey(s.key)
	res, err := s.client.GetObject(req)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) &&
			(awsErr.Code() == s3.ErrCodeNoSuchBucket ||
				awsErr.Code() == s3.ErrCodeNoSuchKey) {
			return false, nil
		}
		return false, err
	}
	res.Body.Close()
	return true, nil
}

func (s s3Backend) Save(data []byte) error {
	req := (&s3.PutObjectInput{}).
		SetBucket(s.bucket).
		SetKey(s.key).
		SetBody(bytes.NewReader(data))
	_, err := s.client.PutObject(req)
	if err != nil {
		return err
	}
	return nil
}

func (s s3Backend) Load() ([]byte, error) {
	req := (&s3.GetObjectInput{}).
		SetBucket(s.bucket).
		SetKey(s.key)
	res, err := s.client.GetObject(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

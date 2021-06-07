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

	"github.com/apex/log"
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
	s3FlagSet.String("s3-bucket-name", "", "name of the S3 bucket (required)")
	s3FlagSet.String("s3-key", "", "path of the store object in the bucket (required)")
	s3FlagSet.String("s3-region", "", "region of the S3 storage")
	s3FlagSet.String("s3-endpoint-url", "", "override default S3 endpoint URL")
}

type s3Backend struct {
	bucket, key string
	client      s3iface.S3API
}

type s3Factory struct{}

func (f s3Factory) New(conf map[string]interface{}) (Backend, error) {
	return newS3(conf)
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

func newS3(conf map[string]interface{}) (Backend, error) {
	cfgs := []*aws.Config{}

	logEntry := log.NewEntry(log.Log.(*log.Logger))

	opt := readOpt("s3", "bucket-name", conf)
	if opt == nil || opt == "" {
		return nil, fmt.Errorf("missing bucket name")
	}
	bucket, ok := opt.(string)
	if !ok {
		return nil, fmt.Errorf("bucket name is not a string: (%T)%s", bucket, bucket)
	}
	logEntry = logEntry.WithField("bucket", bucket)

	opt = readOpt("s3", "key", conf)
	if opt == nil || opt == "" {
		return nil, fmt.Errorf("missing key")
	}
	key, ok := opt.(string)
	if !ok {
		return nil, fmt.Errorf("key is not a string: (%T)%s", key, key)
	}
	logEntry = logEntry.WithField("key", key)

	opt = readOpt("s3", "endpoint-url", conf)
	if opt != nil && opt != "" {
		endpoint, ok := opt.(string)
		if !ok {
			return nil, fmt.Errorf("S3 endpoint url is not a string")
		}
		cfg := aws.NewConfig().WithEndpoint(endpoint).WithS3ForcePathStyle(true)
		cfgs = append(cfgs, cfg)
		logEntry = logEntry.WithField("endpoint URL", endpoint)
	}

	opt = readOpt("s3", "region", conf)
	if opt != nil && opt != "" {
		region, ok := opt.(string)
		if !ok {
			return nil, fmt.Errorf("S3 region is not a string")
		}
		cfg := aws.NewConfig().WithRegion(region)
		cfgs = append(cfgs, cfg)
		logEntry = logEntry.WithField("region", region)
	}

	logEntry.Info("using S3 object")

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	client := s3.New(sess, cfgs...)

	return s3Backend{
		bucket: bucket,
		key:    key,
		client: client,
	}, nil
}

func (s s3Backend) Exists() (bool, error) {
	log.
		WithField("bucket", s.bucket).
		WithField("key", s.key).
		Info("checking store existence")

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
	log.
		WithField("bucket", s.bucket).
		WithField("key", s.key).
		Info("writing encrypted data to S3 storage")

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
	log.
		WithField("bucket", s.bucket).
		WithField("key", s.key).
		Info("reading encrypted data from S3 storage")

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

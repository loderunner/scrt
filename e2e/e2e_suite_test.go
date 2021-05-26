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

package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var envReplacer = strings.NewReplacer("-", "_", ".", "_")

func TestE2e(t *testing.T) {
	if os.Getenv("SCRT_TEST_E2E") != "y" {
		t.Log("Skipping e2e tests. Set environment variable SCRT_TEST_E2E=y to run e2e tests.")
		return
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "End-to-end tests")
}

var executablePath string
var tmpLocalDir string

var _ = BeforeSuite(func() {
	var err error

	executablePath, err = gexec.Build("github.com/loderunner/scrt")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	os.RemoveAll(tmpLocalDir)
})

var _ = Describe("scrt", func() {
	Context("for local backend", func() {
		var err error
		tmpLocalDir, err = os.MkdirTemp("", "scrt-e2e-*")
		Expect(err).NotTo(HaveOccurred())

		extraArgs := [4]map[string]string{
			{"local.path": filepath.Join(tmpLocalDir, "store-args.scrt")},
			{"local.path": filepath.Join(tmpLocalDir, "store-env.scrt")},
			{"local.path": filepath.Join(tmpLocalDir, "store-implicit-conf.scrt")},
			{"local.path": filepath.Join(tmpLocalDir, "store-explicit-conf.scrt")},
		}

		runTestsForStorage(
			"local",
			"toto",
			extraArgs,
		)
	})
	Context("for s3 backend", func() {
		extraArgs := [4]map[string]string{
			{
				"s3.bucket-name":  "test-bucket",
				"s3.key":          "/store-args.scrt",
				"s3.endpoint-url": os.Getenv("SCRT_TEST_E2E_S3_ENDPOINT_URL"),
				"s3.region":       os.Getenv("SCRT_TEST_E2E_S3_REGION"),
			},

			{
				"s3.bucket-name":  "test-bucket",
				"s3.key":          "/store-env.scrt",
				"s3.endpoint-url": os.Getenv("SCRT_TEST_E2E_S3_ENDPOINT_URL"),
				"s3.region":       os.Getenv("SCRT_TEST_E2E_S3_REGION"),
			},

			{
				"s3.bucket-name":  "test-bucket",
				"s3.key":          "/store-implicit-conf.scrt",
				"s3.endpoint-url": os.Getenv("SCRT_TEST_E2E_S3_ENDPOINT_URL"),
				"s3.region":       os.Getenv("SCRT_TEST_E2E_S3_REGION"),
			},
			{
				"s3.bucket-name":  "test-bucket",
				"s3.key":          "/store-explicit-conf.scrt",
				"s3.endpoint-url": os.Getenv("SCRT_TEST_E2E_S3_ENDPOINT_URL"),
				"s3.region":       os.Getenv("SCRT_TEST_E2E_S3_REGION"),
			},
		}

		runTestsForStorage(
			"s3",
			"toto",
			extraArgs,
		)
	})
	Context("for git backend", func() {
		extraArgs := [4]map[string]string{
			{
				"git.url":  os.Getenv("SCRT_TEST_E2E_GIT_REPOSITORY_URL"),
				"git.path": "store-args.scrt",
			},
			{
				"git.url":  os.Getenv("SCRT_TEST_E2E_GIT_REPOSITORY_URL"),
				"git.path": "store-env.scrt",
			},
			{
				"git.url":  os.Getenv("SCRT_TEST_E2E_GIT_REPOSITORY_URL"),
				"git.path": "store-implicit-conf.scrt",
			},
			{
				"git.url":  os.Getenv("SCRT_TEST_E2E_GIT_REPOSITORY_URL"),
				"git.path": "store-explicit-conf.scrt",
			},
		}

		runTestsForStorage(
			"git",
			"toto",
			extraArgs,
		)
	})
})

func execute(args []string, env []string) *gexec.Session {
	cmd := exec.Command(executablePath, args...)
	cmd.Env = append(os.Environ(), env...)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}

func runTestsForStorage(storage, password string, extraArgs [4]map[string]string) {
	Context("with args", func() {
		args := []string{
			"--storage=" + storage,
			"--password=" + password,
		}
		for k, v := range extraArgs[0] {
			args = append(args, "--"+strings.ReplaceAll(k, ".", "-")+"="+v)
		}
		runTests(args, []string{})
	})
	Context("with environment variables", func() {
		env := []string{
			"SCRT_STORAGE=" + storage,
			"SCRT_PASSWORD=" + password,
		}
		for k, v := range extraArgs[1] {
			k = "SCRT_" + strings.ToUpper(envReplacer.Replace(k))
			env = append(env, k+"="+v)
		}
		runTests([]string{}, env)
	})
	Context("with .scrt.yml implicit configuration file", func() {
		BeforeEach(func() {
			conf := map[string]interface{}{
				"password": password,
			}
			for k, v := range extraArgs[2] {
				i := strings.Index(k, ".")
				if i == -1 {
					conf[k] = v
				} else {
					// Handle nested conf 1-layer deep
					newK := k[:i]
					subK := k[i+1:]
					if subConf, ok := conf[newK]; ok {
						(subConf.(map[string]interface{}))[subK] = v
					} else {
						conf[newK] = map[string]interface{}{subK: v}
					}
				}
			}
			yamlData, err := yaml.Marshal(conf)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(".scrt.yml", yamlData, 0600)
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			os.Remove(".scrt.yml")
		})
		runTests([]string{}, []string{})
	})
	Context("with scrt.yml explicit configuration file", func() {
		BeforeEach(func() {
			conf := map[string]interface{}{
				"password": password,
			}
			for k, v := range extraArgs[3] {
				i := strings.Index(k, ".")
				if i == -1 {
					conf[k] = v
				} else {
					// Handle nested conf 1-layer deep
					newK := k[:i]
					subK := k[i+1:]
					if subConf, ok := conf[newK]; ok {
						(subConf.(map[string]interface{}))[subK] = v
					} else {
						conf[newK] = map[string]interface{}{subK: v}
					}
				}
			}
			yamlData, err := yaml.Marshal(conf)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile("scrt.yml", yamlData, 0600)
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			os.Remove("scrt.yml")
		})
		runTests([]string{"--config=scrt.yml"}, []string{})
	})
}

func runTests(args []string, env []string) {
	It("inits a new store", func() {
		session := execute(append(args, "init"), env)
		Eventually(session).Should(gexec.Exit(0))
		session.Wait()
	})

	It("fails when initing an existing store", func() {
		session := execute(append(args, "init"), env)
		Eventually(session).ShouldNot(gexec.Exit(0))
		session.Wait()
	})

	It("sets a new value", func() {
		session := execute(append(args, "set", "hello", "world"), env)
		Eventually(session).Should(gexec.Exit(0))
		session.Wait()
	})

	It("gets the value", func() {
		session := execute(append(args, "get", "hello"), env)
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say("world"))
		session.Wait()
	})

	It("fails to get a non-existing value", func() {
		session := execute(append(args, "get", "toto"), env)
		Eventually(session).ShouldNot(gexec.Exit(0))
		session.Wait()
	})

	It("fails to overwrite a value", func() {
		session := execute(append(args, "set", "hello", "world2"), env)
		Eventually(session).ShouldNot(gexec.Exit(0))
		session.Wait()
	})

	It("checks the value has not been overwritten", func() {
		session := execute(append(args, "get", "hello"), env)
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say("world"))
		session.Wait()
	})

	It("overwrites a value with --overwrite", func() {
		session := execute(append(args, "set", "--overwrite", "hello", "world2"), env)
		Eventually(session).Should(gexec.Exit(0))
		session.Wait()
	})

	It("checks the value has been overwritten", func() {
		session := execute(append(args, "get", "hello"), env)
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say("world2"))
		session.Wait()
	})

	It("lists the values", func() {
		session := execute(append(args, "list"), env)
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say("hello"))
		session.Wait()
	})

	It("unsets the value in the store", func() {
		session := execute(append(args, "unset", "hello"), env)
		Eventually(session).Should(gexec.Exit(0))
		session.Wait()
	})

	It("unsets a non-existing value in the store", func() {
		session := execute(append(args, "unset", "toto"), env)
		Eventually(session).Should(gexec.Exit(0))
		session.Wait()
	})

	It("fails to get a unset value", func() {
		session := execute(append(args, "get", "toto"), env)
		Eventually(session).ShouldNot(gexec.Exit(0))
		session.Wait()
	})
}

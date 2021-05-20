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

		tmpLocalLocations := [4]string{
			filepath.Join(tmpLocalDir, "store-args.scrt"),
			filepath.Join(tmpLocalDir, "store-env.scrt"),
			filepath.Join(tmpLocalDir, "store-implicit-conf.scrt"),
			filepath.Join(tmpLocalDir, "store-explicit-conf.scrt"),
		}

		runTestsForStorage(
			"local",
			"toto",
			tmpLocalLocations,
			nil,
		)
	})
	Context("for s3 backend", func() {
		s3Locations := [4]string{
			"s3://test-bucket/store-args.scrt",
			"s3://test-bucket/store-env.scrt",
			"s3://test-bucket/store-implicit-conf.scrt",
			"s3://test-bucket/store-explicit-conf.scrt",
		}

		runTestsForStorage(
			"s3",
			"toto",
			s3Locations,
			map[string]string{
				"s3-endpoint-url": os.Getenv("SCRT_TEST_E2E_S3_ENDPOINT_URL"),
				"s3-region":       os.Getenv("SCRT_TEST_E2E_S3_REGION"),
			},
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

func runTestsForStorage(storage, password string, locations [4]string, extraArgs map[string]string) {
	Context("with args", func() {
		args := []string{
			"--storage=" + storage,
			"--password=" + password,
			"--location=" + locations[0],
		}
		for k, v := range extraArgs {
			args = append(args, "--"+k+"="+v)
		}
		runTests(args, []string{})
	})
	Context("with environment variables", func() {
		env := []string{
			"SCRT_STORAGE=" + storage,
			"SCRT_PASSWORD=" + password,
			"SCRT_LOCATION=" + locations[1],
		}
		for k, v := range extraArgs {
			env = append(env, "SCRT_"+strings.ToUpper(k)+"="+v)
		}
		runTests([]string{}, env)
	})
	Context("with .scrt.yml implicit configuration file", func() {
		BeforeEach(func() {
			conf := map[string]string{
				"storage":  storage,
				"password": password,
				"location": locations[2],
			}
			for k, v := range extraArgs {
				conf[k] = v
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
			conf := map[string]string{
				"storage":  storage,
				"password": password,
				"location": locations[3],
			}
			for k, v := range extraArgs {
				conf[k] = v
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

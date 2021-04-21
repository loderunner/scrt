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
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/loderunner/scrt/cmd"
)

var (
	version = "dev"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			handlePanic(r)
		}
	}()

	cmd.RootCmd.Version = version
	err := cmd.RootCmd.Execute()
	if err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	var posixErr syscall.Errno
	if errors.As(err, &posixErr) {
		os.Exit(int(posixErr))
	}
	os.Exit(-1)
}

func handlePanic(err interface{}) {
	fmt.Println(err)
	os.Exit(-1)
}

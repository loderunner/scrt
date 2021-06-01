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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/kevinburke/ssh_config"
	"github.com/spf13/pflag"
)

var gitFlagSet *pflag.FlagSet

const defaultCommitMessage = "update secrets"

func init() {
	gitFlagSet = pflag.NewFlagSet("git", pflag.ContinueOnError)
	gitFlagSet.String("git-url", "", "URL of the git repository")
	gitFlagSet.String("git-path", "", "path of the store in the repository")
	gitFlagSet.String("git-branch", "", "branch to checkout, commit and push to")
}

type gitBackend struct {
	url    string
	path   string
	branch string
	repo   *git.Repository
	fs     billy.Filesystem
}

type gitFactory struct{}

func (f gitFactory) New(conf map[string]interface{}) (Backend, error) {
	return newGit(conf)
}

func (f gitFactory) Name() string {
	return "Git"
}

func (f gitFactory) Description() string {
	return "store secrets to a git repository"
}

func (f gitFactory) Flags() *pflag.FlagSet {
	return gitFlagSet
}

func init() {
	Backends["git"] = gitFactory{}
}

func newGit(conf map[string]interface{}) (Backend, error) {
	opt := readOpt("git", "url", conf)
	if opt == nil || opt == "" {
		return nil, fmt.Errorf("missing repository URL")
	}
	url, ok := opt.(string)
	if !ok {
		return nil, fmt.Errorf("repository URL is not a string: (%T)%s", opt, opt)
	}

	opt = readOpt("git", "path", conf)
	if opt == nil || opt == "" {
		return nil, fmt.Errorf("missing path")
	}
	path, ok := opt.(string)
	if !ok {
		return nil, fmt.Errorf("path is not a string: (%T)%s", opt, opt)
	}

	var branch string
	opt = readOpt("git", "branch", conf)
	if opt != nil && opt != "" {
		branch, ok = opt.(string)
		if !ok {
			return nil, fmt.Errorf("branch is not a string: (%T)%s", opt, opt)
		}
	}

	g := gitBackend{
		url:    url,
		path:   path,
		branch: branch,
	}

	err := g.clone()
	// If the repo is empty, init a new repo
	if err == transport.ErrEmptyRemoteRepository {
		err = g.init()
	}
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g gitBackend) Exists() (bool, error) {
	_, err := g.fs.Stat(g.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (g gitBackend) Save(data []byte) error {
	f, err := g.fs.OpenFile(g.path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0700)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := f.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("wrote %d bytes, expected %d", n, len(data))
	}
	err = f.Close()
	if err != nil {
		return err
	}
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(g.path)
	if err != nil {
		return err
	}
	authorCommitter := &object.Signature{
		Name:  "scrt",
		Email: "scrt@scrt.run",
		When:  time.Now(),
	}
	_, err = w.Commit(defaultCommitMessage, &git.CommitOptions{
		Author:    authorCommitter,
		Committer: authorCommitter,
	})
	if err != nil {
		return err
	}
	err = g.repo.Push(&git.PushOptions{
		RemoteName: git.DefaultRemoteName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g gitBackend) Load() ([]byte, error) {
	f, err := g.fs.OpenFile(g.path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (g *gitBackend) clone() error {
	auths, err := buildAuths(g.url)
	if err != nil {
		return err
	}
	storage := memory.NewStorage()
	g.fs = memfs.New()
	if len(auths) > 0 {
		for _, auth := range auths {
			g.repo, err = git.Clone(
				storage,
				g.fs,
				&git.CloneOptions{
					URL:           g.url,
					ReferenceName: plumbing.NewBranchReferenceName(g.branch),
					Depth:         1,
					Auth:          auth,
				},
			)
			if err == nil {
				return nil
			}
		}
	} else {
		g.repo, err = git.Clone(
			storage,
			g.fs,
			&git.CloneOptions{
				URL:           g.url,
				ReferenceName: plumbing.NewBranchReferenceName(g.branch),
				Depth:         1,
			},
		)
	}
	if err != nil {
		g.fs = nil
	}
	return err
}

func (g *gitBackend) init() error {
	var err error

	storage := memory.NewStorage()
	g.fs = memfs.New()
	g.repo, err = git.Init(storage, g.fs)
	if err != nil {
		return err
	}

	_, err = g.repo.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{g.url},
	})
	if err != nil {
		return err
	}

	// Set default branch name, if not configured
	if g.branch == "" {
		g.branch = "main"
		gitConfig, err := g.repo.ConfigScoped(config.SystemScope)
		if err == nil {
			n := gitConfig.Init.DefaultBranch
			if n != "" {
				g.branch = n
			}
		}
	}

	ref := plumbing.NewSymbolicReference(
		plumbing.HEAD,
		plumbing.NewBranchReferenceName(g.branch),
	)
	err = storage.SetReference(ref)
	if err != nil {
		return err
	}

	return nil
}

func buildAuths(url string) ([]ssh.AuthMethod, error) {
	e, err := transport.NewEndpoint(url)
	if err != nil {
		return nil, err
	}
	if e.Protocol != "ssh" {
		return nil, nil
	}

	sshConfig := ssh_config.DefaultUserSettings
	if sshConfig == nil {
		defaultAuth, err := ssh.DefaultAuthBuilder(e.User)
		if err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{defaultAuth}, nil
	}

	auths := make([]ssh.AuthMethod, 0, 2)

	identitiesOnly := sshConfig.Get(e.Host, "IdentitiesOnly")
	if identitiesOnly != "yes" {
		sshAgentAuth, err := ssh.NewSSHAgentAuth(e.User)
		if err == nil {
			auths = append(auths, sshAgentAuth)
		}
	}

	idFiles := sshConfig.GetAll(e.Host, "IdentityFile")
	for _, idFile := range idFiles {
		publicKeyAuth, err := ssh.NewPublicKeysFromFile(e.User, idFile, "")
		if err == nil {
			auths = append(auths, publicKeyAuth)
		}
	}

	if len(auths) > 0 {
		return auths, nil
	}

	return nil, fmt.Errorf("no valid authentication method")
}

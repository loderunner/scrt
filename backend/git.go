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
	"context"
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
	gitFlagSet.String("git-url", "", "URL of the git repository (required)")
	gitFlagSet.String("git-path", "", "path of the store in the repository (required)")
	gitFlagSet.String("git-branch", "", "branch to checkout, commit and push to on updates")
	gitFlagSet.String("git-checkout", "", "tree-ish revision to checkout, e.g. commit or tag")
	gitFlagSet.String("git-message", "", "commit message when updating the store")
}

type gitBackend struct {
	path    string
	message string
	repo    *git.Repository
	fs      billy.Filesystem
}

type gitFactory struct{}

func (f gitFactory) New(conf map[string]interface{}) (Backend, error) {
	return f.NewContext(context.Background(), conf)
}

func (f gitFactory) NewContext(ctx context.Context, conf map[string]interface{}) (Backend, error) {
	return newGit(ctx, conf)
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

func newGit(ctx context.Context, conf map[string]interface{}) (Backend, error) {
	logger := getLogger(ctx)

	opt := readOpt("git", "url", conf)
	if opt == nil || opt == "" {
		return nil, fmt.Errorf("missing repository URL")
	}
	url, ok := opt.(string)
	if !ok {
		return nil, fmt.Errorf("repository URL is not a string: (%T)%s", opt, opt)
	}
	logger = logger.WithField("url", url)

	opt = readOpt("git", "path", conf)
	if opt == nil || opt == "" {
		return nil, fmt.Errorf("missing path")
	}
	path, ok := opt.(string)
	if !ok {
		return nil, fmt.Errorf("path is not a string: (%T)%s", opt, opt)
	}
	logger = logger.WithField("path", path)

	var branch string
	opt = readOpt("git", "branch", conf)
	if opt != nil && opt != "" {
		branch, ok = opt.(string)
		if !ok {
			return nil, fmt.Errorf("branch is not a string: (%T)%s", opt, opt)
		}
		logger = logger.WithField("branch", branch)
	}

	var checkout string
	opt = readOpt("git", "checkout", conf)
	if opt != nil && opt != "" {
		checkout, ok = opt.(string)
		if !ok {
			return nil, fmt.Errorf("checkout is not a string: (%T)%s", opt, opt)
		}
		logger = logger.WithField("checkout", checkout)
	}

	message := defaultCommitMessage
	opt = readOpt("git", "message", conf)
	if opt != nil && opt != "" {
		message, ok = opt.(string)
		if !ok {
			return nil, fmt.Errorf("message is not a string: (%T)%s", opt, opt)
		}
	}

	logger.Info("using git repository")

	g := gitBackend{
		path:    path,
		message: message,
	}

	err := g.clone(ctx, url, branch)
	// If the repo is empty, init a new repo
	if errors.Is(err, transport.ErrEmptyRemoteRepository) {
		logger.Info("repository is empty")
		err = g.init(ctx, url, branch)
	}
	if err != nil {
		return nil, err
	}

	if checkout != "" {
		err = g.checkout(ctx, checkout)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (g gitBackend) Exists() (bool, error) {
	return g.ExistsContext(context.Background())
}
func (g gitBackend) ExistsContext(ctx context.Context) (bool, error) {
	logger := getLogger(ctx)

	logger.
		WithField("path", g.path).
		Info("checking store existence")

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
	return g.SaveContext(context.Background(), data)
}
func (g gitBackend) SaveContext(ctx context.Context, data []byte) error {
	logger := getLogger(ctx)

	logger = logger.WithField("path", g.path)

	logger.Info("opening file in git repository")
	f, err := g.fs.OpenFile(g.path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0700)
	if err != nil {
		return err
	}
	defer f.Close()

	logger.Info("writing encrypted data to git repository")
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

	logger.Info("staging file in worktree")
	_, err = w.Add(g.path)
	if err != nil {
		return err
	}

	gitConfig, err := g.repo.ConfigScoped(config.SystemScope)
	if err != nil {
		return err
	}
	user := gitConfig.User
	authorCommitter := &object.Signature{
		Name:  user.Name,
		Email: user.Email,
		When:  time.Now(),
	}

	logger.
		WithField("committer", fmt.Sprintf("%s <%s>", authorCommitter.Name, authorCommitter.Email)).
		Infof("committing changes to git repository: \"%s\"", g.message)
	_, err = w.Commit(
		g.message,
		&git.CommitOptions{
			Author:    authorCommitter,
			Committer: authorCommitter,
		},
	)
	if err != nil {
		return err
	}

	logger.Info("pushing changes to git remote")
	err = g.repo.Push(&git.PushOptions{
		RemoteName: git.DefaultRemoteName,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g gitBackend) Load() ([]byte, error) {
	return g.LoadContext(context.Background())
}
func (g gitBackend) LoadContext(ctx context.Context) ([]byte, error) {
	logger := getLogger(ctx)

	logger.
		WithField("path", g.path).
		Info("reading encrypted data from git repository")

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

func (g *gitBackend) clone(ctx context.Context, url, branch string) error {
	logger := getLogger(ctx)
	auths, err := buildAuths(ctx, url)
	if err != nil {
		return err
	}

	logger = logger.WithField("url", url)
	storage := memory.NewStorage()
	g.fs = memfs.New()
	var referenceName plumbing.ReferenceName
	if branch != "" {
		logger = logger.WithField("branch", branch)
		referenceName = plumbing.NewBranchReferenceName(branch)
	}

	logger.Info("cloning git repository in memory")
	if len(auths) > 0 {
		for _, auth := range auths {
			g.repo, err = git.Clone(
				storage,
				g.fs,
				&git.CloneOptions{
					URL:           url,
					ReferenceName: referenceName,
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
				URL:           url,
				ReferenceName: referenceName,
			},
		)
	}

	return err
}

func (g *gitBackend) init(ctx context.Context, url, branch string) error {
	logger := getLogger(ctx)
	var err error

	logger.Info("initializing git repository in memory")
	storage := memory.NewStorage()
	g.fs = memfs.New()
	g.repo, err = git.Init(storage, g.fs)
	if err != nil {
		return err
	}

	logger.WithField("url", url).Info("adding git remote")
	_, err = g.repo.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{url},
	})
	if err != nil {
		return err
	}

	// Set default branch name, if not configured
	if branch == "" {
		branch = "main"
		gitConfig, err := g.repo.ConfigScoped(config.SystemScope)
		if err == nil {
			n := gitConfig.Init.DefaultBranch
			if n != "" {
				branch = n
			}
		}
	}
	logger.WithField("branch", branch).Info("using git branch")

	ref := plumbing.NewSymbolicReference(
		plumbing.HEAD,
		plumbing.NewBranchReferenceName(branch),
	)
	err = storage.SetReference(ref)
	if err != nil {
		return err
	}

	return nil
}

func buildAuths(ctx context.Context, url string) ([]ssh.AuthMethod, error) {
	logger := getLogger(ctx)
	e, err := transport.NewEndpoint(url)
	if err != nil {
		return nil, err
	}
	if e.Protocol != "ssh" {
		return nil, nil
	}

	logger.Info("configuring SSH authentication methods")

	sshConfig := ssh_config.DefaultUserSettings
	if sshConfig == nil {
		defaultAuth, err := ssh.DefaultAuthBuilder(e.User)
		if err != nil {
			return nil, err
		}
		logger.Info("SSH config not found, using default authentication")
		return []ssh.AuthMethod{defaultAuth}, nil
	}

	auths := make([]ssh.AuthMethod, 0, 2)

	identitiesOnly := sshConfig.Get(e.Host, "IdentitiesOnly")
	if identitiesOnly != "yes" {
		logger.Info("using SSH agent authentication (if available)")
		sshAgentAuth, err := ssh.NewSSHAgentAuth(e.User)
		if err == nil {
			auths = append(auths, sshAgentAuth)
		}
	} else {
		logger.Info("using identity files only")
	}

	idFiles := sshConfig.GetAll(e.Host, "IdentityFile")
	for _, idFile := range idFiles {
		logger := logger.WithField("identity_file", idFile)
		publicKeyAuth, err := ssh.NewPublicKeysFromFile(e.User, idFile, "")
		if err == nil {
			logger.Info("identity file found")
			auths = append(auths, publicKeyAuth)
		} else {
			logger.Info("identity file not found")
		}
	}

	if len(auths) > 0 {
		return auths, nil
	}

	return nil, fmt.Errorf("no valid authentication method")
}

func (g *gitBackend) checkout(ctx context.Context, checkout string) error {
	logger := getLogger(ctx)

	logger.
		WithField("revision", checkout).
		Info("resolving revision to commit hash")

	hash, err := g.repo.ResolveRevision(plumbing.Revision(checkout))
	if err != nil {
		return err
	}

	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	logger.WithField("hash", hash).Info("checking out hash on a detached HEAD")
	return w.Checkout(&git.CheckoutOptions{
		Hash: *hash,
	})
}

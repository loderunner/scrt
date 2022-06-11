---
sidebarDepth: 0
---

# Git

Use the `git` storage type to create and access a store in a git repository. `scrt` will clone the repository in memory, checkout the given branch (or the default branch if no branch is given), read the store in the file at the given path, and will commit and push any modifications to the remote.

### Options

**`--git-url`** (required): a git-compatible repository URL. Most git-compatible URLs and protocols can be used. See [`git clone` documentation](https://git-scm.com/docs/git-clone#_git_urls) to learn more.

**`--git-path`** (required): the path to the store file inside the the git repository, relative to the repository root. A repository can contain multiple scrt stores, at different paths.

**`--git-branch`:** the name of the branch to checkout after cloning (or initializing). If no branch is given, the default branch from the remote will be used, or `main` if a new repository is initialized.

**`--git-checkout`:** a git revision to checkout. If this option is specified, the revision will be checked out in a ["detached HEAD"](https://git-scm.com/docs/git-checkout#_detached_head) and pushing will not work; making updates (`init`, `set` or `unset`) will be impossible.

**`--git-message`:** the message of the git commit. A default message will be used if this is not set.

### Example

```shell
scrt init --storage=git \
          --password=p4ssw0rd \
          --git-url=git@github.com:githubuser/secrets.git \
          --git-path=store.scrt
```

::: tip
`scrt` will initialize a new repo if none can be cloned.
:::

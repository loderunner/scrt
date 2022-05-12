![scrt logo](scrt-logo.png)

[![Workflow status](https://img.shields.io/github/workflow/status/loderunner/scrt/Build%20scrt?style=flat-square)](https://github.com/loderunner/scrt/actions/workflows/build.yml) [![Coverage Status](https://img.shields.io/coveralls/github/loderunner/scrt?style=flat-square)](https://coveralls.io/github/loderunner/scrt?branch=main) [![Go reference](https://img.shields.io/static/v1?label=%E2%80%8B&message=reference&color=00add8&logo=go&style=flat-square)](https://pkg.go.dev/github.com/loderunner/scrt)

`scrt` is a command-line secret manager for developers, sysadmins and devops. `scrt` aims to provide command-line users with a secure way of storing and retrieving secrets, while retaining control of the storage.

---

- [Features](#features)
- [Installation](#installation)
  - [Download binary release](#download-binary-release)
  - [`apt` (Debian/Ubuntu)](#apt-debianubuntu)
  - [`yum` (RHEL/CentOS)](#yum-rhelcentos)
  - [Homebrew (macOS)](#homebrew-macos)
  - [`go get`](#go-get)
  - [Build from source](#build-from-source)
- [Example](#example)
  - [Initialization](#initialization)
  - [Configuration](#configuration)
  - [Using the store](#using-the-store)
- [Usage](#usage)
  - [Global options](#global-options)
  - [Listing storage types](#listing-storage-types)
  - [Initializing a store](#initializing-a-store)
  - [Storing a secret](#storing-a-secret)
  - [Retrieving a secret](#retrieving-a-secret)
  - [Listing all secrets](#listing-all-secrets)
  - [Deleting a secret](#deleting-a-secret)
- [Configuration](#configuration-1)
  - [Configuration file](#configuration-file)
  - [Environment variables](#environment-variables)
- [Storage types](#storage-types)
  - [Local](#local)
  - [S3](#s3)
  - [Git](#git)
- [FAQ](#faq)
- [License](#license)

# Features

- Stateless command-line tool for Linux/Windows/Darwin
- All cryptography happens in the client, on your computer: no passwords, keys or plaintext data over the wire, no key management included
- Key/value interface: `get`/`set`/`unset`
- Configuration from command-line, configuration file or environment variables (no unexpected defaults!)
- Multiple backend choices:
  - Local filesystem
  - S3 (or S3-compatible object storage)
  - Git repository
  - More to come...

# Installation

## Download binary release

Download the latest finary release for your platform from the [releases page](https://github.com/loderunner/scrt/releases). Decompress the archive to the desired location. E.g.

```shell
tar xzvf scrt_1.2.3_linux_x86_64.tar.gz
sudo cp scrt_1.2.3_linux_x86_64/scrt /usr/local/bin/scrt
```

## `apt` (Debian/Ubuntu)

Configure the apt repository:

```shell
echo "deb [signed-by=/usr/share/keyrings/scrt-archive-keyring.gpg] https://apt.scrt.run /" | sudo tee /etc/apt/sources.list.d/scrt.list
curl "https://apt.scrt.run/key.gpg" | gpg --dearmor | sudo tee /usr/share/keyrings/scrt-archive-keyring.gpg > /dev/null
```

Install the binary package:

```shell
sudo apt update
sudo apt install scrt
```

## `yum` (RHEL/CentOS)

Configure the yum repository, in `/etc/yum.repos.d/scrt.repo`:

```ini
[scrt]
name=scrt
baseurl=https://yum.scrt.run
repo_gpgcheck=1
gpgcheck=1
enabled=1
gpgkey=https://yum.scrt.run/key.gpg
sslverify=1
metadata_expire=300
```

Install the binary package

```shell
sudo yum update
sudo yum install scrt
```

## Homebrew (macOS)

Configure the Homebrew tap:

```shell
brew tap loderunner/scrt
```

Install the binary package:

```
brew install scrt
```

## `go get`

Use `go get` to download and build the latest version:

```shell
go get github.com/loderunner/scrt
```

`scrt` will be available in the binaries directory of your GOPATH. Add it to your path, and run `scrt`.

## Build from source

Clone the repository and use `go build` to build a binary (requires go >= 1.16):

```shell
git clone https://github.com/loderunner/scrt.git
cd scrt
go build .
```

The built executable will be located at `scrt` at the root of the repository.

# Example

## Initialization

Initialize a new store, with `scrt init`.

```shell
scrt init --storage=local \
          --password=p4ssw0rd \
          --local-path=~/.scrt/store.scrt
# store initialized
```

This will create an empty store, in a `store.scrt` file located in `.scrt` inside your home directory. The file is encrypted using a secret key derived from the given password.

The content of the file is unreadable:

```
00000000  e0 97 af ea 86 f7 6a f0  82 06 47 8f fc 54 47 8e  |......j...G..TG.|
00000010  89 f9 ca f4 00 98 24 f3  85 1e bd 85 e5 c1 66 43  |......$.......fC|
00000020  d8 5d 47 2b 99 b1 99 fa  2c 07 0a ec 8c 11        |.]G+....,.....|
```

## Configuration

Set your configuration in environment variables, so you don't have to type them out each time you run a command.

```shell
export SCRT_STORAGE=local
export SCRT_PASSWORD=p4ssw0rd
export SCRT_LOCAL_PATH=~/.scrt/store.scrt
```

## Using the store

Set and retrieve a value for a key using `scrt set` and `scrt get`.

```shell
scrt set hello 'World!'
scrt get hello
# Output: World!
```

The content of the file is still unreadable, but now contains your value:

```
00000000  1d cc 02 68 c0 e5 d4 a4  9d 8f ff 14 0c 3b 73 71  |...h.........;sq|
00000010  47 54 2a 78 d8 87 63 fd  29 dc b4 e4 72 c7 0e 57  |GT*x..c.)...r..W|
00000020  be 04 ba e9 7d 36 6d e1  64 47 e2 e2 c0 fb 83 30  |....}6m.dG.....0|
00000030  51 9e ad cf 15 d8 7e 35  77 1c 0c f1 70 be cb 91  |Q.....~5w...p...|
```

# Usage

Use `scrt --help` to output a full help message.

```
A secret manager for the command-line

Usage:
  scrt [command]

Available Commands:
  init        Initialize a new store
  set         Associate a key to a value in a store
  get         Retrieve the value associated to key from a store
  list        List all the keys in a store
  unset       Remove the value associated to key in a store
  storage     List storage types and options
  help        Help about any command

Flags:
  -c, --config string     configuration file
  -h, --help              help for scrt
  -p, --password string   master password to unlock the store
      --storage string    storage type
  -v, --version           version for scrt

Use "scrt [command] --help" for more information about a command.
```

## Global options

**`-c`**, **`--config`:** Path to a YAML [Configuration file](#configuration-file)

**`--storage`:** storage type, see [Storage types](#storage-types) for details.

**`-p`**, **`--password`:** password to the store. The argument will be used to derive a key, to decrypt and encrypt the data in the store.

> In the following examples, these options will be sometimes omitted, as they can be [configured](#configuration) using a configuration file or environment variables.

## Listing storage types

```
scrt storage
```

List all available storage types and options

## Initializing a store

```
scrt init [flags]
```

Initialize a new store. If an item is already present at the given location, the initialization will fail unless the `--overwrite` option is set.

### Example

Create a store in a `store.scrt` file in the local filesystem, in the current working directory, using the password `"p4ssw0rd"`.

```shell
scrt init --storage=local --password=p4ssw0rd --local-path=./store.scrt
```

### Options

**`--overwrite`:** when this flag is set, `scrt` will overwrite the item at the given location, if it exists, instead of returning an error. If no item exists at the location, `--overwrite` has no effect.

## Storing a secret

```
scrt set [flags] key [value]
```

Associate a value to a key in the store. If `value` is omitted from the command
line, it will be read from standard input.

If a value is already set for `key`, the command will fail unless the `--overwrite` option is set.

### Example

Associate `Hello World` to the key `greeting` in the store, using implicit store configuration (configuration file or environment variables).

```shell
scrt set greeting "Hello World"
```

### Options

**`--overwrite`:** when this flag is set, `scrt` will overwrite the value for `key` in the store, if it exists, instead of returning an error. If no value is associated to `key`, `--overwrite` has no effect.

## Retrieving a secret

```
scrt get [flags] key
```

Retrieve the value associated to the key in the store, if it exists. Returns an error if no value is associated to the key.

### Example

Retrieve the value associated to the key `greeting` in the store, using implicit store configuration (configuration file or environment variables).

```shell
scrt get greeting
# Output: Hello World
```

## Listing all secrets

```
scrt list
```

List the keys of all the secrets in the store.

### Example

List the keys of all the secrets in the store, using implicit store configuration (configuration file or environment variables).

```shell
scrt list
# Output: greeting
```

## Deleting a secret

```
scrt unset [flags] key
```

Disassociate the value associated to a key in the store. If no value is associated to the key, does nothing.

### Example

Remove the value associated to the key. After this command, no value will be associated to the key `greeting` in the store.

```shell
scrt unset greeting
```

# Configuration

Repeating the global options every time the `scrt` command is invoked can be verbose. Also, some options–like the store password–shouldn't be used on the command line on a shared computer, to avoid security issues.

To prevent this, `scrt` can be configured with a configuration file or using environment variables.

`scrt` uses the following precedence order. Each item takes precedence over the item below it:

- flags
- environment variables
- configuration file

> Configuration options can be considered to be chosen from "most explicit" (flags) to "least explicit" (configuration file).

## Configuration file

The `scrt` configuration file is a YAML file with the configuration options as keys.

Example:

```yaml
storage: local
password: p4ssw0rd
local:
  path: ~/.scrt/store.scrt
```

If the `--config` option is given to the command-line, `scrt` will try to load the configuration from a file at the given path. Otherwise, it looks for any file named `.scrt`, `.scrt.yml` or `.scrt.yaml` in the current working directory, then recursively in the parent directory up to the root of the filesystem. If such a file is found, its values are loaded as configuration.

This can be useful in configuring the location of a store for a project, by adding a `.scrt` file at the root of the project repository. `scrt` can then be used in CI and other DevOps tools.

> :warning: Don't add the password to a configuration file in a public git repository! :warning:

Storage type (`storage`) can be ignored in a configuration file. `scrt` will read the configuration under the key for the storage type (e.g. `local:`). _Defining configurations for multiple storage types in a single file will result in undefined behavior._

## Environment variables

Each global option has an environment variable counterpart. Environment variables use the same name as the configuration option, in uppercase letters, prefixed with `SCRT_`.

- `storage` ⇒ `SCRT_STORAGE`
- `password` ⇒ `SCRT_PASSWORD`
- `local-path` ⇒ `SCRT_LOCAL_PATH`

To configure a default store on your system, add the following to your `.bashrc` file (if using `bash`):

```bash
export SCRT_STORAGE=local
export SCRT_PASSWORD=p4ssw0rd
export SCRT_LOCAL_PATH=~/.scrt/store.scrt
```

> Refer to your shell interpreter's documentation to set environment variables if you don't use `bash` (`zsh`, `dash`, `tcsh`, etc.)

# Storage types

```
Local:
  local       store secrets to local filesystem
Flags:
      --local-path string   path to the store in the local filesystem
                            (required)

S3:
  s3          store secrets to AWS S3 or S3-compatible object storage
Flags:
      --s3-bucket-name string    name of the S3 bucket (required)
      --s3-endpoint-url string   override default S3 endpoint URL
      --s3-key string            path of the store object in the bucket
                                 (required)
      --s3-region string         region of the S3 storage

Git:
  git         store secrets to a git repository
Flags:
      --git-branch string     branch to checkout, commit and push to on updates
      --git-checkout string   tree-ish revision to checkout, e.g. commit or tag
      --git-message string    commit message when updating the store
      --git-path string       path of the store in the repository (required)
      --git-url string        URL of the git repository (required)
```

`scrt` supports various storage backends, independent of the secrets engine. Each storage type has a name, and configuration options vary according to the chosen type.

Storage types may support additional options. See the documentation below for details.

## Local

Use the `local` storage type to create and access a store on your local filesystem.

Example:

```shell
scrt init --storage=local --password=p4ssw0rd --local-path=/tmp/store.scrt
```

### Options

**`--local-path`** (required): the path to the store file on the local filesystem.

## S3

Use the `s3` storage type to create and access a store using [AWS S3](https://aws.amazon.com/s3/) or any compatible object storage (such as [MinIO](https://min.io/)).

Example:

```shell
scrt init --storage=s3 \
          --password=p4ssw0rd \
          --s3-bucket-name=scrt-bucket \
          --s3-key=/store.scrt
```

> `scrt` uses your [AWS configuration (config files, environment variables)](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) if it can be found.

### Options

**`--s3-bucket-name`** (required): the name of the bucket to save to store to

**`--s3-key`** (required): the key to the store object

**`--s3-region`:** set the region for the S3 bucket

**`--s3-endpoint-url`:** when using an S3-compatible object storage other than AWS, `scrt` requires the URL of the S3 API endpoint.

## Git

Use the `git` storage type to create and access a store in a git repository. `scrt` will clone the repository in memory, checkout the given branch (or the default branch if no branch is given), read the store in the file at the given path, and will commit and push any modifications to the remote.

Example:

```shell
scrt init --storage=git \
          --password=p4ssw0rd \
          --git-url=git@github.com:githubuser/secrets.git \
          --git-path=store.scrt
```

> `scrt` will initialize a new repo if none can be cloned.

### Options

**`--git-url`** (required): a git-compatible repository URL. Most git-compatible URLs and protocols can be used. See [`git clone` documentation](https://git-scm.com/docs/git-clone#_git_urls) to learn more.

**`--git-path`** (required): the path to the store file inside the the git repository, relative to the repository root. A repository can contain multiple scrt stores, at different paths.

**`--git-branch`:** the name of the branch to checkout after cloning (or initializing). If no branch is given, the default branch from the remote will be used, or `main` if a new repository is initialized.

**`--git-checkout`:** a git revision to checkout. If specified, the revision will be checked out in a ["detached HEAD"](https://git-scm.com/docs/git-checkout#_detached_head) and pushing will not work; making updates (`init`, `set` or `unset`) will be impossible.

**`--git-message`:** the message of the git commit. A default message will be used if this is not set.

# FAQ

### How do you pronounce `scrt`?

Nobody knows. It's either "secret" without the e's; or "skrrt" like a [Migos](https://genius.com/artists/Migos) ad-lib.

### What is the cryptography behind `scrt`?

`scrt` relies on the industry-standard [AES](https://csrc.nist.gov/publications/detail/fips/197/final) symmetric encryption algorithm with 256-bit keys, with GCM [mode of operation](https://csrc.nist.gov/publications/detail/sp/800-38a/final) (AES-256-GCM, in OpenSSL parlance).

The encryption keys are derived from the password using the [Argon2id](https://www.password-hashing.net/#argon2) key derivation function. A new random salt is used every time the store is written to, preventing reuse of existing cryptographic keys.

### Does `scrt` store my keys? Should I be worried about my secrets being intercepted?

`scrt` does not save keys in the store, nor does it transfer any plaintext over the wire. All decryption and encryption happens on your computer while the program is running. This is the only way to provide full privacy and zero-trust security.

The downside to this is that the entire store must be loaded into memory, possibly downloading it through the network, decrypted, and possibly reencrypted (on a mutating operation like `set` or `unset`) every time you run `scrt`. If the size of your store becomes an issue, there are workarounds like splitting your store into multiple stores, or downloading the entire store to the local filesystem before using it.

### I lost my password, how can I recover my secrets?

I've got some good news and some bad news.

The bad news: you're doomed. Your secrets are encrypted with a key that can only be derived from your password. `scrt` does not store or manage keys. There is no way to recover your secrets without your password.

The good news: you probably won't lose your password again.

# License

[Apache 2.0](https://choosealicense.com/licenses/apache-2.0/)

[![Workflow status](https://img.shields.io/github/workflow/status/loderunner/scrt/scrt?style=flat-square)](https://github.com/loderunner/scrt/actions/workflows/workflow.yml) [![Go reference](https://img.shields.io/static/v1?label=%E2%80%8B&message=reference&color=00add8&logo=go&style=flat-square)](https://pkg.go.dev/github.com/loderunner/scrt) [![Coverage Status](https://img.shields.io/coveralls/github/loderunner/scrt?style=flat-square)](https://coveralls.io/github/loderunner/scrt?branch=main)

# scrt

`scrt` is a command-line secret manager for developers, sysadmins and devops. `scrt` aims to provide command-line users with a secure way of storing and retrieving secrets, while retaining control of the storage.

- [Features](#features)
- [Installation](#installation)
  - [`go get`](#go-get)
  - [Build from source](#build-from-source)
- [Example](#example)
  - [Initialization](#initialization)
  - [Configuration](#configuration)
  - [Using the store](#using-the-store)
- [Usage](#usage)
  - [Global options](#global-options)
  - [Initializing a store](#initializing-a-store)
  - [Storing a secret](#storing-a-secret)
  - [Retrieving a secret](#retrieving-a-secret)
  - [Deleting a secret](#deleting-a-secret)
- [Configuration](#configuration-1)
  - [Configuration file](#configuration-file)
  - [Environment variables](#environment-variables)
- [Storage types](#storage-types)
  - [Local](#local)
  - [S3](#s3)
- [FAQ](#faq)
- [License](#license)

## Features

- Stateless command-line tool for Linux/Windows/Darwin
- All cryptography happens in the client: no passwords, keys or plaintext data over the wire, no key management included
- Key/value interface: `get`/`set`/`unset`
- Configuration from command-line, configuration file or environment variables (no unexpected defaults!)
- Multiple backend choices:
  - Local filesystem
  - S3 (or S3-compatible object storage)
  - More to come...

## Installation

### `go get`

Use `go get` to download and build the latest version:

```sh
go get github.com/loderunner/scrt
```

`scrt` will be available in the binaries directory of your GOPATH. Add it to your path, and run `scrt`.

### Build from source

Clone the repository and use `go build` to build a binary (requires go >= 1.16):

```sh
git clone https://github.com/loderunner/scrt.git
cd scrt
go build .
```

The built executable will be located at `scrt` at the root of the repository.

## Example

### Initialization

Initialize a new store, with `scrt init`.

```
$ scrt init --storage=local \
            --location=~/.scrt/store.scrt \
            --password=p4ssw0rd
local store initialized at ~/.scrt/store.scrt
```

This will create an empty store, in a `store.scrt` file located in `.scrt` inside your home directory. The file is encrypted using a secret key derived from the given password.

The content of the file is unreadable:

```
00000000  e0 97 af ea 86 f7 6a f0  82 06 47 8f fc 54 47 8e  |......j...G..TG.|
00000010  89 f9 ca f4 00 98 24 f3  85 1e bd 85 e5 c1 66 43  |......$.......fC|
00000020  d8 5d 47 2b 99 b1 99 fa  2c 07 0a ec 8c 11        |.]G+....,.....|
```

### Configuration

Set your configuration in environment variables, so you don't have to type them out each time you run a command.

```
$ export SCRT_STORAGE=local
$ export SCRT_LOCATION=~/.scrt/store.scrt
$ export SCRT_PASSWORD=p4ssw0rd
```

### Using the store

Set and retrieve a value for a key using `scrt set` and `scrt get`.

```
$ scrt set hello 'World!'
$ scrt get hello
World!
```

The content of the file is still unreadable, but now contains your value:

```
00000000  1d cc 02 68 c0 e5 d4 a4  9d 8f ff 14 0c 3b 73 71  |...h.........;sq|
00000010  47 54 2a 78 d8 87 63 fd  29 dc b4 e4 72 c7 0e 57  |GT*x..c.)...r..W|
00000020  be 04 ba e9 7d 36 6d e1  64 47 e2 e2 c0 fb 83 30  |....}6m.dG.....0|
00000030  51 9e ad cf 15 d8 7e 35  77 1c 0c f1 70 be cb 91  |Q.....~5w...p...|
```

## Usage

Use `scrt --help` to output a full help message.

```
A secret manager for the command-line

Usage:
  scrt [command]

Available Commands:
  get         Retrieve the value associated to key from a store
  help        Help about any command
  init        Initialize a new store
  set         Associate a key to a value in a store
  unset       Remove the value associated to key in a store

Flags:
  -c, --config string     configuration file
  -h, --help              help for scrt
      --location string   store location
  -p, --password string   master password to unlock the store
      --storage string    storage type
  -v, --version           version for scrt

Use "scrt [command] --help" for more information about a command.
```

### Global options

**`-c`**, **`--config`:** Path to a YAML [Configuration file](#configuration-file)

**`--storage`:** storage type, see [Storage types](#storage-types) for details.

**`--location`:** location of the store, [storage](#storage-types)-dependent.

**`-p`**, **`--password`:** password to the store. The argument will be used to derive a key, to decrypt and encrypt the data in the store.

> In the following examples, these options will be sometimes omitted, as they can be [configured](#configuration) using a configuration file or environment variables.

### Initializing a store

```
scrt init [flags]
```

Initialize a new store at the given location. If an item is already present at the location, the initialization will fail unless the `--overwrite` option is set.

#### Example

Create a store in a `store.scrt` file in the local filesystem, in the current working directory, using the password `"p4ssw0rd"`.

```sh
scrt init --storage=local --location=./store.scrt --password=p4ssw0rd
```

#### Options

**`--overwrite`:** when this flag is set, `scrt` will overwrite the item at the given location, if it exists, instead of returning an error. If no item exists at the given location, `--overwrite` has no effect.

### Storing a secret

```
scrt set [flags] key [value]
```

Associate a value to a key in the store. If `value` is omitted from the command
line, it will be read from standard input.

If a value is already set for `key`, the command will fail unless the `--overwrite` option is set.

#### Example

Associate `Hello World` to the key `greeting` in the store, using implicit store configuration (configuration file or environment variables).

```sh
scrt set greeting "Hello World"
```

#### Options

**`--overwrite`:** when this flag is set, `scrt` will overwrite the value for `key` in the store, if it exists, instead of returning an error. If no value is associated to `key`, `--overwrite` has no effect.

### Retrieving a secret

```
scrt get [flags] key
```

Retrieve the value associated to the key in the store, if it exists. Returns an error if no value is associated to the key.

#### Example

Retrieve the value associated to the key `greeting` in the store, using implicit store configuration (configuration file or environment variables).

```sh
scrt get greeting
# Output: Hello World
```

### Deleting a secret

```
scrt unset [flags] key
```

Disassociate the value associated to a key in the store. If no value is associated to the key, does nothing.

#### Example

Remove the value associated to the key. After this command, no value will be associated to the key `greeting` in the store.

```sh
scrt unset greeting
```

## Configuration

Repeating the global options every time the `scrt` command is invoked can be verbose. Also, some options–like the store password–shouldn't be used on the command line on a shared computer, to avoid security issues.

To prevent this, `scrt` can be configured with a configuration file or using environment variables.

`scrt` uses the following precedence order. Each item takes precedence over the item below it:

- flags
- environment variables
- configuration file

> Configuration options can be considered to be chosen from "most explicit" (flags) to "least explicit" (configuration file).

### Configuration file

The `scrt` configuration file is a YAML file with the configuration options as keys.

Example:

```yaml
storage: local
location: ~/.scrt/store.scrt
password: p4ssw0rd
```

If the `--config` option is set, `scrt` will try to load the configuration from a file at the given path. otherwise, it looks for any file named `.scrt`, `.scrt.yml` or `.scrt.yaml` in the current working directory, then recursively in the parent directory up to the root of the filesystem. If such a file is found, its values are loaded as configuration.

This can be useful in configuring the location of a store for a project, by adding a `.scrt` file at the root of the project repository. `scrt` can then be used in CI and other DevOps tools.

> :warning: Don't add the password to a public git repository! :warning:

### Environment variables

Each global option has an environment variable counterpart. Environment variables use the same name as the configuration option, in uppercase letters, prefixed with `SCRT_`.

- `storage` ⇒ `SCRT_STORAGE`
- `location` ⇒ `SCRT_LOCATION`
- `password` ⇒ `SCRT_PASSWORD`

To configure a default store on your system, add the following to your `.bashrc` file (if using `bash`):

```bash
export SCRT_STORAGE=local
export SCRT_LOCATION=~/.scrt/store.scrt
export SCRT_PASSWORD=p4ssw0rd
```

> Refer to your shell interpreter's documentation to set environment variables if you don't use `bash` (`zsh`, `dash`, `tcsh`, etc.)

## Storage types

`scrt` supports various storage backends, independent of the secrets engine. Each storage type has a name, and location strings vary according to the chosen type.

Storage types may support additional options. See the documentation below for details.

### Local

Use the `local` storage type to create and access a store on your local filesystem.

Example:

```sh
scrt init --storage=local --location=/tmp/store.scrt --password=p4ssw0rd
```

#### Location

With the `local` backend, the location string is a path on the filesystem.

### S3

Use the `s3` storage type to create and access a store using [AWS S3](https://aws.amazon.com/s3/) or any compatible object storage (such as [MinIO](https://min.io/)).

Example:

```sh
scrt init --storage=s3 --location=s3://scrt-bucket/store.scrt --password=p4ssw0rd
```

> `scrt` uses your [AWS configuration (config files, environment variables)](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) if it can be found.

#### Location

With the `s3` backend, the location string is an S3 URI of the form `s3://mybucket/mykey`.

#### Extra options

**`s3-endpoint-url`:** When using an S3-compatible object storage other than AWS, `scrt` requires the URL of the S3 API endpoint. Can be configured in the configuration file, or with the `SCRT_S3_ENDPOINT_URL` environment variable.

## FAQ

#### How do you pronounce `scrt`?

Nobody knows. It's either "secret" without the e's; or "skrrt" like a [Migos](https://genius.com/artists/Migos) ad-lib.

#### What is the cryptography behind `scrt`?

`scrt` relies on the industry-standard [AES](https://csrc.nist.gov/publications/detail/fips/197/final) symmetric encryption algorithm with 256-bit keys, with GCM [mode of operation](https://csrc.nist.gov/publications/detail/sp/800-38a/final) (AES-256-GCM, in OpenSSL parlance).

The encryption keys are derived from the password using the [Argon2id](https://www.password-hashing.net/#argon2) key derivation function. A new random salt is used every time the store is written to, preventing reuse of existing cryptographic keys.

#### Does `scrt` store my keys? Should I be worried about my secrets being intercepted?

`scrt` does not save keys in the store, nor does it transfer any plaintext over the wire. All decryption and encryption happens on your computer while the program is running. This is the only way to provide full privacy and zero-trust security.

The downside to this is that the entire store must be loaded into memory, possibly downloading it through the network, decrypted, and possibly reencrypted (on a mutating operation like `set` or `unset`) every time you run `scrt`. If the size of your store becomes an issue, there are workarounds like splitting your store into multiple stores, or downloading the entire store to the local filesystem before using it.

#### I lost my password, how can I recover my secrets?

I've got some good news and some bad news.

The bad news: you're doomed. Your secrets are encrypted with a key that can only be derived from your password. `scrt` does not store or manage keys. There is no way to recover your secrets without your password.

The good news: you probably won't lose your password again.

## License

[Apache 2.0](https://choosealicense.com/licenses/apache-2.0/)

# Getting Started

This examples in this section will help you set up and start using scrt with a store in a local file.

## Initialization

Initialize a new store, with `scrt init`

```shell
scrt init --storage=local \
          --password=p4ssw0rd \
          --local-path=store.scrt
```

This will create an empty store, in a `store.scrt` file. The file is encrypted using a secret key derived from the given password. (The password in these examples is very weak. In a production setting, do not use such a simple password. Follow the [NIST recommandations](https://auth0.com/blog/dont-pass-on-the-new-nist-password-guidelines/) for good password creation.)

::: warning
The password is the key to all your secrets. If you lose your password, there is no way to recover your secrets.
:::

## Configuration

Set your configuration in environment variables, so you don't have to type them out each time you run a command.

```shell
export SCRT_STORAGE=local
export SCRT_PASSWORD=p4ssw0rd
export SCRT_LOCAL_PATH=store.scrt
```

In the following examples, we assume the environment variables have been set. See the [Configuration reference](/guide/configuration.md) for advanced configuration options.

## Using the store

### Setting and retrieving

Set and retrieve a value for a key using `scrt set` and `scrt get`

```shell
scrt set hello 'World!'
scrt get hello

# Output: World!
```

### Adding more secrets

Add another secret by using `scrt set` for another key

```shell
scrt set barbes 'rochechouart'
scrt get barbes

# Output: rochechouart
```

### Updating

Update a secret with `scrt set` for an existing key

```shell
scrt set hello 'Brooklyn'
scrt get hello

# Output: Brooklyn
```

### Deleting

Delete a secret with `scrt unset`

```shell
scrt unset hello
scrt get hello

# Error: no value for key: "hello"
```

## Getting help

Get more information about scrt commands with

```shell
scrt --help
```

### Storage options

scrt supports several storage backends. Find out more about the storage options with

```shell
scrt storage
```

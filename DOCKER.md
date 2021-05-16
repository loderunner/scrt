# What is scrt?

[`scrt`](https://github.com/loderunner/scrt) is a secret manager for the command-line. Use `scrt` to encrypt and decrypt secrets from a variety of storage backends.

## How to use this image?

See the [documentation](https://github.com/loderunner/scrt/README.md) for details on how to run `scrt`.

### Basic usage

```sh
docker run loderuner/scrt --help
```

### Configuration

`scrt` can be configured in a container with the command-line argument, with environment variables or using a `.scrt.yml` configuration file.

### Command-line options

Run the docker image and pass the options as to the executable.

```sh
docker run loderunner/scrt --storage=s3 \
                           --location=s3://scrt-bucket/store.scrt \
                           --password=p4ssw0rd \
                           init
```

#### Environment variables

Use the [environment variables](https://github.com/loderunner/scrt/README.md#environment-variables) to configure `scrt`:

```sh
docker run --env SCRT_STORAGE=s3 \
           --env SCRT_LOCATION=s3://scrt-bucket/store.scrt \
           --env SCRT_PASSWORD=p4ssw0rd \
           init
```

#### Configuration file

Create a `.scrt.yml` configuration file on the host:

```yaml
storage: s3
location: s3://scrt-bucker/store.scrt
password: p4ssw0rd
```

then mount the configuration file at the root of the container filesystem:

```shell
docker run -v .scrt.yml:/.scrt.yml loderunner/scrt init
```

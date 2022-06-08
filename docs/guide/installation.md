---
description: Find out how to install scrt secret manager on Linux, macOS or Windows. Download a binary release, install from an apt or yum repository, install with Homebrew, or build from source.
---

# Installation

## Download binary release

Download the latest binary release for your platform from the [releases page](https://github.com/loderunner/scrt/releases). Decompress the archive to the desired location. E.g.

```shell
tar xzvf scrt_0.3.3_linux_x86_64.tar.gz
sudo cp scrt_0.3.3_linux_x86_64/scrt /usr/local/bin/scrt
```

## apt (Debian/Ubuntu)

Configure the apt repository:

```shell
echo "deb [signed-by=/usr/share/keyrings/scrt-archive-keyring.gpg] https://apt.scrt.run /" \
  | sudo tee /etc/apt/sources.list.d/scrt.list
curl "https://apt.scrt.run/key.gpg" \
  | gpg --dearmor \
  | sudo tee /usr/share/keyrings/scrt-archive-keyring.gpg > /dev/null
```

Install the binary package:

```shell
sudo apt update
sudo apt install scrt
```

## yum (RHEL/CentOS/Rocky Linux)

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

## go get

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

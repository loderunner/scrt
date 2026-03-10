---
description: Contributing to scrt — how to build from source, run tests, and verify changes locally.
---

# Contributing <!-- omit in toc -->

- [Prerequisites](#prerequisites)
- [Clone a fork of the Repository](#clone-a-fork-of-the-repository)
- [Build](#build)
- [Run Tests](#run-tests)
- [Install Locally](#install-locally)
- [Verify Your Changes](#verify-your-changes)

## Prerequisites

- Go >= 1.18
- Git
- SSH key configured for GitHub (if testing git storage)

## Clone a fork of the Repository

```bash
git clone https://github.com/loderunner/scrt.git

```bash
cd scrt
```

## Build

**Compile without producing a binary** (useful for quick syntax/type checks):

```bash
go build ./...
```

**Build the binary:**

```bash
go build -o scrt .
```

The binary will be at `./scrt` in the repository root.

## Run Tests

```bash
go test ./...
```

## Install Locally

**Option A — Replace the Homebrew binary (macOS):**

```bash
cp scrt /opt/homebrew/bin/scrt
```

**Option B — Install to `$GOPATH/bin`:**

```bash
go install .
```

Make sure `$GOPATH/bin` is in your `$PATH`.

## Verify Your Changes

**Quick smoke test with local storage:**

```bash
./scrt init --storage=local --password=test --local-path=test.scrt
```

```bash
./scrt set --storage=local --password=test --local-path=test.scrt mykey "myvalue"
```

```bash
./scrt get --storage=local --password=test --local-path=test.scrt mykey
```

**Test with git storage:**

```bash
./scrt init --storage=git \
    --password=test \
    --git-url=git@github.com:<user>/<repo>.git \
    --git-path=store.scrt \
    --verbose
```

**Clean up test artifacts:**

```bash
rm -f test.scrt
```

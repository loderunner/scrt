---
home: true
title: Home
heroText: A command-line secret manager for developers, sysadmins, and devops
tagline: null
description: A command-line secret manager for developers, sysadmins, and devops
heroImage: /images/hero.png
actions:
  - text: Get Started
    link: /guide/installation.html
    type: primary
  - text: GitHub
    link: //github.com/loderunner/scrt
    type: secondary
features:
  - title: No Background Process
    image: /images/no-background.svg
    details: No service or daemon to install. Read and write secrets with a single command, on Linux/Windows/macOS.
  - title: Intuitive Interface
    image: /images/command-line.svg
    details: 'Straightforward key-value interface: use set/get/unset commands to manipulate secrets'
  - title: Local Encryption/Decryption
    image: /images/laptop-user.svg
    details: 'All cryptography happens in the client, on your computer: no passwords, keys or plaintext data over the Internet.'
  - title: Highly Configurable
    image: /images/settings.svg
    details: Fully configure options with command line arguments, configuration files or environment variables (no unexpected defaults!)
  - title: Choose Your Own Storage
    image: /images/cloud-storage.svg
    details: 'Multiple storage backend options: local filesystem, AWS S3 or S3-compatible object storage, git repository...'
  - title: Robust Secrecy
    image: /images/lock.svg
    details: 'NIST-approved, industry-grade cryptography standards: AES-256 encryption and Argon2id key derivation'
footer: Copyright © Charles Francoise 2021
---

### Secure secret management from the command line

```sh
# Set store password in environment
$ export SCRT_PASSWORD=*******

# Add a secret to the store
$ scrt set greeting 'Good news, everyone!'

# Retrieve the secret from the store
$ scrt get greeting
Good news, everyone!
```

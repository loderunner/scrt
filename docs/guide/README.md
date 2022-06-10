---
description: |
scrt is a secret manager for the command line. You can encrypt and decrypt secrets from the command line, choosing between a variety of secure storage backends, locally, or remotely over the network.

The purpose of scrt is to provide developers, sysadmins and devops with a reliable and secure tool to manage secrets in projects, while avoiding GUI-centered tools, expensive cloud-based services and complex self-hosted solutions.
---

# Introduction

**scrt** is a secret manager for the command line. You can encrypt and decrypt secrets from the command line, choosing between a variety of secure storage backends, locally, or remotely over the network.

The purpose of scrt is to provide developers, sysadmins and devops with a reliable and secure tool to manage secrets in projects, while avoiding GUI-centered tools, expensive cloud-based services and complex self-hosted solutions.

## How it works

scrt keeps a collection of secrets inside a _store_, a single encrypted file, stored locally on your computer, or remotely over the network.

When performing operations on your secrets, scrt loads the store into memory, decrypts the payload, creates, retrieves or updates secrets, and, if necessary, encrypts the changes back to the store. Each secret is referenced by a name, and can be an arbitrary string of bytes, of any length.

scrt uses [AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) symmetric encryption, and derives its 256-bit keys from a password of your choosing using [Argon2id](https://en.wikipedia.org/wiki/Argon2) key derivation. A new key is derived from the password each time the store is re-encrypted, avoid key re-use and improving security.

The store data can be stored on a file on your computer's hard drive, or use one of the remote storage backends, such as AWS S3 (or any S3-compatible object storage), or a git remote repository.

## When should I use scrt?

The main purpose of scrt is to provide engineering teams with a straightforward way to share secrets, but actual usages can vary. You can:

- Store secrets in a file at the root of a project repository, retrieving them as needed during CI
- Use a git repository to provide quick and easy secret sharing between team members, without relying on an unfamiliar solution
- Retrieve a store from S3, and use the secrets in Ansible or Chef scripts during deployment
- Use as a replacement for complex, expensive SaaS-based secret solutions

## Why not&hellip;?

### sops

sops is actually very similar to scrt, relying on a single file to encrypt and decrypt secrets. A distinguishable feature from sops is its ability to interface with key management systems to fetch encryption keys at run time. In order to accomodate for secrets encrypted with different keys, sops does not encrypt the entire store as a single opaque object, keeping secret keys in plaintext, trading a bit of secrecy for versatility.

### Vault

Hashicorp Vault provides engineering teams with both the tools to handle secrets and the storage engine as a server. A fantastic solution for large teams with many secrets. Yet installing and maintaining a Vault server for your orgnization can require many hours, and usually a dedicated IT team. scrt and Vault are intended for organizations and projects of different scales.

### Cloud secret managers

Most cloud providers offer a secret management solution (along with key management). Very similar in features and operation to Vault, in this case the solution is managed by the cloud provider, allowing an organization to reduce maintenance to nearly nothing. This choice usually comes with a cost on your cloud bill.

### 1Password

While SaaS password managers sometimes offer a command-line or API interface,they are built with the non-technical user in mind, and are usually complicated to adapt to engineering practices.

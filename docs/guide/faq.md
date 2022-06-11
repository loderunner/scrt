---
description: Frequently Asked Questions about scrt, a command-line secret manager for developers, sysadmins and devops.
---

# Frequently Asked Questions

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

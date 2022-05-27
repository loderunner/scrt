---
sidebarDepth: 0
---

# Global

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
  completion  Generate the autocompletion script for the specified shell

Flags:
  -c, --config string     configuration file
  -h, --help              help for scrt
  -p, --password string   master password to unlock the store
      --storage string    storage type
  -v, --verbose           verbose output
      --version           version for scrt
```

### Global options

**`-c`**, **`--config`:** Path to a YAML [Configuration file](/guide/configuration.md)

**`--storage`:** storage type, see [Storage types](/reference/storage.md) for details.

**`-p`**, **`--password`:** password to the store. The argument will be used to derive a key, to decrypt and encrypt the data in the store.

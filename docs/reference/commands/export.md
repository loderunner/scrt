---
sidebarDepth: 0
---

# export

```
scrt export [flags]
```

Exports the keys and values of the store to a file. Please note non-string values are not supported and are considered undefined behavior.

### Options

**`--format`:** the format of the output file. Supported formats are `yaml` `dotenv`, and `json`. Required.

**`--output`:** the path to the output file. Defaults to `stdout`.

### Example

Export the keys and values of the store to a file named `secrets.yaml` in the current directory.

```shell
scrt export --output secrets.yaml --format yaml
```
---
sidebarDepth: 0
---

# set

```
scrt set [flags] key [value]
```

Associate a value to a key in the store. If `value` is omitted from the command
line, it will be read from standard input.

If a value is already set for `key`, the command will fail unless the `--overwrite` option is set.

### Options

**`--overwrite`:** when this flag is set, `scrt` will overwrite the value for `key` in the store, if it exists, instead of returning an error. If no value is associated to `key`, `--overwrite` has no effect.

### Example

Associate `Hello World` to the key `greeting` in the store, using implicit store configuration (configuration file or environment variables).

```shell
scrt set greeting "Hello World"
```

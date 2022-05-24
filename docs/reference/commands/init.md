---
sidebarDepth: 0
---

# init

```
scrt init [flags]
```

Initialize a new store. If an item is already present at the given location, the initialization will fail unless the `--overwrite` option is set.

### Options

**`--overwrite`:** when this flag is set, `scrt` will overwrite the item at the given location, if it exists, instead of returning an error. If no item exists at the location, `--overwrite` has no effect.

### Example

Create a store in a `store.scrt` file in the local filesystem, in the current working directory, using the password `"p4ssw0rd"`.

```shell
scrt init --storage=local --password=p4ssw0rd --local-path=./store.scrt
```

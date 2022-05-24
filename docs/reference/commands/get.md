---
sidebarDepth: 0
---

# get

```
scrt get [flags] key
```

Retrieve the value associated to the key in the store, if it exists. Returns an error if no value is associated to the key.

### Example

Retrieve the value associated to the key `greeting` in the store, using implicit store configuration (configuration file or environment variables).

```shell
scrt get greeting

# Output: Hello World
```

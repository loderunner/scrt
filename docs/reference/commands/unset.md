---
sidebarDepth: 0
---

# unset

```
scrt unset [flags] key
```

Disassociate the value associated to a key in the store. If no value is associated to the key, does nothing.

### Example

Remove the value associated to the key. After this command, no value will be associated to the key `greeting` in the store.

```shell
scrt unset greeting
```

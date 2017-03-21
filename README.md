# Consul launcher

This is a simplified version of consul-template.

The main differences is that it doesn't use any local template. Instead all files with a consul prefix will be downloaded.

Usage:

```
consul-launcher --prefix prefix --destination custom command --with args
```

Where `prefix` is the consul prefix. All the keys with this prefix will be downloaded.

`destination` is a directory where the files will be saved to.

The rest of the arguments could any command (with argument) which will be started (and restarted on any change in the consul).




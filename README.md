parent
------

A minimal process wrapper that simply forwards all signals to the child process and waits for the child to exit.

Useful for containers.

### Usage
```
parent <command> [<args>...]
```

### Use Case

Sample use case for a `Dockerfile` to use env vars in entrypoint.

```Dockerfile
...

ENV MY_ENV=value

ENTRYPOINT [/usr/bin/parent]
CMD["echo", "$MY_ENV"]

```

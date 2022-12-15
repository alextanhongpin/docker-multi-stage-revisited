# docker-multi-stage-revisited

## What is the difference between the syntax below?

```bash
# SHELL form, /bin/sh -c /go/bin/…
CMD ["/go/bin/app"] # output of 47facfd6c8ea

# EXEC form, "/go/bin/app"
CMD /go/bin/app # output of 16554a390f0d
```


> The SHELL form runs the command as a child process (on a shell).

> The EXEC form runs the executable on the main process (the one that has PID 1). [^1]


The `COMMAND` output of `$ docker ps -a`

```
CONTAINER ID   IMAGE                      COMMAND                  CREATED              STATUS                          PORTS     NAMES
16554a390f0d   alextanhongpin/app:0.0.2   "/bin/sh -c /go/bin/…"   About a minute ago   Exited (0) 39 seconds ago                 competent_merkle
47facfd6c8ea   1a73a3d36064               "/go/bin/app"            2 minutes ago        Exited (0) About a minute ago             optimistic_newton
```


This is especially important for _graceful termination_ of the server. If the shell form is used, the app will not receive the correct termination signal.

## Using scratch vs alpine

The issue with `scratch` is the lack of debugging capability (same for `distrolless`, see below).

> FROM scratch literally is an empty, zero-byte image / filesystem, where you add everything yourself [^2].

So when running the following command, we will get an error:

```
$ docker exec -it $(docker ps -q) ps -ef
```
Output:

```
OCI runtime exec failed: exec failed: unable to start container process: exec: "ps": executable file not found in $PATH: unknown
```

The solution is to replace `FROM scratch` with:

```Dockerfile
FROM alpine:latest

# Allow `$ docker exec -it (pid) bash` instead of `$ docker exec -it (pid) /bin/sh`
RUN apk add bash
```


## Distroless

`Distroless` images contains only your application and its runtime dependencies. They do no contain package manager, shells or any other programs you would expect to find in standard Linux distribution. [^3]


## Build


Note that we don't have to specify `GOOS` or `GOARCH` because go will detect it based on the given OS.

```diff
- RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app
+ RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/app
```

## Docker stop vs kill

For graceful shutdown, use `docker stop` instead of `docker kill`. This can be observed by the messages logged by the application (`server exiting`).

```bash
$ docker ps # get the id of the running container
$ docker stop <container> # kill it (gracefully)
$ docker kill <container> # does not stop the process gracefully
```

## Entrypoint vs cmd

See here [^4].

## References

[^1]: https://engineering.pipefy.com/2021/07/30/1-docker-bits-shell-vs-exec/#:~:text=The%20SHELL%20form%20runs%20the,process%20(on%20a%20shell).&text=The%20EXEC%20form%20runs%20the,one%20that%20has%20PID%201).
[^2]: https://github.com/moby/moby/issues/17896
[^3]: https://github.com/GoogleContainerTools/distroless
[^4]: https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#entrypoint

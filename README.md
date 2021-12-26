# Remote Code Execution

## Getting Started

Build the image.
```bash
$ docker build . -t rce-engine
```

Run the container with the image you have created before.
```bash
$ docker run -it -v /var/run/docker.sock:/var/run/docker.sock rce-engine
```

## Build the all-in-one-ubuntu image

Build the ubuntu docker image to run code with different languages.

```bash
$ docker build build/all-in-one-ubuntu -t all-in-one-ubuntu
$ docker run -dit all-in-one-ubuntu
```

Open ubuntu container from command line.
```
$ docker ps # find the container id
$ docker exec -it <container-id> bash
```

Test the compiler/interpreters.

1. __Python -->__ `python3 --version`
2. __Java -->__ `java -version`
2. __Javac -->__ `javac -version`
3. __NodeJS -->__ `nodejs --version`
4. __Golang -->__ `/usr/local/go/bin/go version`
5. __C++ -->__ `g++ --version`
6. __C -->__ `gcc --version`

## Execute multiple commands with docker exec

The command below will work for ubuntu container. You need to change the `bash` command according to container you use. For example for alpine container it should be `/bin/sh`. 

```bash
$ docker exec -w <workdir> -it <container-id> bash -c "<command> && <command>"
```
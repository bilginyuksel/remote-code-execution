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
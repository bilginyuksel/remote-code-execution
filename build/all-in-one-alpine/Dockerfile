FROM alpine:3.15

# echo "installing go version 1.17..."

RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go
RUN wget -O go.tgz https://dl.google.com/go/go1.17.src.tar.gz
RUN tar -C /usr/local -xzf go.tgz
RUN cd /usr/local/go/src && ./make.bash
RUN export PATH="/usr/local/go/bin:$PATH"
RUN apk del .build-deps

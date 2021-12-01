FROM docker:dind

RUN apk add go
RUN apk add python3

WORKDIR /app
COPY . /app

# build phase
FROM golang:1.17

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=0

RUN go build -o rce-engine .

# execution phase
FROM docker:dind

WORKDIR /

RUN mkdir app target .config

COPY --from=0 /app/rce-engine .
COPY --from=0 /app/.config/prod.yml .config/prod.yml

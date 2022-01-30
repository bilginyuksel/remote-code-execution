# build phase
FROM golang:1.17

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=0

RUN go build -o rce-engine .

# execution phase
FROM alpine:latest

USER root

RUN apk --no-cache add ca-certificates

WORKDIR /
RUN mkdir app target .config

COPY --from=0 /app/rce-engine .
COPY --from=0 /app/.config/prod.yml .config/prod.yml
# TODO: Change the app environment when it is prod
ENV APP_ENV=prod

CMD ["./rce-engine", "serve"]

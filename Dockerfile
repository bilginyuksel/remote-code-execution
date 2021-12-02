# FROM docker:dind

# RUN apk add go
# RUN apk add python3

# WORKDIR /app
# COPY . /app/

# RUN ls
# RUN go run .

# RUN systemctl start docker
# RUN dockerd

FROM golang:1.17

WORKDIR /app
COPY . /app

CMD ["go", "run", "."]
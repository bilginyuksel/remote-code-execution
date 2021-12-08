FROM golang:1.17

WORKDIR /app
COPY . /app

CMD ["go", "run", "."]
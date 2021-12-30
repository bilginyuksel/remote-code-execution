FROM golang:1.17

RUN mkdir /app /target

WORKDIR /rce
COPY . /rce
ENV APP_ENV=dev

CMD ["go", "run", ".", "exec", "-p", "testdata/maingo", "-l", "golang"]
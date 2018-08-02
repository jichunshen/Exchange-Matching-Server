FROM golang:1.10.1-alpine3.7
RUN mkdir /code
WORKDIR /code
ADD . /code/
RUN apk add --no-cache git mercurial
RUN go get -u github.com/go-sql-driver/mysql
RUN go get github.com/jmoiron/sqlx
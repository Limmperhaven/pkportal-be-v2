FROM golang:latest
MAINTAINER Mr. Artem Chulaevskiy
LABEL version="0.1"

WORKDIR /app

COPY ./ ./
RUN go mod download

RUN go build -o /main

CMD ["/main"]
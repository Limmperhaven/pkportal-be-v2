FROM golang:latest
MAINTAINER Mr. Artem Chulaevskiy

WORKDIR /app

COPY ./ ./
RUN go mod download

RUN go build -o /main

CMD ["/main"]
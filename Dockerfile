FROM golang:latest
MAINTAINER Mr. Artem Chulaevskiy

WORKDIR /app

COPY ./ ./
RUN apt-get update \
    && apt-get install -y wkhtmltopdf

RUN go mod download

RUN go build -o /main

CMD ["/main"]
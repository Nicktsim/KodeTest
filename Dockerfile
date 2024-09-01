
FROM golang:1.22.4-alpine  

RUN apk update && apk add git bash

COPY . /app

WORKDIR /app

RUN go mod tidy && go build -o main cmd/main.go

EXPOSE 8080/tcp

# точка запуска приложения не содержит параметров - все передается через ENV переменные
ENTRYPOINT [ "/app/main"]
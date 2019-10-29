FROM golang

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 5556
CMD ["./api-cmangos", "contrib/config.ini.dist"]
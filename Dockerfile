FROM golang

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 5556
ENTRYPOINT ["./api-cmangos"]
CMD ["contrib/config.ini.dist"]
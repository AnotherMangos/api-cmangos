all: cmangos-api

cmangos-api:
	go get github.com/gorilla/mux
	go get gopkg.in/ini.v1
	go get github.com/google/uuid
	go get github.com/go-sql-driver/mysql
	go build

clean:
	rm -rf cmangos-api


re: clean all
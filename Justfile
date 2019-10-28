all: cmangos-api

cmangos-api:
	go build

clean:
	rm -rf cmangos-api


re: clean all

docker-build:
	docker build . -t steakhouse.sysroot.ovh/cmangos-api
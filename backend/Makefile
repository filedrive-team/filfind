.PHONY: build

VERSION=./buildversion.go
Major=1
Minjor=0
Patch=0

TARGET=./filfind-backend

build: clean buildversion.go
	go build -o ${TARGET} .

run: build
	${TARGET}

buildversion.go:
	echo package main > ${VERSION}
	echo "const Major = \"${Major}\"" >> ${VERSION}
	echo "const Minjor = \"${Minjor}\"" >> ${VERSION}
	echo "const Patch = \"${Patch}\"" >> ${VERSION}
	echo "const BuildVersion = \"`git --no-pager log --pretty="%h" -n 1`\"" >> ${VERSION}

.PHONY: clean
clean:
	-rm -f ${TARGET}
	-rm -f buildversion.go

swag:
#    go get -u github.com/swaggo/swag/cmd/swag
#	swag init --parseDependency
	swag init -g=api/server/api_server.go

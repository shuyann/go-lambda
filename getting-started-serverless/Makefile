.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go

clean:
	rm -rf ./bin

test:
	go test -v ./...

fmt:
	gofmt -w -s ./

lint:
	golint --set_exit_status ./...

deploy: clean build
	sls deploy --verbose

remove:
	sls remove --verbose

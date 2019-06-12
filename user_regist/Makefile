.PHONY: build clean deploy

BUCKET = awscloudbook

s3put:
	aws s3 cp index.html s3://$(BUCKET)/

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

remove:
	sls remove --verbose

build-sam:
	sam build

local-with-sam: build-sam
	sam local start-api

build-zip:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
	sleep 3
	zip bootstrap.zip bootstrap

build-sam:
	sam build

local-with-sam: build-sam
	sam local start-api

debug-with-sam: build-sam
	sam local start-api -d 5858 --debugger-path $${HOME}/go/bin/linux_amd64 --debug-args="-delveAPI=2" --skip-pull-image

build-zip:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
	sleep 3
	zip bootstrap.zip bootstrap

start-db:
	java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb

create-table:
	aws dynamodb create-table --endpoint-url "http://localhost:8000" --cli-input-json file://init-dynamodb.json

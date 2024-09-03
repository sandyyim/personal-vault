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


aws dynamodb create-table \
  --table-name personal-vault-dynamodb \
  --attribute-definitions \
      AttributeName=id,AttributeType=S \
      AttributeName=name,AttributeType=S \
	  AttributeName=secret,AttributeType=S \
  --key-schema \
      AttributeName=id,KeyType=HASH \
  --endpoint-url=http://localhost:4556 \
  --region=us-east-1

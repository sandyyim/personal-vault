# personal-vault

GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
GOOS=linux GOARCH=amd64 go install github.com/go-delve/delve/cmd/dlv@latest

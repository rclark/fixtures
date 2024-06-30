init:
	go mod tidy
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

doc:
	gomarkdoc -o readme.md .

test:
	go test -v ./...

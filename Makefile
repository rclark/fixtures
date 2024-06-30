init:
	go mod tidy
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
	go install gotest.tools/gotestsum@latest
	go install github.com/axw/gocov/gocov@latest
	go install github.com/matm/gocov-html/cmd/gocov-html@latest
	@if [ ! -e .git/hooks/pre-commit ]; then \
		ln -s $(PWD)/.githooks/pre-commit .git/hooks/pre-commit; \
	fi

doc:
	@gomarkdoc \
		--output readme.md  \
		--template-file file=templates/file.md \
		--template-file package=templates/package.md \
		.

test:
	gotestsum --format testname -- -coverprofile=coverage.out ./...
	gocov convert coverage.out | gocov-html > coverage.html

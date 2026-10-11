TEMPL_CMD := $(HOME)/go/bin/templ

templ:
	$(TEMPL_CMD) generate

run: templ
	go run .

build: templ
	go build -o ./bin/todos .

test:
	go test ./...

vet:
	go vet ./...

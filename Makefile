.PHONY: build clean test daemon ui

build: daemon ui

daemon:
	go build -o bin/defensed ./cmd/defensed

ui:
	go build -o bin/defense-ui ./cmd/defense-ui

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

run-daemon: daemon
	./bin/defensed

run-ui: ui
	./bin/defense-ui

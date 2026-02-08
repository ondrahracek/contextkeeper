.PHONY: build test clean

build:
	go build -o contextkeeper .

test:
	go test ./...

clean:
	rm -f contextkeeper

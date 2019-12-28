.PHONY: build clean

build: clean
	env GOOS=darwin go build -ldflags="-s -w" -o bin/darwin/awsstsgen src/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/linux/awsstsgen src/main.go

clean:
	rm -rf ./bin
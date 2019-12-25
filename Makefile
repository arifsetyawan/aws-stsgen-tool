.PHONY: build clean

build: clean
	env GOOS=darwin go build -ldflags="-s -w" -o bin/darwin/stsgen src/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/linux/stsgen src/main.go

clean:
	rm -rf ./bin
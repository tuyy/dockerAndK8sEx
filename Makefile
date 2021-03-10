.PHONY: build clean

build :
	mkdir -p dist
	env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o dist/first .

clean:
	rm -rf dist

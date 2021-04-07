.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/login login/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/callback callback/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/authorizer authorizer/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

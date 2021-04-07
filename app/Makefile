.PHONY: build clean deploy

build: 
	npm run build

clean:
	rm -rf ./build

deploy: clean build
	aws s3 cp build/ s3://serverless-go-frontend --recursive --acl public-read

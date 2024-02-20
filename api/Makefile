CUR_DIR = $(shell pwd)
LAMBDA_DIR = $(CUR_DIR)/cmd/lambda

# go build variables
GOVARS = GOOS=linux GOARCH=amd64
VERSION = $(shell cat ./VERSION)
BUILT_AT = $(shell TZ=UTC date +%FT%TZ)
LDFLAGS=  -w -s -X \"main.Version=$(VERSION)\" -X \"main.BuiltAt=$(BUILT_AT)\"

test:
	@eval go test -cover $$(go list ./... | grep -v /mocks/) -coverprofile .test_coverage.txt

--build-lambda:
	cd $(LAMBDA_DIR) \
	&& echo "we're on $$(pwd)" \
	&& echo 'cleaning lambda...' && find . -type f -not \( -name '*go' -or -name '*bootstrap' \) -delete \
	&& echo 'building lambda...' && $(GOVARS) go build -tags lambda.norpc -o ./main -ldflags "-s -w" main.go \
	&& echo 'zipping lambda...' && chmod 755 main bootstrap && zip -FS main.zip main

build: --build-lambda
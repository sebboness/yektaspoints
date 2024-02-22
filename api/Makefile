CUR_DIR = $(shell pwd)
LAMBDA_DIR = $(CUR_DIR)/cmd/lambda

# go build variables
GOVARS = GOOS=linux GOARCH=amd64 CGO_ENABLED=0
VERSION = $(shell cat ./VERSION)
BUILT_AT = $(shell TZ=UTC date +%FT%TZ)
LDFLAGS=  -s -w -X \"main.Version=$(VERSION)\" -X \"main.BuiltAt=$(BUILT_AT)\"

test:
	@eval go test -cover $$(go list ./... | grep -v /mocks/) -coverprofile .test_coverage.txt

--build-lambda:
	cd $(LAMBDA_DIR) \
	&& echo "we're on $$(pwd)" \
	&& echo 'cleaning lambda...' && find . -type f -not \( -name '*go' -or -name '*go' \) -delete \
	&& echo 'building lambda...' && $(GOVARS) go build -tags lambda.norpc -o ./bootstrap -ldflags "$(LDFLAGS)" . \
	&& echo 'zipping lambda...' && chmod 755 * && zip -FS bootstrap.zip bootstrap

build: --build-lambda
CUR_DIR = $(shell pwd)
LAMBDA_DIR = $(CUR_DIR)/cmd/lambda

# go build variables
GOVARS = GOOS=linux GOARCH=arm64 CGO_ENABLED=0
VERSION = $(shell cat ./VERSION)
BUILT_AT = $(shell TZ=UTC date +%FT%TZ)
LDFLAGS=  -X \"main.Version=$(VERSION)\" -X \"main.BuiltAt=$(BUILT_AT)\"

# array of all lambda names separated by space
LAMBDAS := request_points get_user_points

test:
	@eval go test -cover $$(go list ./... | grep -v /mocks/) -coverprofile .test_coverage.txt

lambda-check:
ifneq ($(filter $(l),$(LAMBDAS)),)
	$(info $(l) exists in [$(LAMBDAS)])
else
	$(error "$(l)" is not valid. run target with argument l=XXX (valid values are: [$(LAMBDAS)]))
endif

--build-lambda:
	cd $(LAMBDA_DIR)/$(l) \
	&& echo "we're on $$(pwd)" \
	&& echo 'cleaning $(l) lambda...' && find . -type f -not -name '*go' -delete \
	&& echo 'building $(l) lambda...' && $(GOVARS) go build -tags lambda.norpc -a -o ./bootstrap -ldflags "$(LDFLAGS)" . \
	&& echo 'zipping $(l) lambda...' && zip $(l).zip bootstrap

# run with argument l=
# for example: make build l=request_points
build: lambda-check --build-lambda
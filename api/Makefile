CUR_DIR = $(shell pwd)
LAMBDA_DIR = $(CUR_DIR)/cmd/lambda
GOVARS = GOOS=linux GOARCH=arm64 CGO_ENABLED=0
VERSION = $(shell cat ./VERSION)
BUILT_AT = $(shell TZ=UTC date +%FT%TZ)
LDFLAGS=  -X \"main.Version=$(VERSION)\" -X \"main.BuiltAt=$(BUILT_AT)\"
ALL_EXCEPT_MAIN = $(shell ls | grep -v "^main.go$$")

xyz = echo

abc:
	@$(xyz) 'i say abc.' && $(xyz) 'crap'

def:
	@echo $(ALL_EXCEPT_MAIN)

ghi: abc def

clean_request_points:
	@cd $(LAMBDA_DIR)/request_points && rm *.zip

build_request_points:
	@cd $(LAMBDA_DIR)/request_points && $(GOVARS) go build -tags lambda.norpc -a -o ./bootstrap -ldflags "$(LDFLAGS)" .

zip_request_points:
	@cd $(LAMBDA_DIR)/request_points && zip request_points.zip bootstrap

test:
	@eval go test -cover $$(go list ./... | grep -v /mocks/) -coverprofile .test_coverage.txt
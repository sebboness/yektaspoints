test:
	@eval go test -cover $$(go list ./... | grep -v /mocks/) -coverprofile .test_coverage.txt
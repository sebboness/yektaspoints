package error

import (
	"errors"
	"net/http"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go"
	"github.com/sebboness/yektaspoints/util/log"
)

func GetAwsError(err error) *ApiError {
	innerErr := err
	logger := log.Get()
	svc := ""
	op := ""

	// try get http status code
	statusCode := http.StatusBadRequest
	var respErr *awshttp.ResponseError
	if errors.As(err, &respErr) {
		statusCode = respErr.HTTPStatusCode()
	}

	// try get base error
	var oe *smithy.OperationError
	if errors.As(err, &oe) {
		svc = oe.Service()
		op = oe.Operation()
		innerErr = oe.Unwrap()
	}

	logger.AddFields(map[string]any{
		"aws_error":       innerErr.Error(),
		"aws_operation":   op,
		"aws_service":     svc,
		"aws_status_code": statusCode,
	}).Errorf("aws operation failed")

	return New(innerErr).WithStatus(statusCode)
}

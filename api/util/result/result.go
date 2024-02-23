package result

type ResultStatus string

const ResultSuccess ResultStatus = "SUCCESS"
const ResultFailure ResultStatus = "FAILURE"

type Result struct {
	Data    any          `json:"data,omitempty"`
	Errors  []string     `json:"errors"`
	Message string       `json:"message"`
	Status  ResultStatus `json:"status"`
}

func ErrorResult(err error) *Result {
	return &Result{
		Errors: []string{err.Error()},
		Status: ResultFailure,
	}
}

func SuccessResult(data any) *Result {
	return &Result{
		Data:   data,
		Errors: []string{},
		Status: ResultSuccess,
	}
}

func (r *Result) WithMessage(message string) *Result {
	r.Message = message
	return r
}

func (r Result) IsSuccess() bool {
	return r.Status == ResultSuccess
}

package handlers

type Result struct {
	Status  string   `json:"status"`
	Errors  []string `json:"errors"`
	Message string   `json:"message"`
	Data    any      `json:"data,omitempty"`
}

func ResultWithErrors()

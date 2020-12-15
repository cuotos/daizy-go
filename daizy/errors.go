package daizy

type ResponseError struct {
	Success bool `json:"success"`
	Status  int  `json:"status"`
	Errors  []struct {
		Field   string `json:"field"`
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"errors"`
}

func (e *ResponseError) Error() string {
	return e.Errors[0].Message
}

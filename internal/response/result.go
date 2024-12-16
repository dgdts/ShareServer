package response

type Result struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewResultFromError(err error) *Result {
	ret := &Result{}
	if err == nil {
		return ret
	}

	ret.Status = -1
	ret.Message = err.Error()

	return ret
}

func NewResultWithData(data any) *Result {
	return &Result{
		Status: 0,
		Data:   data,
	}
}

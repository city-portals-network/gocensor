package main

// CensorResponse defines http response
type CensorResponse struct {
	Query  string `json:"query"`
	Result bool   `json:"result"`
	Err    error  `json:"error"`
}

func (cr *CensorResponse) setResult(result bool) *CensorResponse {
	cr.Result = result
	return cr
}

func (cr *CensorResponse) setError(err error) *CensorResponse {
	cr.Err = err
	return cr
}

// NewCensorResponse creates new http censor response
func NewCensorResponse(query string) *CensorResponse {
	return &CensorResponse{
		Query:  query,
		Result: false,
	}
}

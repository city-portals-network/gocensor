package main

// CensorResponse defines http response
type CensorResponse struct {
	Query       string `json:"query"`
	Result      bool   `json:"result"`
	Err         error  `json:"error"`
	StatusCode  int    `json:"StatusCode"`
	RequestTime string `json:"RequestTime"`
}

func (cr *CensorResponse) setRequestTime(requestTime string) {
	cr.RequestTime = requestTime
}

// NewCensorResponse creates new http censor response
func NewCensorResponse(query string, result bool, err error) *CensorResponse {
	statusCode := 200
	if err != nil {
		statusCode = 500
	}
	return &CensorResponse{
		Query:      query,
		Result:     result,
		StatusCode: statusCode,
	}
}

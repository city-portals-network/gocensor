package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type ServiceResponse struct {
	Query  string `json:"query"`
	Result bool   `json:"result"`
	// Err         error  `json:"error"`
	StatusCode  int    `json:"StatusCode"`
	RequestTime string `json:"RequestTime"`
}

func (s *server) index(ctx *fasthttp.RequestCtx) {
	t := time.Now()
	ctx.SetContentType("application/json")
	query := string(ctx.QueryArgs().Peek("query"))

	result := s.censor.run(string(query))
	serverResponse := &ServiceResponse{
		Query:       query,
		Result:      result,
		StatusCode:  200,
		RequestTime: time.Since(t).String(),
	}

	body, err := json.Marshal(serverResponse)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	ctx.SetBody(body)
}

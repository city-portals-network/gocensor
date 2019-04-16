package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

func (s *Server) index(ctx *fasthttp.RequestCtx) {
	t := time.Now()
	ctx.SetContentType("application/json")
	query := string(ctx.QueryArgs().Peek("query"))

	result := s.censor.run(string(query))
	serverResponse := NewCensorResponse(query, result, nil)
	serverResponse.setRequestTime(time.Since(t).String())

	body, err := json.Marshal(serverResponse)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	ctx.SetBody(body)
}

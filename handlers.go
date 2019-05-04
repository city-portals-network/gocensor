package main

import (
	"github.com/valyala/fasthttp"
)

func (s *Server) check(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	query := string(ctx.QueryArgs().Peek("query"))
	result := s.censor.check(string(query))
	body := sendSuccessOKWithResult("success", result)
	if result.Err != nil {
		body = sendProblemBadRequest("error", result)
	}
	ctx.SetBody(body)
}

func (s *Server) delete(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	word := ctx.PostArgs().Peek("word")

	if err := s.censor.delete(string(word)); err != nil {
		ctx.SetBody(sendProblemBadRequest("delete error", err))
	} else {
		ctx.SetBody(sendSuccessOK("delete success"))
	}
}
func (s *Server) reload(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	if err := s.censor.reload(); err != nil {
		ctx.SetBody(sendProblemBadRequest("reload error", err))
	} else {
		ctx.SetBody(sendSuccessOK("reload success"))
	}
}

func (s *Server) append(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	word := ctx.PostArgs().Peek("word")

	if err := s.censor.append(string(word)); err != nil {
		ctx.SetBody(sendProblemBadRequest("append error", err))
	} else {
		ctx.SetBody(sendSuccessOK("append success"))
	}
}

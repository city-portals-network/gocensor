package main

// Routes defines all routes in gocensor server
func (s *Server) Routes() {
	s.router.GET("/v1/censor/check/", s.check)
	s.router.GET("/v1/censor/reload/", s.reload)
	s.router.POST("/v1/censor/append/", s.append)
	s.router.POST("/v1/censor/delete/", s.delete)
}

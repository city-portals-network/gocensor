package main

// Routes defines all routes in gocensor server
func (s *Server) Routes() {
	s.router.GET("/", s.index)
	// s.router.GET("/long", s.longRunning
}

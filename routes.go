package main

func (s *server) Routes() {

	// Basics
	s.router.GET("/", s.index)
	// s.router.GET("/long", s.longRunning)

	// // Kubernetes
	// s.router.GET("/healthz", s.healthz)
	// s.router.GET("/readyz", s.readyz)

	// // Monitoring
	// s.router.GET("/stats", expvarhandler.ExpvarHandler)

}

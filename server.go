package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/buaazp/fasthttprouter"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

// Server defines all data
type Server struct {
	censor     *Censor
	HTTPServer *fasthttp.Server
	router     *fasthttprouter.Router
	cfg        *Config
}

// NewServer creates a new HTTP Server
func NewServer(censor *Censor, cfg *Config) *Server {
	// define router
	router := fasthttprouter.New()
	// compression
	handler := router.Handler
	handler = fasthttp.CompressHandler(handler)
	return &Server{
		censor:     censor,
		HTTPServer: newHTTPServer(handler, cfg),
		router:     router,
		cfg:        cfg,
	}
}

// newHTTPServer creates a new HTTP Server
func newHTTPServer(h fasthttp.RequestHandler, cfg *Config) *fasthttp.Server {
	return &fasthttp.Server{
		Handler:              h,
		ReadTimeout:          cfg.Server.ReadTimeout,
		WriteTimeout:         cfg.Server.WriteTimeout,
		MaxKeepaliveDuration: cfg.Server.MaxKeepalive,
	}
}

// Run starts the HTTP server and performs a graceful shutdown
func (s *Server) Run() error {
	ln, err := reuseport.Listen("tcp4", "localhost:"+s.cfg.Server.Port)
	if err != nil {
		return errors.Wrap(err, "error in reuseport listener")
	}

	// create a graceful shutdown listener
	graceful := NewGracefulListener(ln, s.cfg.Server.GracefulTimeout)
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "hostname reuse port")
	}

	// Error handling
	listenErr := make(chan error, 1)

	/// Run server
	go func() {
		log.Printf("%s - Web server starting on port http://%v", hostname, graceful.Addr())
		log.Printf("%s - Press Ctrl+C to stop", hostname)
		// listenErr <- s.HTTPServer.ListenAndServe(":" + cfg.Port)
		listenErr <- s.HTTPServer.Serve(graceful)

	}()

	// SIGINT/SIGTERM handling
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// Handle channels/graceful shutdown
	for {
		select {
		case err := <-listenErr:
			if err != nil {
				return errors.Wrap(err, "port in use")
			}
		// handle termination signal
		case <-osSignals:
			fmt.Printf("\n")
			log.Printf("%s - Shutdown signal received.\n", hostname)

			// Servers in the process of shutting down should disable KeepAlives
			// FIXME: This causes a data race
			s.HTTPServer.DisableKeepalive = true

			// Attempt the graceful shutdown by closing the listener
			// and completing all inflight requests.
			if err := graceful.Close(); err != nil {
				log.Fatalf("error with graceful close: %s", err)
			}

			log.Printf("%s - Server gracefully stopped.\n", hostname)
			os.Exit(0)
		}
	}
}

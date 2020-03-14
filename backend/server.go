package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	// "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/atomic"
	"google.golang.org/grpc"

	"go.uber.org/zap"

	"github.com/sergivb01/acmecopy/api"
)

import _ "net/http/pprof"

type Server struct {
	router *mux.Router

	clientConn *grpc.ClientConn
	cli        api.CompilerClient

	// db  *sqlx.DB
	log *zap.Logger
	cfg Config

	healthy *atomic.Bool
}

func NewServer() (*Server, error) {
	c, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load config: %w", err)
	}

	// db, err := sqlx.Open("postgres", c.PostgresURI)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not connect to database: %w", err)
	// }

	logConfig := zap.NewDevelopmentConfig()
	if c.Production {
		logConfig = zap.NewProductionConfig()
	}
	logConfig.DisableStacktrace = true
	logConfig.DisableCaller = true

	logger, err := logConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("could not create logger: %w", err)
	}

	return &Server{
		cfg: *c,
		// db:      db,
		log:     logger,
		healthy: atomic.NewBool(false),
	}, nil
}

func (s *Server) routes() {
	router := mux.NewRouter()

	router.HandleFunc("/test", s.handleIndex).Methods("GET")
	router.HandleFunc("/api/health", s.handleHealthz()).Methods("GET")
	router.HandleFunc("/api/submit", s.handleSubmit()).Methods("POST")

	s.router = router
}

func (s *Server) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set correct RemoteAddr for Reverse Proxies (Cloudflare)
		hdr := r.Header.Get("X-Forwarded-For")
		if hdr != "" {
			r.RemoteAddr = hdr
		}

		defer func() {
			s.log.Info("http server",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("address", r.RemoteAddr))
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Listen() {
	// Configure routes
	s.routes()

	if err := s.createGRPCClient(); err != nil {
		s.log.Fatal("couldn't create GRPC client", zap.Error(err))
		return
	}

	if !s.cfg.Production {
		go registerDebugServer()
		s.log.Debug("registered pprof http server", zap.String("address", ":8083"))
	}

	srv := &http.Server{
		Addr:         s.cfg.Listen,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 15,
		Handler:      s.loggerMiddleware(s.router),
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2"},
		},
	}

	// Run our server in a goroutine so that it doesn't block
	go func() {
		// Should initialize Atomic Health, but is being done in NewServer
		s.log.Info("started listening HTTP server", zap.String("address", s.cfg.Listen))
		s.healthy.Store(true)

		if err := srv.ListenAndServeTLS(s.cfg.TLSCert, s.cfg.TLSKey); err != nil {
			s.log.Fatal("failed to start server", zap.Error(err))
			s.healthy.Store(false)
			return
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal
	<-c

	s.log.Info("closing HTTP server and database pool")
	// if err := s.db.Close(); err != nil {
	// 	s.log.Fatal("failed to close database connection", zap.Error(err))
	// }

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Set healthy to false for monitoring and disable keep alive to close keep-alive connections
	s.healthy.Store(false)
	srv.SetKeepAlivesEnabled(false)

	if err := s.clientConn.Close(); err != nil {
		s.log.Fatal("failed to close GRPC client conn", zap.Error(err))
	}

	if err := srv.Shutdown(ctx); err != nil {
		s.log.Fatal("failed to shut down server gracefully", zap.Error(err))
	}
}

func registerDebugServer() {
	server := &http.Server{
		Addr:         ":8888",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

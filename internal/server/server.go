package server

import (
	"context"
	"github.com/Limmperhaven/pkportal-be-v2/internal/config"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/middlewares"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server http.Server
	ch     chan os.Signal
}

func NewServer(cfg *config.Server, c *controllers.ControllerStorage, m *middlewares.MiddlewareStorage) *Server {
	router := gin.Default()
	initRoutes(router, c, m)
	handler := initCors(router)
	return &Server{
		server: http.Server{
			Addr:    cfg.Host + ":" + cfg.Port,
			Handler: handler,
			//Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		ch: make(chan os.Signal, 1),
	}
}

func (s *Server) Run() {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("There's an error with the server: %s", err.Error())
		}
	}()
	log.Println("Server started on", s.server.Addr)
	s.wait()
}

func (s *Server) wait() {
	defer func() {
		s.ch <- os.Interrupt
	}()
	signal.Notify(s.ch,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-s.ch
	s.close()
}

func (s *Server) close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdwon error: %s", err.Error())
	}
}

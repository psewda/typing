package server

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/psewda/typing/pkg/controllers"
	"github.com/psewda/typing/pkg/log"
)

// Server is the http server implementation
// using 'echo' framework. It is central object
// for RESTful APIs in this app.
type Server struct {
	Running bool
	echo    *echo.Echo
	logger  *log.Logger
}

// RegisterController adds the specified controller and its routes
// in the server runtime environment.
func (s *Server) RegisterController(ctrl controllers.Controller) {
	if ctrl != nil {
		ctrl.AddRoutes(s.echo)
	}
}

// Run starts the http server and returns
// the error object for any failure.
func (s *Server) Run(port uint16) error {
	if !s.Running {
		if err := s.start(port); err != nil {
			return err
		}
		s.Running = true
	}
	return nil
}

// Shutdown stops the http server gracefully
// and resets the server state.
func (s *Server) Shutdown() {
	if s.Running {
		const waitTime = time.Second * 5
		ctx, cancel := context.WithTimeout(context.Background(), waitTime)
		defer cancel()
		if err := s.echo.Shutdown(ctx); err != nil {
			s.logger.Fatal("error in server shutdown, stopping the service", err)
		}
		s.Running = false
	}
}

// New creates a new instance of http server.
func New(banner bool, logger *log.Logger) *Server {
	e := echo.New()

	// hide banner if flag is false
	if !banner {
		e.HideBanner = true
		e.HidePort = true
	}

	// set log header
	e.Logger.SetHeader(log.Header)

	// add middlewares here
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return &Server{
		Running: false,
		echo:    e,
		logger:  logger,
	}
}

// GetRandPort generates a random port number from 1024 to 65535.
func GetRandPort() uint16 {
	src := rand.NewSource(time.Now().UnixNano())

	// nolint:gosec // we don't need secure random number but faster one, there
	// is no security involved here. In most cases user will pass port
	// number, so this will not be used.
	rnd := rand.New(src)
	return uint16(rnd.Intn(math.MaxUint16-1024) + 1024)
}

func (s *Server) start(port uint16) error {
	const waitTime = time.Second
	waitChan := time.After(waitTime)
	errChan := make(chan error)

	if port < 1024 {
		err := errors.New("port number should be in private/dynamic port range")
		return err
	}

	go func() {
		if err := s.echo.Start(fmt.Sprintf(":%d", port)); err != http.ErrServerClosed {
			if s.Running {
				s.logger.Fatal("error on 'echo' server, stoping the service", err)
			}
			errChan <- err
		}
	}()

	// timeout or wait for server to start
	var err error = nil
	select {
	case <-waitChan:
	case err = <-errChan:
	}
	return err
}

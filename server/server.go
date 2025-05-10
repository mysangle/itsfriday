package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
	
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"

	apiv1 "itsfriday/server/router/api/v1"
	"itsfriday/server/profile"
	"itsfriday/store"
)

type Server struct {
	Secret     string
	Profile    *profile.Profile
	Store      *store.Store

	echoServer *echo.Echo
}

func NewServer(ctx context.Context, profile *profile.Profile, store *store.Store) (*Server, error) {
	s := &Server{
		Store:   store,
		Profile: profile,
	}

	secret := "itsfriday"
	s.Secret = secret

	echoServer := echo.New()
	echoServer.Debug = true
	echoServer.HideBanner = true
	echoServer.HidePort = true
	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogRemoteIP: true,
		LogError:    true,
		LogLatency:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			var errorMessage string
			if v.Error == nil {
				errorMessage = "_"
			} else {
				errorMessage = v.Error.Error()
			}
			slog.Debug("REQUEST: ",
				"remote_ip: ", v.RemoteIP,
				", uri: ", v.URI,
				", status: ", v.Status,
				", error: ", errorMessage,
				", latency: ", v.Latency.String(),
			)
			return nil
		},
	}))
	authHandler := apiv1.NewAuthHandler(store, secret, "user")
	echoServer.Use(echojwt.WithConfig(echojwt.Config{
		ContextKey: authHandler.ContextKey,
		ContinueOnIgnoredError: true,
		ParseTokenFunc: authHandler.ParseTokenFunc,
		ErrorHandler: authHandler.ErrorHandler,
	  }))
	  echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins:     []string{fmt.Sprintf("http://localhost:%s", profile.Port)},
        AllowMethods:     []string{http.MethodGet, http.MethodPost},
        AllowCredentials: true,
    }))
	s.echoServer = echoServer

	echoServer.GET("/monitor/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "{\"status\":\"UP\"}")
	})

	apiv1.NewAPIV1Service(s.Secret, profile, store, echoServer)

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.Profile.Addr, s.Profile.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		s.echoServer.Listener = listener
		if err := s.echoServer.Start(address); err != nil {
			slog.Error("failed to start echo server", "error", err)
		}
	}()
    return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Shutdown echo server.
	if err := s.echoServer.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	// Close database connection.
	if err := s.Store.Close(); err != nil {
		slog.Error("failed to close database", slog.String("error", err.Error()))
	}

	slog.Info("itsfriday stopped properly")
}

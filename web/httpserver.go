package web

import (
	"net/http"
	"spotifies-be/config"
	"spotifies-be/externalapi"
	"spotifies-be/web/handler"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
)

type httpServer struct {
	cfg      config.Config
	engine   *gin.Engine
	handlers httpHandlers
}

type httpHandlers struct {
	auth    handler.Auth
	profile handler.Profile
}

func NewHTTPServer(cfg config.Config, sessStore sessions.Store, spotifyAuthClient spotifyAuth.Authenticator, spotifyClient externalapi.SpotifyClient) *httpServer {
	if cfg.Gin.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	ginEngine := gin.Default()

	handlers := httpHandlers{
		*handler.NewAuth(cfg, sessStore, spotifyAuthClient),
		*handler.NewProfile(cfg, spotifyClient),
	}

	srv := httpServer{
		cfg:      cfg,
		engine:   ginEngine,
		handlers: handlers,
	}

	srv.setupRouter()

	return &srv
}

func (h *httpServer) Start() error {
	return h.engine.Run(h.cfg.ServerAddress)
}

func (h *httpServer) setupRouter() {
	router := h.engine

	router.Use(RequestID(), LogRequest())

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	router.GET("/login", h.handlers.auth.Login)

	router.GET("/callback", h.handlers.auth.Callback)

	router.GET("/me", Token(), h.handlers.profile.Me)
}

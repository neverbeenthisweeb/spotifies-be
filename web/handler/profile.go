package handler

import (
	"net/http"
	"spotifies-be/config"
	"spotifies-be/util/log"
	"spotifies-be/web/response"

	"spotifies-be/externalapi"

	"github.com/gin-gonic/gin"
)

type Profile struct {
	cfg           config.Config
	spotifyClient externalapi.SpotifyClient
}

func NewProfile(cfg config.Config, spotifyClient externalapi.SpotifyClient) *Profile {
	return &Profile{
		cfg:           cfg,
		spotifyClient: spotifyClient,
	}
}

func (h *Profile) Me(c *gin.Context) {
	// FIXME: Add use case layer?
	res, err := h.spotifyClient.GetMe(c)
	if err != nil {
		log.Error(c, "Failed to get me", err)
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

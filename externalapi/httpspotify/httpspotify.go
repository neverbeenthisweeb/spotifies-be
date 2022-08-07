package httpspotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"spotifies-be/config"
	"spotifies-be/model"
	"spotifies-be/util/log"

	"github.com/go-resty/resty/v2"
)

const (
	baseUrl = "https://api.spotify.com/v1"
	timeout = 60 * time.Second
)

type HTTPSpotifyClient struct {
	restyClient *resty.Client
}

func NewHTTPSpotifyClient(cfg config.Config) *HTTPSpotifyClient {
	rc := resty.New()
	rc.SetBaseURL(baseUrl)
	rc.SetTimeout(timeout)
	rc.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		if r.IsError() {
			var v ErrorResponse
			err := json.Unmarshal(r.Body(), &v)
			if err != nil {
				log.Error(r.Request.Context(), fmt.Sprintf("Failed to unmarshal response: %s", string(r.Body())), err)
				return err
			}
			return errors.New(v.Error.Message)
		}

		return nil
	})

	return &HTTPSpotifyClient{
		restyClient: rc,
	}
}

func (h *HTTPSpotifyClient) GetMe(ctx context.Context) (*model.Profile, error) {
	tok, err := getToken(ctx)
	if err != nil {
		log.Error(ctx, "Failed to get token", err)
		return nil, err
	}

	var result model.Profile

	_, err = h.restyClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(tok).
		SetResult(&result).
		Get("/me")
	if err != nil {
		return nil, err
	}

	return &result, err
}

func getToken(ctx context.Context) (string, error) {
	tok := ctx.Value("access_token")
	if tok == nil {
		err := errors.New("missing access token")
		log.Error(ctx, "Missing access token from context", err)
		return "", err
	}

	return tok.(string), nil
}

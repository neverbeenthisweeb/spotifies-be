package externalapi

import (
	"context"
	"spotifies-be/model"
)

type SpotifyClient interface {
	GetMe(ctx context.Context) (*model.Profile, error)
}

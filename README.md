# Spotifies BE

## How to run

1. Copy `.env.example` to `.env`
1. Set the `SESSION_KEY` in your `.env` ([how to generate secure session key](https://github.com/gorilla/sessions#sessions))
1. Set the `CLIENT_ID` for your Spotify app ([how to setup Spotify app](https://developer.spotify.com/documentation/general/guides/authorization/app-settings/))
1. Install dependencies `go mod tidy`
1. Run web server `go run .`
1. Hit `/login` to get your token
1. Hit `/me` to get your profile (`/me` endpoint requires token)

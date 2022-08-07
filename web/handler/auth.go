package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"spotifies-be/config"
	"spotifies-be/util/log"
	"spotifies-be/web/constant"
	"spotifies-be/web/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	maxAgeBufferTimeSec = 5 * 60
)

type Auth struct {
	cfg               config.Config
	sessStore         sessions.Store
	spotifyAuthClient spotifyAuth.Authenticator
}

func NewAuth(cfg config.Config, sessStore sessions.Store, spotifyAuthClient spotifyAuth.Authenticator) *Auth {
	return &Auth{
		cfg:               cfg,
		sessStore:         sessStore,
		spotifyAuthClient: spotifyAuthClient,
	}
}

// func (h *AuthHandler) PreLogin(c *gin.Context) {}

func (h *Auth) Login(c *gin.Context) {
	sess, err := h.sessStore.New(c.Request, constant.CookieName)
	if err != nil {
		log.Error(c, "Failed to create new session", err)
		response.Error(c, err)
		return
	}

	// Generate state to provide protection against attacks such as cross-site request forgery
	// https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.1
	state := fmt.Sprintf("%x", securecookie.GenerateRandomKey(8))
	sess.Values[constant.SessKeyState] = state

	codeVerifier := "w0HfYrKnG8AihqYHA9_XUPTIcqEXQvCQfOF2IitRgmlF43YWJ8dy2b49ZUwVUOR.YnvzVoTBL57BwIhM4ouSa~tdf0eE_OmiMC_ESCcVOe7maSLIk9IOdBhRstAxjCl7"
	sess.Values[constant.SessKeyCodeVerifier] = codeVerifier

	codeChallenge := generateCodeChallenge(codeVerifier)
	sess.Values[constant.SessKeyCodeChallenge] = codeChallenge

	authURL := h.spotifyAuthClient.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	// Set login callback redirect URL (to know here does this user's previous Spotifes' page)
	// sess.Values[keyFromPage] = getRequestFullURL(c)

	err = sess.Save(c.Request, c.Writer)
	if err != nil {
		log.Error(c, "Failed to save session", err)
		response.Error(c, err)
		return
	}

	c.Redirect(http.StatusFound, authURL)
}

func (h *Auth) Callback(c *gin.Context) {
	sess, err := h.sessStore.Get(c.Request, constant.CookieName)
	if err != nil {
		log.Error(c, "Failed to get session", err)
		response.Error(c, err)
		return
	}

	state := fmt.Sprintf("%v", sess.Values[constant.SessKeyState])
	codeVerifier := fmt.Sprintf("%v", sess.Values[constant.SessKeyCodeVerifier])

	tok, err := h.spotifyAuthClient.Token(c, state, c.Request,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		log.Error(c, "Failed to get token, deleting cookies", err)
		sess.Options.MaxAge = -1
		err2 := sess.Save(c.Request, c.Writer)
		if err2 != nil {
			log.Error(c, "Failed to delete cookies", err)
			response.Error(c, err)
			return
		}
		response.Error(c, err)
		return
	}

	sess.Values[constant.SessKeyAccessToken] = tok.AccessToken
	sess.Options.MaxAge = int(tok.Expiry.Unix()-time.Now().Unix()) - maxAgeBufferTimeSec

	// Remove unused values
	delete(sess.Values, constant.SessKeyState)
	delete(sess.Values, constant.SessKeyCodeVerifier)
	delete(sess.Values, constant.SessKeyCodeChallenge)

	err = sess.Save(c.Request, c.Writer)
	if err != nil {
		log.Error(c, "Failed to save session", err)
		response.Error(c, err)
		return
	}

	c.String(http.StatusOK, "OK")
}

func generateCodeChallenge(codeVerifier string) string {
	bCodeChallenge := sha256.Sum256([]byte(codeVerifier))

	codeChallenge := base64.StdEncoding.EncodeToString(bCodeChallenge[:])
	// We also need to do some formatting to the base64
	// https://community.spotify.com/t5/Spotify-for-Developers/Unable-to-use-PKCE-authorization-code-verifier-was-incorrect/td-p/5006416
	codeChallenge = strings.Replace(codeChallenge, "+", "-", -1)
	codeChallenge = strings.Replace(codeChallenge, "/", "_", -1)
	codeChallenge = strings.TrimSuffix(codeChallenge, "=")

	return codeChallenge
}

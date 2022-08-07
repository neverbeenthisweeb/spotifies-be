package web

import (
	"errors"
	"fmt"
	"spotifies-be/config"
	"spotifies-be/util/log"
	"spotifies-be/web/constant"
	"spotifies-be/web/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

var sessStore = sessions.NewCookieStore([]byte(config.LoadConfig().SessionKey))

func Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := sessStore.Get(c.Request, constant.CookieName)
		if err != nil {
			c.Abort()
			log.Error(c, "Failed to get session", err)
			response.Error(c, err)
			return
		}

		if sess.Values[constant.SessKeyAccessToken] == nil {
			c.Abort()
			vErr := errors.New("missing access token")
			log.Error(c, "Failed to get session", vErr)
			response.Error(c, vErr)
			return
		}

		accessToken := fmt.Sprintf("%v", sess.Values[constant.SessKeyAccessToken])

		c.Set(string(constant.SessKeyAccessToken), accessToken)
		c.Next()
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// FIXME: For future improvement, also consider if FE sends request ID too
		if c.Value(constant.CtxKeyRequestID) == nil {
			c.Set(constant.CtxKeyRequestID, uuid.NewString())
			c.Next()
		}
	}
}

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqLogger := log.Logger

		startTime := time.Now()
		c.Next()
		latency := float64(time.Since(startTime).Nanoseconds()) / 10e6
		urlPath := c.Request.URL.String()
		requestID := c.GetString(constant.CtxKeyRequestID)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()

		reqLogger = reqLogger.With(zap.String("latency", fmt.Sprintf("%.2f ms", latency)),
			zap.String("url_path", urlPath),
			zap.String("request_id", requestID),
			zap.Int("status_code", statusCode),
			zap.String("client_ip", clientIP),
			zap.String("client_user_agent", clientUserAgent),
			zap.Int("data_size", dataSize),
		)

		if len(c.Errors) > 0 {
			ginMdlErrors := aggregateErrors(c)
			reqLogger = reqLogger.With(zap.String("gin_middleware_errors", ginMdlErrors))
		}

		switch {
		case statusCode > 499:
			reqLogger.Error("-")
		case statusCode > 399:
			reqLogger.Warn("-")
		default:
			reqLogger.Info("-")
		}
	}
}

func aggregateErrors(c *gin.Context) string {
	var errMsg string
	for _, v := range c.Errors {
		errMsg += v.Error() + "; "
	}
	return errMsg
}

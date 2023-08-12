package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/liquiddev99/dropbyte-backend/token"
)

func authMiddleware(token token.Token) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("authorization")

		if len(authHeader) == 0 {
			err := errors.New("Authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, responseError(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("Invalid Authorization header formar")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, responseError(err))
			return
		}

		if authType := strings.ToLower(fields[0]); authType != "bearer" {
			err := fmt.Errorf("Unsupport authorization type %s", authType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, responseError(err))
			return
		}

		access_token := fields[1]
		payload, err := token.VerifyToken(access_token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, responseError(err))
			return
		}

		err = payload.CheckExpired()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, responseError(err))
			return
		}

		ctx.Set("payload", payload)
		log.Println("Next")
		ctx.Next()
	}
}

package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hfleury/bk_globalshot/internal/dto"
	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/pkg/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.FromError(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.FromError(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.FromError(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.FromError(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func RequireRoles(roles ...model.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, exists := ctx.Get(authorizationPayloadKey)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedResponse("User not authorized"))
			return
		}

		tokenPayload, ok := payload.(*token.Payload)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.InternalServerErrorResponse())
			return
		}

		for _, role := range roles {
			if model.Role(tokenPayload.Role) == role {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, dto.ForbiddenResponse("Access denied: insufficient permissions"))
	}
}

func GetAuthPayload(ctx *gin.Context) *token.Payload {
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return nil
	}
	return payload.(*token.Payload)
}

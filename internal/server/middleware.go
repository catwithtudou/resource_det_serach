package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"resource_det_search/api"
	"resource_det_search/internal/utils"
	"strings"
	"time"
)

func AuthJwtMw(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			logger.Errorf("[AuthJwtMw]the header is empty")
			c.JSON(http.StatusOK, api.UserAuthErr)
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			logger.Errorf("[AuthJwtMw]the header form is wrong")
			c.JSON(http.StatusOK, api.UserAuthErr)
			c.Abort()
			return
		}

		claim, err := utils.ParseJwtToken(parts[1])
		if err != nil {
			logger.Errorf("[AuthJwtMw]failed to parse token:err=[%+v]", err)
			c.JSON(http.StatusOK, api.UserAuthErr)
			c.Abort()
			return
		}

		if time.Now().Unix() >= claim.ExpiresAt {
			logger.Errorf("[AuthJwtMw]token expire")
			c.JSON(http.StatusOK, api.UserAuthErr)
			c.Abort()
			return
		}

		if claim.Uid <= 0 {
			logger.Errorf("[AuthJwtMw]invalid uid")
			c.JSON(http.StatusOK, api.UserAuthErr)
			c.Abort()
			return
		}

		c.Set("uid", claim.Uid)
		c.Set("loginTs", claim.LoginTs)
		c.Next()
	}
}

package api

import (
	"fmt"
	"net/http"

	"github.com/Zhousiru/inker/config"
	"github.com/Zhousiru/inker/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if debugMode {
			c.Next()
			return
		}

		query, err := autoQuery(c, map[string]bool{"token": true})
		if err != nil {
			c.Abort()
			return
		}

		token, err := jwt.Parse(query["token"], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err := fmt.Errorf("invalid token: unexpected signing method: %v", token.Header["alg"])
				c.Abort()
				response(c, http.StatusForbidden, "invalid token", nil, err)
				return nil, err
			}

			return []byte(config.Conf.Key), nil
		})

		if err != nil {
			c.Abort()
			response(c, http.StatusForbidden, "invalid token", nil, err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			c.Abort()
			response(c, http.StatusForbidden, "invalid token", nil, nil)
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			c.Abort()
			response(c, http.StatusForbidden, "invalid token", nil, nil)
			return
		}

		isMatched, err := db.CheckUsername(username)
		if err != nil {
			c.Abort()
			stdResponse(c, typeInternalError, err)
			return
		}

		if !isMatched {
			c.Abort()
			response(c, http.StatusForbidden, "invalid token", nil, nil)
			return
		}

		c.Set("tokenClaims", claims)
		c.Next()
	}
}

package api

import (
	"net/http"
	"time"

	"github.com/Zhousiru/inker/config"
	"github.com/Zhousiru/inker/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func genToken(username string, lifetime int64) (string, error) {
	exp := time.Now().Unix() + lifetime
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username, "exp": exp})

	return token.SignedString([]byte(config.Conf.Key))
}

func login(c *gin.Context) {
	query, err := autoQuery(c, map[string]bool{"username": true, "password": true, "remember": false})
	if err != nil {
		return
	}

	var lifetime int64
	switch query["remember"] {
	case "", "false":
		lifetime = config.Conf.NormalTokenLifetime
	case "true":
		lifetime = config.Conf.RememberTokenLifetime
	default:
		stdResponse(c, typeInvalidParameter, nil)
		return
	}

	isMatched, err := db.Auth(query["username"], query["password"])
	if err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}
	if !isMatched {
		stdResponse(c, typeInvalidParameter, nil)
		return
	}

	token, err := genToken(query["username"], lifetime)
	if err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}
	response(c, http.StatusOK, "", gin.H{"token": token}, nil)
	return
}

func updateUser(c *gin.Context) {
	username := c.MustGet("tokenClaims").(jwt.MapClaims)["username"].(string)
	query, _ := autoQuery(c, map[string]bool{"newUsername": false, "newPassword": false})

	if query["newUsername"] == "" && query["newPassword"] == "" {
		stdResponse(c, typeMissingParameter, nil)
		return
	}

	if err := db.UpdateUser(username, query["newUsername"], query["newPassword"]); err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}

	response(c, http.StatusOK, "", nil, nil)
	return
}

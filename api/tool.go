package api

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

type respType struct {
	StatusCode int
	Msg        string
}

const (
	typeMissingParameter = iota
	typeInvalidParameter
	typeInternalError
)

var stdTypeMap = map[int]respType{
	typeMissingParameter: respType{StatusCode: http.StatusBadRequest, Msg: "missing required parameter"},
	typeInvalidParameter: respType{StatusCode: http.StatusBadRequest, Msg: "invalid parameter"},
	typeInternalError:    respType{StatusCode: http.StatusInternalServerError, Msg: "internal server error"}}

func response(c *gin.Context, code int, msg string, payload interface{}, err error) {
	if debugMode {
		if err != nil {
			c.JSON(code, gin.H{"code": code, "msg": msg, "payload": payload, "debug": gin.H{"stack": fmt.Sprintf("%s", debug.Stack()), "error": err.Error()}})
			return
		}
		c.JSON(code, gin.H{"code": code, "msg": msg, "payload": payload, "debug": gin.H{"stack": fmt.Sprintf("%s", debug.Stack()), "error": nil}})
		return
	}
	c.JSON(code, gin.H{"code": code, "msg": msg, "payload": payload})
}

func stdResponse(c *gin.Context, stdRespType int, err error) {
	response(c, stdTypeMap[stdRespType].StatusCode, stdTypeMap[stdRespType].Msg, nil, err)
}

func autoQuery(c *gin.Context, key map[string]bool) (map[string]string, error) {
	result := map[string]string{}
	for k, required := range key {
		v := c.Query(k)
		if v == "" && required {
			stdResponse(c, typeMissingParameter, nil)
			return nil, errors.New(stdTypeMap[typeMissingParameter].Msg)
		}
		result[k] = v
	}

	return result, nil
}

func autoPostForm(c *gin.Context, key map[string]bool) (map[string]string, error) {
	result := map[string]string{}
	for k, required := range key {
		v := c.PostForm(k)
		if v == "" && required {
			stdResponse(c, typeMissingParameter, nil)
			return nil, errors.New(stdTypeMap[typeMissingParameter].Msg)
		}
		result[k] = v
	}

	return result, nil
}

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Zhousiru/inker/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func newArticle(c *gin.Context) {
	postForm, err := autoPostForm(c, map[string]bool{
		"name": true, "title": true, "content": true, "attr": false})
	if err != nil {
		return
	}

	exist, err := db.ExistArticle(postForm["name"])
	if err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}
	if exist {
		response(c, http.StatusBadRequest, "name already exists", nil, nil)
		return
	}

	var attrMap map[string]interface{}
	if postForm["attr"] != "" {
		err := json.Unmarshal([]byte(postForm["attr"]), &attrMap)
		if err != nil {
			stdResponse(c, typeInvalidParameter, err)
			return
		}
	}

	if err := db.NewArticle(postForm["name"], postForm["title"], postForm["content"], attrMap); err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}
	response(c, http.StatusOK, "", nil, nil)
}

func getArticle(c *gin.Context) {
	query, err := autoQuery(c, map[string]bool{"name": true})
	if err != nil {
		return
	}

	article, err := db.GetArticle(query["name"])
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response(c, http.StatusNotFound, "name not exists", nil, nil)
			return
		}
		stdResponse(c, typeInternalError, err)
		return
	}

	response(c, http.StatusOK, "", article, nil)
}

func paginateHome(c *gin.Context) {
	query, _ := autoQuery(c, map[string]bool{"skip": false, "limit": false})

	var skip int64
	var err error
	if query["skip"] == "" {
		skip = 0
	} else {
		skip, err = strconv.ParseInt(query["skip"], 10, 64)
		if err != nil {
			stdResponse(c, typeInvalidParameter, err)
			return
		}
	}

	var limit int64
	if query["limit"] == "" {
		limit = 0
	} else {
		limit, err = strconv.ParseInt(query["limit"], 10, 64)
		if err != nil {
			stdResponse(c, typeInvalidParameter, err)
			return
		}
	}

	var articleSlice []map[string]interface{}
	articleSlice, err = db.GetRecentArticleByRange(skip, limit)

	if err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}

	response(c, http.StatusOK, "", articleSlice, nil)
}

func paginateSearch(c *gin.Context) {
	query, err := autoQuery(c, map[string]bool{"search": true, "skip": false, "limit": false})
	if err != nil {
		return
	}

	var skip int64
	if query["skip"] == "" {
		skip = 0
	} else {
		skip, err = strconv.ParseInt(query["skip"], 10, 64)
		if err != nil {
			stdResponse(c, typeInvalidParameter, err)
			return
		}
	}

	var limit int64
	if query["limit"] == "" {
		limit = 0
	} else {
		limit, err = strconv.ParseInt(query["limit"], 10, 64)
		if err != nil {
			stdResponse(c, typeInvalidParameter, err)
			return
		}
	}

	var articleSlice []map[string]interface{}
	articleSlice, err = db.GetRecentArticleByRangeSearch(skip, limit, query["search"])

	if err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}

	response(c, http.StatusOK, "", articleSlice, nil)
}

func deleteArticle(c *gin.Context) {
	query, err := autoQuery(c, map[string]bool{"name": true})
	if err != nil {
		return
	}

	exist, err := db.ExistArticle(query["name"])
	if err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}
	if !exist {
		response(c, http.StatusNotFound, "name not exists", nil, nil)
		return
	}

	if err := db.DeleteArticle(query["name"]); err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}

	response(c, http.StatusOK, "", nil, nil)
}

func updateArticle(c *gin.Context) {
	postForm, err := autoPostForm(c, map[string]bool{"name": true, "updateData": true})
	if err != nil {
		return
	}

	exist, err := db.ExistArticle(postForm["name"])
	if err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}
	if !exist {
		response(c, http.StatusBadRequest, "name not exists", nil, nil)
		return
	}

	var updateData map[string]interface{}
	err = json.Unmarshal([]byte(postForm["updateData"]), &updateData)
	if err != nil {
		stdResponse(c, typeInvalidParameter, err)
		return
	}

	err = db.UpdateArticle(postForm["name"], updateData)
	if err != nil {
		stdResponse(c, typeInvalidParameter, err)
		return
	}

	response(c, http.StatusOK, "", nil, nil)
}

package api

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Zhousiru/inker/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func uploadFile(c *gin.Context) {
	postForm, err := autoPostForm(c, map[string]bool{"name": true})
	if err != nil {
		return
	}

	exist, err := db.ExistFile(postForm["name"])
	if err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}
	if exist {
		response(c, http.StatusBadRequest, "name already exists", nil, nil)
		return
	}

	file, _, err := c.Request.FormFile("upload")
	if err != nil {
		stdResponse(c, typeMissingParameter, err)
	}

	bytesData, err := ioutil.ReadAll(file)
	if err != nil {
		stdResponse(c, typeInternalError, err)
	}

	err = db.UploadFile(postForm["name"], bytesData)
	if err != nil {
		stdResponse(c, typeInternalError, err)
	}

	response(c, http.StatusOK, "", nil, nil)
}

func getFile(c *gin.Context) {
	query, err := autoQuery(c, map[string]bool{"name": true})
	if err != nil {
		return
	}

	bytesData, err := db.GetFile(query["name"])
	if err != nil {
		if err == gridfs.ErrFileNotFound {
			response(c, http.StatusNotFound, "name not exists", nil, err)
			return
		}
		stdResponse(c, typeInternalError, err)
		return
	}

	c.Data(http.StatusOK, "", bytesData)
}

func deleteFile(c *gin.Context) {
	query, err := autoQuery(c, map[string]bool{"name": true})
	if err != nil {
		return
	}

	err = db.DeleteFile(query["name"])
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response(c, http.StatusNotFound, "name not exists", nil, err)
			return
		}
		stdResponse(c, typeInternalError, err)
		return
	}

	response(c, http.StatusOK, "", nil, nil)
}

func paginateFile(c *gin.Context) {
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

	var uploadSlice []map[string]interface{}
	uploadSlice, err = db.GetRecentFileByRange(skip, limit)

	if err != nil {
		stdResponse(c, typeInternalError, err)
		return
	}

	response(c, http.StatusOK, "", uploadSlice, nil)
}

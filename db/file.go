package db

import (
	"bytes"
	"io/ioutil"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getIDByName(name string) (interface{}, error) {
	err := connDB()
	if err != nil {
		return "", err
	}

	uploadColl := getColl("upload.files")
	result := uploadColl.FindOne(nil, bson.M{"filename": name})

	var upload map[string]interface{}
	if err := result.Decode(&upload); err != nil {
		return "", err
	}

	return upload["_id"], nil
}

// ExistFile 检查文件是否存在
func ExistFile(name string) (bool, error) {
	_, err := getIDByName(name)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// UploadFile 上传文件
func UploadFile(name string, bytesData []byte) error {
	err := connDB()
	if err != nil {
		return err
	}

	bucketOpts := new(options.BucketOptions)
	bucketOpts.SetName("upload")

	uploadOpts := new(options.UploadOptions)

	bucket, err := gridfs.NewBucket(getDB(), bucketOpts)
	uploadStream, err := bucket.OpenUploadStream(name, uploadOpts)
	if err != nil {
		return err
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(bytesData)
	if err != nil {
		return err
	}

	return nil
}

// GetFile 获取文件
func GetFile(name string) ([]byte, error) {
	err := connDB()
	if err != nil {
		return nil, err
	}

	bucketOpts := new(options.BucketOptions)
	bucketOpts.SetName("upload")

	bucket, err := gridfs.NewBucket(getDB(), bucketOpts)

	w := bytes.NewBuffer(make([]byte, 0))
	_, err = bucket.DownloadToStreamByName(name, w)
	if err != nil {
		return nil, err
	}

	bytesData, err := ioutil.ReadAll(w)

	return bytesData, nil
}

// DeleteFile 删除文件
func DeleteFile(name string) error {
	err := connDB()
	if err != nil {
		return err
	}

	id, err := getIDByName(name)
	if err != nil {
		return err
	}

	bucketOpts := new(options.BucketOptions)
	bucketOpts.SetName("upload")

	bucket, err := gridfs.NewBucket(getDB(), bucketOpts)

	err = bucket.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// GetRecentFileByRange 按照 uploadDate 降序排列，获取指定范围内的文件
func GetRecentFileByRange(skip int64, limit int64) ([]map[string]interface{}, error) {
	err := connDB()
	if err != nil {
		return nil, err
	}

	uploadColl := getColl("upload.files")

	opt := options.Find()
	opt.SetSort(bson.D{{Key: "uploadDate", Value: -1}})
	opt.SetSkip(skip)
	opt.SetLimit(limit)

	cursor, err := uploadColl.Find(nil, bson.M{}, opt)
	if err != nil {
		return nil, err
	}

	var uploadSlice []map[string]interface{}
	for cursor.Next(nil) {
		var article map[string]interface{}
		if err := cursor.Decode(&article); err != nil {
			return nil, err
		}
		uploadSlice = append(uploadSlice, article)
	}

	return uploadSlice, nil
}

package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ExistArticle 检查文章是否存在
func ExistArticle(name string) (bool, error) {
	err := connDB()
	if err != nil {
		return false, err
	}

	_, err = GetArticle(name)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// NewArticle 创建一篇新文章
func NewArticle(name string, title string, content string, attr map[string]interface{}) error {
	err := connDB()
	if err != nil {
		return err
	}

	articleColl := getColl("articles")

	postDate := time.Now()

	_, err = articleColl.InsertOne(nil, bson.M{
		"name":     name,
		"title":    title,
		"content":  content,
		"attr":     attr,
		"postDate": postDate,
		"editDate": postDate})

	return err
}

// GetArticle 获取一篇文章
func GetArticle(name string) (map[string]interface{}, error) {
	err := connDB()
	if err != nil {
		return nil, err
	}

	articleColl := getColl("articles")
	result := articleColl.FindOne(nil, bson.M{"name": name})

	var article map[string]interface{}
	if err := result.Decode(&article); err != nil {
		return nil, err
	}
	return article, nil
}

// GetRecentArticleByRange 按照 editDate 降序排列，获取指定范围内的文章
func GetRecentArticleByRange(skip int64, limit int64) ([]map[string]interface{}, error) {
	err := connDB()
	if err != nil {
		return nil, err
	}

	articleColl := getColl("articles")

	opt := options.Find()
	opt.SetSort(bson.D{{Key: "editDate", Value: -1}})
	opt.SetSkip(skip)
	opt.SetLimit(limit)

	cursor, err := articleColl.Find(nil, bson.M{}, opt)
	if err != nil {
		return nil, err
	}

	var articleSlice []map[string]interface{}
	for cursor.Next(nil) {
		var article map[string]interface{}
		if err := cursor.Decode(&article); err != nil {
			return nil, err
		}
		articleSlice = append(articleSlice, article)
	}

	return articleSlice, nil
}

// DeleteArticle 删除一篇文章
func DeleteArticle(name string) error {
	err := connDB()
	if err != nil {
		return err
	}

	articleColl := getColl("articles")
	_, err = articleColl.DeleteOne(nil, bson.M{"name": name})

	return err
}

// UpdateArticle 更新一篇文章
func UpdateArticle(name string, updateData map[string]interface{}) error {
	err := connDB()
	if err != nil {
		return err
	}

	articleColl := getColl("articles")

	updateData["editDate"] = time.Now()

	_, err = articleColl.UpdateOne(nil, bson.M{"name": name}, bson.M{"$set": updateData})

	return err
}

// GetRecentArticleByRangeSearch 按照 editDate 降序排列，获取指定范围内符合正则表达式的文章
func GetRecentArticleByRangeSearch(skip int64, limit int64, search string) ([]map[string]interface{}, error) {
	err := connDB()
	if err != nil {
		return nil, err
	}

	articleColl := getColl("articles")

	opt := options.Find()
	opt.SetSort(bson.D{{Key: "editDate", Value: -1}})
	opt.SetSkip(skip)
	opt.SetLimit(limit)

	cursor, err := articleColl.Find(nil, bson.M{"$or": []bson.M{
		bson.M{"content": bson.M{"$regex": search, "$options": "i"}},
		bson.M{"title": bson.M{"$regex": search, "$options": "i"}},
		bson.M{"name": bson.M{"$regex": search, "$options": "i"}},
	}}, opt)
	if err != nil {
		return nil, err
	}

	var articleSlice []map[string]interface{}
	for cursor.Next(nil) {
		var article map[string]interface{}
		if err := cursor.Decode(&article); err != nil {
			return nil, err
		}
		articleSlice = append(articleSlice, article)
	}

	return articleSlice, nil
}

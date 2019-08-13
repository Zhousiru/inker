package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// CheckUsername 检查是否存在指定用户名
func CheckUsername(username string) (bool, error) {
	err := connDB()
	if err != nil {
		return false, err
	}

	userColl := getColl("users")
	result := userColl.FindOne(nil, bson.M{"username": username})

	if _, err := result.DecodeBytes(); err != nil {
		return false, err
	}

	return true, nil
}

// UpdateUser 更新用户数据
func UpdateUser(username string, newUsername string, newPassword string) error {
	err := connDB()
	if err != nil {
		return err
	}

	userColl := getColl("users")
	updateData := bson.M{}
	if newUsername != "" {
		updateData["username"] = newUsername
	}
	if newPassword != "" {
		encrypt, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updateData["password"] = string(encrypt)
	}

	_, err = userColl.UpdateOne(nil, bson.M{"username": username}, bson.M{"$set": updateData})

	return err
}

// NewUser 创建用户
func NewUser(username string, password string) error {
	err := connDB()
	if err != nil {
		return err
	}

	encrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userColl := getColl("users")

	_, err = userColl.InsertOne(nil, bson.M{
		"username": username,
		"password": string(encrypt)})
	if err != nil {
		return err
	}

	return nil
}

// Auth 检查用户名与密码是否匹配
func Auth(username string, password string) (bool, error) {
	err := connDB()
	if err != nil {
		return false, err
	}

	userColl := getColl("users")
	result := userColl.FindOne(nil, bson.M{"username": username})

	user := bson.M{}

	err = result.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

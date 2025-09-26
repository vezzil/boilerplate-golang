package services

import (
    "boilerplate-golang/mysql"
    "boilerplate-golang/src/entity"
)

func GetAllUsers() ([]entity.User, error) {
    var users []entity.User
    result := mysqldb.DB.Find(&users)
    return users, result.Error
}

func CreateUser(user entity.User) (entity.User, error) {
    result := mysqldb.DB.Create(&user)
    return user, result.Error
}

package services

import (
    "go-mvcs-boilerplate/mysqldb"
    "go-mvcs-boilerplate/models"
)

func GetAllUsers() ([]models.User, error) {
    var users []models.User
    result := mysqldb.DB.Find(&users)
    return users, result.Error
}

func CreateUser(user models.User) (models.User, error) {
    result := mysqldb.DB.Create(&user)
    return user, result.Error
}

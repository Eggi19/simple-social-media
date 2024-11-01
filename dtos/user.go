package dtos

import "github.com/Eggi19/simple-social-media/entities"

type UserRegisterData struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,excludes= ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890,containsany=!#$%&'()*+0x2C-./:\"\\;<=>?@[]^_{0x7C}~,min=8,max=128"`
}

func ConvertUserRegisterData(userData *UserRegisterData) *entities.User {
	return &entities.User{
		Name:     userData.Name,
		Email:    userData.Email,
		Password: userData.Password,
	}
}

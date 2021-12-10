package auth

import(
	"auction_auth_service/models"
)


type UserRepository interface{
	CreateUser(user *models.User) error
	GetUserID(email, password string) (string, error)
}
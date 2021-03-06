package authhttp

import (
	"auction_auth_service/auth"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	UseCase auth.UseCase
}

func NewAuthMiddleware(usecase auth.UseCase) *AuthMiddleware {
	return &AuthMiddleware{
		UseCase: usecase,
	}
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	token, err := c.Cookie("access-token")
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, http.ErrNoCookie)
		return
	}

	_, err = m.UseCase.ParseToken(token)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

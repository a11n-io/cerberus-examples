package routes

import (
	"cerberus-examples/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name" `
}

type authRoutes struct {
	userService services.UserService
}

func NewAuthRoutes(userService services.UserService) Routable {
	return &authRoutes{
		userService: userService,
	}
}

func (r *authRoutes) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("auth/register", func(c *gin.Context) { r.Register(c) })
	rg.POST("auth/login", func(c *gin.Context) { r.Login(c) })
}

func (r *authRoutes) Register(c *gin.Context) {
	var authData AuthData

	if err := c.Bind(&authData); err != nil {
		c.AbortWithStatusJSON(400, jsonError(err))
		return
	}

	user, err := r.userService.Register(
		c,
		authData.Email,
		authData.Password,
		authData.Name,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusCreated, jsonData(user))
}

func (r *authRoutes) Login(c *gin.Context) {
	email, password, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("invalid credentials")))
	}

	user, err := r.userService.Login(
		c,
		email,
		password,
	)
	if err != nil {
		c.AbortWithStatusJSON(400, jsonError(err))
		return
	}

	c.JSON(http.StatusCreated, jsonData(user))
}

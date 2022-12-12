package routes

import (
	"cerberus-examples/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	RoleId   string `json:"roleId"`
}

type userRoutes struct {
	userService services.UserService
}

func NewUserRoutes(userService services.UserService) Routable {
	return &userRoutes{userService: userService}
}

func (r *userRoutes) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("users", func(c *gin.Context) { r.Add(c) })
	rg.GET("users", func(c *gin.Context) { r.GetAll(c) })
}

func (r *userRoutes) Add(c *gin.Context) {

	var userData UserData

	if err := c.Bind(&userData); err != nil {
		c.AbortWithStatusJSON(400, jsonError(err))
		return
	}

	user, err := r.userService.Add(
		c,
		userData.Email,
		userData.Password,
		userData.Name,
		userData.RoleId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusCreated, jsonData(user))
}

func (r *userRoutes) GetAll(c *gin.Context) {

	user, err := r.userService.GetAll(
		c,
	)
	if err != nil {
		c.AbortWithStatusJSON(400, jsonError(err))
		return
	}

	c.JSON(http.StatusCreated, jsonData(user))
}

package routes

import (
	"github.com/gin-gonic/gin"
)

type Routable interface {
	RegisterRoutes(r *gin.RouterGroup)
}

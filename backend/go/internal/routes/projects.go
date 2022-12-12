package routes

import (
	"cerberus-examples/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProjectData struct {
	Name        string `json:"name" `
	Description string `json:"description"`
}

type projectRoutes struct {
	service services.ProjectService
}

func NewProjectRoutes(service services.ProjectService) Routable {
	return &projectRoutes{service: service}
}

func (r *projectRoutes) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("accounts/:accountId/projects", func(c *gin.Context) { r.Create(c) })
	rg.GET("accounts/:accountId/projects", func(c *gin.Context) { r.FindAll(c) })
	rg.GET("projects/:projectId", func(c *gin.Context) { r.Get(c) })
	rg.DELETE("projects/:projectId", func(c *gin.Context) { r.Delete(c) })
}

func (r *projectRoutes) Create(c *gin.Context) {

	var projectData ProjectData

	accountId := c.Param("accountId")
	if accountId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing accountId")))
		return
	}

	if err := c.Bind(&projectData); err != nil {
		c.AbortWithStatusJSON(400, jsonError(err))
		return
	}

	project, err := r.service.Create(
		c,
		accountId,
		projectData.Name,
		projectData.Description,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusCreated, jsonData(project))
}

func (r *projectRoutes) FindAll(c *gin.Context) {
	accountId := c.Param("accountId")
	if accountId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing accountId")))
		return
	}

	projects, err := r.service.FindAll(
		c,
		accountId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(projects))
}

func (r *projectRoutes) Get(c *gin.Context) {

	projectId := c.Param("projectId")
	if projectId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing projectId")))
		return
	}

	project, err := r.service.Get(
		c,
		projectId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(project))
}

func (r *projectRoutes) Delete(c *gin.Context) {

	projectId := c.Param("projectId")
	if projectId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing projectId")))
		return
	}

	err := r.service.Delete(
		c,
		projectId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(true))
}

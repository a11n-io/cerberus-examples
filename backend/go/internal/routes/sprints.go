package routes

import (
	"cerberus-examples/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SprintData struct {
	Goal string `json:"goal"`
}

type sprintRoutes struct {
	service services.SprintService
}

func NewSprintRoutes(service services.SprintService) Routable {
	return &sprintRoutes{service: service}
}

func (r *sprintRoutes) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("projects/:projectId/sprints", func(c *gin.Context) { r.Create(c) })
	rg.GET("projects/:projectId/sprints", func(c *gin.Context) { r.FindByProject(c) })
	rg.GET("sprints/:sprintId", func(c *gin.Context) { r.Get(c) })
	rg.POST("sprints/:sprintId/start", func(c *gin.Context) { r.Start(c) })
	rg.POST("sprints/:sprintId/end", func(c *gin.Context) { r.End(c) })
}

func (r *sprintRoutes) Create(c *gin.Context) {

	var resourceTypeData SprintData

	projectId := c.Param("projectId")
	if projectId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing projectId")))
		return
	}

	if err := c.Bind(&resourceTypeData); err != nil {
		c.AbortWithStatusJSON(400, jsonError(err))
		return
	}

	sprint, err := r.service.Create(
		c,
		projectId,
		resourceTypeData.Goal,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusCreated, jsonData(sprint))
}

func (r *sprintRoutes) FindByProject(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.AbortWithStatusJSON(401, jsonError(fmt.Errorf("unauthorized")))
	}

	log.Println("User:", userId)

	projectId := c.Param("projectId")
	if projectId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing projectId")))
		return
	}

	sprints, err := r.service.FindByProject(
		c,
		projectId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(sprints))
}

func (r *sprintRoutes) Start(c *gin.Context) {

	sprintId := c.Param("sprintId")
	if sprintId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing sprintId")))
		return
	}

	rts, err := r.service.Start(
		c,
		sprintId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(rts))
}

func (r *sprintRoutes) End(c *gin.Context) {

	sprintId := c.Param("sprintId")
	if sprintId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing sprintId")))
		return
	}

	rts, err := r.service.End(
		c,
		sprintId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(rts))
}

func (r *sprintRoutes) Get(c *gin.Context) {

	sprintId := c.Param("sprintId")
	if sprintId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing sprintId")))
		return
	}

	sprint, err := r.service.Get(
		c,
		sprintId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(sprint))
}

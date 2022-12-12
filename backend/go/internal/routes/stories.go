package routes

import (
	"cerberus-examples/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type StoryData struct {
	Description string `json:"description"`
	Estimation  string `json:"estimation"`
	Status      string `json:"status"`
	UserId      string `json:"userId"`
}

type storyRoutes struct {
	service services.StoryService
	//cerberusClient cerberus.CerberusClient
}

func NewStoryRoutes(service services.StoryService /*, cerberusClient cerberus.CerberusClient*/) Routable {
	return &storyRoutes{
		service: service,
		//cerberusClient: cerberusClient,
	}
}

func (r *storyRoutes) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("sprints/:sprintId/stories", func(c *gin.Context) { r.Create(c) })
	rg.GET("sprints/:sprintId/stories", func(c *gin.Context) { r.FindBySprint(c) })
	rg.GET("stories/:storyId", func(c *gin.Context) { r.Get(c) })
	rg.POST("stories/:storyId/estimate", func(c *gin.Context) { r.Estimate(c) })
	rg.POST("stories/:storyId/status", func(c *gin.Context) { r.ChangeStatus(c) })
	rg.POST("stories/:storyId/assign", func(c *gin.Context) { r.Assign(c) })
}

func (r *storyRoutes) Create(c *gin.Context) {

	var data StoryData

	sprintId := c.Param("sprintId")
	if sprintId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing sprintId")))
		return
	}

	//hasAccess, err := r.cerberusClient.HasAccess(c, sprintId, common.CreateStory_A)
	//if err != nil || !hasAccess {
	//	c.AbortWithStatusJSON(http.StatusForbidden, jsonError(err))
	//	return
	//}

	if err := c.Bind(&data); err != nil {
		c.AbortWithStatusJSON(400, jsonError(err))
		return
	}

	story, err := r.service.Create(
		c,
		sprintId,
		data.Description,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusCreated, jsonData(story))
}

func (r *storyRoutes) FindBySprint(c *gin.Context) {

	sprintId := c.Param("sprintId")
	if sprintId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing sprintId")))
		return
	}

	stories, err := r.service.FindBySprint(
		c,
		sprintId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(stories))
}

func (r *storyRoutes) Get(c *gin.Context) {

	storyId := c.Param("storyId")
	if storyId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing storyId")))
		return
	}

	//hasAccess, err := r.cerberusClient.HasAccess(c, storyId, common.ReadStory_A)
	//if err != nil || !hasAccess {
	//	c.AbortWithStatusJSON(http.StatusForbidden, jsonError(err))
	//	return
	//}

	story, err := r.service.Get(
		c,
		storyId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(story))
}

func (r *storyRoutes) Estimate(c *gin.Context) {

	storyId := c.Param("storyId")
	if storyId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing storyId")))
		return
	}

	//hasAccess, err := r.cerberusClient.HasAccess(c, storyId, common.EstimateStory_A)
	//if err != nil || !hasAccess {
	//	c.AbortWithStatusJSON(http.StatusForbidden, jsonError(err))
	//	return
	//}

	var data StoryData

	if err := c.Bind(&data); err != nil {
		c.AbortWithStatusJSON(400, jsonError(err))
		return
	}

	estimation, err := strconv.ParseInt(data.Estimation, 10, 32)
	if err != nil {
		estimation = 0
	}

	story, err := r.service.Estimate(
		c,
		storyId,
		int(estimation),
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(story))
}

func (r *storyRoutes) ChangeStatus(c *gin.Context) {

	storyId := c.Param("storyId")
	if storyId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing storyId")))
		return
	}

	//hasAccess, err := r.cerberusClient.HasAccess(c, storyId, common.ChangeStoryStatus_A)
	//if err != nil || !hasAccess {
	//	c.AbortWithStatusJSON(http.StatusForbidden, jsonError(err))
	//	return
	//}

	var data StoryData

	if err := c.Bind(&data); err != nil {
		c.AbortWithStatusJSON(400, jsonError(err))
		return
	}

	story, err := r.service.ChangeStatus(
		c,
		storyId,
		data.Status,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(story))
}

func (r *storyRoutes) Assign(c *gin.Context) {

	storyId := c.Param("storyId")
	if storyId == "" {
		c.AbortWithStatusJSON(400, jsonError(fmt.Errorf("missing storyId")))
		return
	}

	//hasAccess, err := r.cerberusClient.HasAccess(c, storyId, common.ChangeStoryAssignee_A)
	//if err != nil || !hasAccess {
	//	c.AbortWithStatusJSON(http.StatusForbidden, jsonError(err))
	//	return
	//}

	var data StoryData

	if err := c.Bind(&data); err != nil {
		c.AbortWithStatusJSON(400, jsonError(err))
		return
	}

	story, err := r.service.Assign(
		c,
		storyId,
		data.UserId,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, jsonError(err))
		return
	}

	c.JSON(http.StatusOK, jsonData(story))
}

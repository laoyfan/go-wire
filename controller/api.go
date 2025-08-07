package controller

import (
	"encoding/json"
	"go-wire/controller/dto"
	"go-wire/logger"
	"go-wire/service"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

type ApiController struct {
	Controller
	service *service.ApiService
}

func NewApiController(service *service.ApiService, log logger.Logger, trans ut.Translator) *ApiController {
	return &ApiController{
		Controller: Controller{
			log:   log,
			trans: trans,
		},
		service: service,
	}
}

func (c *ApiController) RegisterRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/api")
	userGroup.GET("test", c.Test)
}

func (c *ApiController) Test(ctx *gin.Context) {
	var testReq dto.TestRequest
	if err := c.Valid(ctx, &testReq); err != nil {
		return
	}
	c.InfoLog(ctx, "req", logger.KeyValue("req", testReq))
	value, err := c.service.Test(ctx, testReq.Id)
	if err != nil {
		c.Error(ctx, "获取用户失败", err)
	}
	var testResp dto.TestResponse
	err = json.Unmarshal([]byte(value), &testResp)
	if err != nil {
		c.Error(ctx, "testResp 解析失败", err)
	}
	c.Success(ctx, testResp)
}

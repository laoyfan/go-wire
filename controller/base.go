package controller

import (
	"context"
	"errors"
	"go-wire/constant"
	"go-wire/logger"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	znTranslations "github.com/go-playground/validator/v10/translations/zh"
)

type RouteRegistrar interface {
	RegisterRoutes(group *gin.RouterGroup)
}

type Controller struct {
	trans ut.Translator
	log   logger.Logger
}

func (c *Controller) Result(ctx *gin.Context, code int, msg string, data any) {
	ctx.JSON(http.StatusOK, constant.Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func (c *Controller) Success(ctx *gin.Context, data any) {
	c.Result(ctx, constant.SUCCESS, "请求成功", data)
}

func (c *Controller) Error(ctx *gin.Context, msg string, err error) {
	c.log.Error(ctx, msg, logger.Error(err))
	panic(constant.ErrorResponse{
		Code: constant.ERROR,
		Msg:  msg,
	})
}

func (c *Controller) InfoLog(ctx *gin.Context, msg string, filed ...logger.Field) {
	c.log.Info(ctx, msg, filed...)
}

func (c *Controller) ErrorLog(ctx *gin.Context, msg string, filed ...logger.Field) {
	c.log.Error(ctx, msg, filed...)
}
func NewTrans(log logger.Logger) (ut.Translator, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		uni := ut.New(zhT, zhT)

		Trans, o := uni.GetTranslator("zh")
		if !o {
			log.Error(context.TODO(), "翻译器获取失败")
			return nil, errors.New("翻译器获取失败")
		}

		if err := znTranslations.RegisterDefaultTranslations(v, Trans); err != nil {
			log.Error(context.TODO(), "翻译器注册失败", logger.Error(err))
			return nil, err
		}
		return Trans, nil
	}
	return nil, errors.New("翻译器注册失败")
}

// Valid 参数校验
func (c *Controller) Valid(ctx *gin.Context, valid interface{}) error {
	if err := ctx.ShouldBind(valid); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			c.ErrorLog(ctx, "参数检验失败",
				logger.StringAny("url", ctx.Request.URL.Path),
				logger.StringAny("validationErrors", errs.Translate(c.trans)),
			)
			c.Result(ctx, constant.VALID, "请求参数校验失败", c.removeTopStruct(errs.Translate(c.trans)))
		} else {
			c.ErrorLog(ctx, "请求解析失败",
				logger.StringAny("url", ctx.Request.URL.Path),
				logger.Error(err),
			)
			c.Result(ctx, http.StatusBadRequest, err.Error(), nil)
		}
		return err
	}
	return nil
}

func (c *Controller) removeTopStruct(fields map[string]string) map[string]string {
	res := make(map[string]string, len(fields))
	for field, err := range fields {
		if idx := strings.Index(field, "."); idx != -1 {
			res[field[idx+1:]] = err
		}
	}
	return res
}

package service

import (
	"go-wire/logger"
	"go-wire/repo"

	"github.com/gin-gonic/gin"
)

type ApiService struct {
	repo *repo.ApiRepo
	log  logger.Logger
}

func NewApiService(repo *repo.ApiRepo, log logger.Logger) *ApiService {
	return &ApiService{repo: repo, log: log}
}

func (s *ApiService) Test(ctx *gin.Context, id string) (string, error) {
	user, err := s.repo.Test(ctx, id)
	if err != nil {
		s.log.Error(ctx, "Test Service", logger.Error(err))
		return "", err
	}
	return user, nil
}

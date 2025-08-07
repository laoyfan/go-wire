package dto

type TestRequest struct {
	Id string `json:"id" form:"id" binding:"required"`
}

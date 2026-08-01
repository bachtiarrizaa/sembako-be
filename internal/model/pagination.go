package model

type PaginationRequest struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}

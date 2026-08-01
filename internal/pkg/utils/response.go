package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
)

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalData  int64 `json:"totalData"`
	TotalPages int   `json:"totalPages"`
}

type ApiResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func sendJSON(c *gin.Context, statusCode int, response ApiResponse) {
	bytes, err := json.Marshal(response)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Data(statusCode, "application/json; charset=utf-8", bytes)
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	sendJSON(c, statusCode, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessResponseWithPagination(c *gin.Context, statusCode int, message string, data interface{}, pagination Pagination) {
	sendJSON(c, statusCode, ApiResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: &pagination,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	sendJSON(c, statusCode, ApiResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func HandleError(c *gin.Context, err error) {
	appErr := errs.ToAppError(err)
	ErrorResponse(c, appErr.Code, appErr.Message)
}

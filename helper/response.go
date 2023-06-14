package helper

import (
	"strings"

	"github.com/aldisaputra17/book-store/entities"
)

type Response struct {
	Success bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type ResponseWithPagination struct {
	Success bool                `json:"status"`
	Message string              `json:"message"`
	Errors  interface{}         `json:"errors"`
	Data    interface{}         `json:"data"`
	Total   entities.Pagination `json:"total"`
}

type EmptyObj struct{}

func BuildResponse(success bool, message string, data interface{}) Response {
	res := Response{
		Success: success,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	res := Response{
		Success: false,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
	return res
}

func BuildReadWithPagination(success bool, message string, data interface{}, total entities.Pagination) ResponseWithPagination {
	res := ResponseWithPagination{
		Success: success,
		Message: message,
		Errors:  nil,
		Data:    data,
		Total:   total,
	}
	return res
}

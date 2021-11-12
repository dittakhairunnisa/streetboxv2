package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Pagination ..
type Pagination struct {
	TotalRecords int         `json:"totalRecords"`
	TotalPages   int         `json:"totalPages"`
	Data         interface{} `json:"data"`
	Offset       int         `json:"offset"`
	Limit        int         `json:"limit"`
	Page         int         `json:"page"`
	PrevPage     int         `json:"prevPage"`
	NextPage     int         `json:"nextPage"`
}

// ResponseSuccess model for  swagger
type ResponseSuccess struct {
	Data struct{} `json:"data"`
}

// ResponseSuccessArray ...
type ResponseSuccessArray struct {
	Data []struct{} `json:"data"`
}

// ResponseErrors model for swagger
type ResponseErrors struct {
	Error errorResponse `json:"error"`
}

// ResponseJSON --
func ResponseJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
	return
}

// ResponsePagination ..
func ResponsePagination(c *gin.Context, data Pagination) {
	c.JSON(http.StatusOK, data)
	return
}

// ResponseCreated --
func ResponseCreated(c *gin.Context, message interface{}) {
	c.JSON(http.StatusCreated, message)
	return
}

// ResponseUpdated -- Set response for update process
func ResponseUpdated(c *gin.Context, message interface{}) {
	c.JSON(http.StatusNoContent, message)
	return
}

// ResponseDeleted -- Set response for delete process
func ResponseDeleted(c *gin.Context, message interface{}) {
	if message == "" {
		message = "Resource Deleted"
	}
	c.JSON(http.StatusNoContent, gin.H{"data": message})
	return
}

// ResponseError -- Set response for error
func ResponseError(c *gin.Context, message interface{}, statusCode int) {
	c.JSON(statusCode, gin.H{"error": gin.H{"code": statusCode, "message": message}})
	return
}

// ResponseFailValidation -- Set response for fail validation
func ResponseFailValidation(c *gin.Context, message interface{}) {
	ResponseError(c, message, 422)
	return
}

// ResponseUnauthorized -- Set response not authorized
func ResponseUnauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": errorResponse{
		Code:    http.StatusUnauthorized,
		Message: message,
	}})
	return
}

// ResponseNotFound -- Set response not found
func ResponseNotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource Not Found"
	}
	c.JSON(http.StatusNotFound, gin.H{"error": errorResponse{
		Code:    http.StatusNotFound,
		Message: message,
	}})
	return
}

// ResponseMethodNotAllowed -- Set response method not allowed
func ResponseMethodNotAllowed(c *gin.Context, message string) {
	if message == "" {
		message = "Method Not Allowed"
	}
	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": errorResponse{
		Code:    http.StatusNotFound,
		Message: message,
	}})
	return
}

// ResponseRedirect --
func ResponseRedirect(c *gin.Context, url string) {
	if url == "" {
		return
	}
	c.Redirect(http.StatusFound, url)
	return
}

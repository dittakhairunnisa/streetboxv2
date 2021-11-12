package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"streetbox.id/model"
	"streetbox.id/util"
)

// AuthMiddleware ...
//
// Permission:
//
// 1. all		: all authorized user
//
// 2. superadmin: superadmin
//
// 3. merchant	: foodtruck -> admin -> superadmin
//
// 4. admin		: admin -> superadmin
//
// 5. consumer	: consumer
func AuthMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := util.TokenValid(c.Request, permission)
		if err != nil {
			model.ResponseError(c, err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}

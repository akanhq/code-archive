// middleware/casbin_rbac.go
package middleware

import (
	"casbin_demo/pkg/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("roleID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		obj := c.FullPath()     // e.g. /products/add
		act := c.Request.Method // e.g. POST

		ok, err := casbin.Enforcer.Enforce(roleID, obj, act)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission"})
			c.Abort()
			return
		}

		c.Next()
	}
}

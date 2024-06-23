package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Role string

const (
	Admin      Role = "admin"
	Instructor Role = "instructor"
	Student    Role = "student"
)

var rolePermissions = map[Role][]string{
	Admin:      {"*"},
	Instructor: {"/courses", "/lessons", "/availability"},
	Student:    {"/courses", "/lessons"},
}

func CheckRolePermission(role Role, path string) bool {
	permissions, exists := rolePermissions[role]
	if !exists {
		return false
	}

	for _, perm := range permissions {
		if perm == "*" || strings.HasPrefix(path, perm) {
			return true
		}
	}
	return false
}

func RBACMiddleware(roles ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		role, ok := userRole.(Role)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		for _, r := range roles {
			if role == r && CheckRolePermission(role, c.Request.URL.Path) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}

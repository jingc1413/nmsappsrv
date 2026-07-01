package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NullStrategy is a pass-through authentication strategy that always
// allows requests without any credential validation.
// Aligned with Java NullAuthentication.
type NullStrategy struct{}

func (s *NullStrategy) Name() string { return "Null" }

func (s *NullStrategy) Authenticate(c *gin.Context, _, _ string) bool {
	return true
}

func (s *NullStrategy) Challenge(c *gin.Context) {
	// Null strategy never challenges
	c.Status(http.StatusUnauthorized)
	c.Abort()
}

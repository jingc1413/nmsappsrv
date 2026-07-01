package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BasicStrategy implements HTTP Basic Authentication for TR-069 CPE.
// Aligned with Java BasicAuthentication.
type BasicStrategy struct {
	Realm string
}

func (s *BasicStrategy) Name() string { return "Basic" }

func (s *BasicStrategy) Authenticate(c *gin.Context, expectedUser, expectedPass string) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(authHeader[6:])
	if err != nil {
		return false
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return false
	}

	return parts[0] == expectedUser && parts[1] == expectedPass
}

func (s *BasicStrategy) Challenge(c *gin.Context) {
	c.Header("WWW-Authenticate", `Basic realm="`+s.Realm+`"`)
	c.Status(http.StatusUnauthorized)
	c.Abort()
}

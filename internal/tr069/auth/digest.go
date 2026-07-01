package auth

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// DigestStrategy implements HTTP Digest Authentication (RFC 2617) for TR-069 CPE.
// Aligned with Java DigestAuthentication — but with actual MD5 hashing
// (Java version has a bug where MD5Utils.toMD5() is a no-op).
type DigestStrategy struct {
	Realm string
}

func (s *DigestStrategy) Name() string { return "Digest" }

func (s *DigestStrategy) Authenticate(c *gin.Context, expectedUser, expectedPass string) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Digest ") {
		return false
	}

	directives := parseDigestHeader(authHeader[7:])

	username := directives["username"]
	realm := directives["realm"]
	nonce := directives["nonce"]
	uri := directives["uri"]
	response := directives["response"]
	qop := directives["qop"]
	nc := directives["nc"]
	cnonce := directives["cnonce"]

	if username == "" || realm == "" || nonce == "" || uri == "" || response == "" {
		return false
	}

	if username != expectedUser {
		return false
	}

	// Compute expected digest per RFC 2617
	A1 := fmt.Sprintf("%s:%s:%s", username, realm, expectedPass)
	HA1 := md5Hex(A1)

	A2 := fmt.Sprintf("POST:%s", uri)
	HA2 := md5Hex(A2)

	var expected string
	if qop == "auth" || qop == "auth-int" {
		expected = md5Hex(fmt.Sprintf("%s:%s:%s:%s:%s:%s", HA1, nonce, nc, cnonce, qop, HA2))
	} else {
		expected = md5Hex(fmt.Sprintf("%s:%s:%s", HA1, nonce, HA2))
	}

	return expected == response
}

func (s *DigestStrategy) Challenge(c *gin.Context) {
	nonce := md5Hex(fmt.Sprintf("%d:%s", time.Now().UnixNano(), s.Realm))
	header := fmt.Sprintf(
		`Digest realm="%s", nonce="%s", algorithm=MD5, qop="auth"`,
		s.Realm, nonce,
	)
	c.Header("WWW-Authenticate", header)
	c.Status(http.StatusUnauthorized)
	c.Abort()
}

// parseDigestHeader parses a Digest authorization header value into key-value pairs.
func parseDigestHeader(header string) map[string]string {
	result := make(map[string]string)
	parts := strings.Split(header, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		eqIdx := strings.Index(part, "=")
		if eqIdx < 0 {
			continue
		}
		key := strings.TrimSpace(part[:eqIdx])
		val := strings.TrimSpace(part[eqIdx+1:])
		val = strings.Trim(val, `"`)
		result[key] = val
	}
	return result
}

// md5Hex computes MD5 hash and returns hex string.
func md5Hex(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

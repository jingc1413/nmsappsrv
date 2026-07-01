package auth

import "github.com/gin-gonic/gin"

// Strategy defines the authentication strategy for TR-069 CPE connections.
// Aligned with Java Authentication interface.
type Strategy interface {
	// Authenticate validates the request. Returns true if authenticated.
	Authenticate(c *gin.Context, deviceUsername, devicePassword string) bool
	// Challenge sends a 401 response with the appropriate WWW-Authenticate header.
	Challenge(c *gin.Context)
	// Name returns the strategy name (Basic, Digest, Null).
	Name() string
}

// Get returns the auth strategy matching the given algorithm name.
// Aligned with Java AuthenticationHandler dispatch logic:
//   - "Basic"  → BasicStrategy
//   - "Digest" → DigestStrategy
//   - default  → NullStrategy (pass-through)
//
// Note: Java's "MD5" strategy is a stub (always returns null), so it's omitted.
func Get(algorithm string) Strategy {
	switch algorithm {
	case "Basic":
		return &BasicStrategy{Realm: "TR-069 ACS"}
	case "Digest":
		return &DigestStrategy{Realm: "TR-069 ACS"}
	default:
		return &NullStrategy{}
	}
}

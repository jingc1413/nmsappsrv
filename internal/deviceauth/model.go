package deviceauth

// DeviceAuthConfig represents the per-tenant device authentication configuration.
// Aligned with Java ACSAuthenticationDTO.
type DeviceAuthConfig struct {
	Algorithm string `json:"algorithm"` // "Basic", "Digest", "Null"
	Enable    bool   `json:"enable"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

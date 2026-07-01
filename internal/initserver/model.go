// Package initserver implements the Zero-Touch Provisioning (ZTP) init server.
//
// When a gNB device boots for the first time, it contacts this init server
// via TR-069 CWMP. The server responds with a SetParameterValues (SPV) message
// containing IPsec, CA, and ACS configuration so the device can establish
// secure tunnels and register with the production ACS.
//
// Configuration is persisted in the system_config table (key = "initserver")
// as a JSON blob, following the same pattern as other modules (ntp, site, etc.).
package initserver

// SystemConfig mirrors the system_config table for key-value config storage.
// The initserver config is stored with id="initserver" and config=<JSON>.
type SystemConfig struct {
	Id     string  `gorm:"primaryKey;column:id;type:varchar(255)" json:"id"`
	Config *string `gorm:"column:config;type:longtext" json:"config"`
}

func (SystemConfig) TableName() string { return "system_config" }

// InitServerConfig holds all init-server configuration parameters.
// Field names and JSON tags mirror the Java ConfigVO for compatibility.
//
// The config covers three functional areas:
//   - CA Server: URL, username, password for certificate authority
//   - ACS Server: management server URL for device registration
//   - IPsec: full IKE/IPsec tunnel parameters (gateway, identity, crypto, ChildSA)
type InitServerConfig struct {
	// --- Switch ---
	// Enable or disable the init server ("Enable" / "Disable").
	InitServerEnable string `json:"initServerEnable"`

	// --- CA Server ---
	CAUrl      string `json:"caUrl"`
	CAUsername string `json:"caUsername"`
	CAPassword string `json:"caPassword"`

	// --- ACS ---
	ACSURL string `json:"acsURL"`

	// --- IPsec General ---
	IPsecEnable              string `json:"ipSecEnable"`
	IPsecAuthenticationMethod string `json:"ipSecAuthenticationMethod"`
	IPsecPreSharedKey        string `json:"ipSecPreSharedKey"`
	IPsecCerts               string `json:"ipSecCerts"`

	// --- IPsec Gateway ---
	IPsecSecGWServer1 string `json:"ipSecSecGWServer1"`
	IPsecSecGWServer2 string `json:"ipSecSecGWServer2"`
	IPsecSecGWServer3 string `json:"ipSecSecGWServer3"`

	// --- IPsec Identity ---
	IPsecLocalId  string `json:"ipSecLocalId"`
	IPsecRemoteId string `json:"ipSecRemoteId"`

	// --- IPsec Ports ---
	IPsecLocalPort     string `json:"ipSecLocalPort"`
	IPsecLocalNattPort string `json:"ipSecLocalNattPort"`
	IPsecRemotePort    string `json:"ipSecRemotePort"`

	// --- IPsec EAP ---
	IPsecLocalEapId  string `json:"ipSecLocalEapId"`
	IPsecRemoteEapId string `json:"ipSecRemoteEapId"`

	// --- IPsec Crypto ---
	IPsecOPC                          string `json:"ipSecOPC"`
	IPsecK                            string `json:"ipSecK"`
	IPsecEncryptionAlgorithms         string `json:"ipSecEncryptionAlgorithms"`
	IPsecIntegrityAlgorithms          string `json:"ipSecIntegrityAlgorithms"`
	IPsecDiffieHellmanGroupTransforms string `json:"ipSecDiffieHellmanGroupTransforms"`

	// --- IPsec VIPS ---
	IPsecEnableVips   string `json:"ipSecEnableVips"`
	IPsecEnableVipsV6 string `json:"ipSecEnableVipsV6"`

	// --- IPsec DPD ---
	IPsecDpdDelay string `json:"ipSecDpdDelay"`

	// --- IPsec ChildSA ---
	IPsecId                              string `json:"ipSecId"`
	IPsecLocalTs                         string `json:"ipSecLocalTs"`
	IPsecRemoteTs                        string `json:"ipSecRemoteTs"`
	IPsecChildSAEncryptionAlgorithms     string `json:"ipSecChildSAEncryptionAlgorithms"`
	IPsecChildSAIntegrityAlgorithms      string `json:"ipSecChildSAIntegrityAlgorithms"`
	IPsecChildSADiffieHellmanGroupTransforms string `json:"ipSecChildSADiffieHellmanGroupTransforms"`
}

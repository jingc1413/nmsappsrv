package tr069

// Device type constants matching the Java NMS convention.
const (
	DeviceTypeCPE = "cpe"
	DeviceTypeENB = "enb" // used for both eNB and gNB
	DeviceTypeNR  = "nr"  // alias for gNB
	DeviceTypeLTE = "lte" // alias for eNB
)

// IGD = InternetGatewayDevice (root path for CPE)
const igd = "InternetGatewayDevice"

// GetBasicParamPaths returns the list of TR-069 parameter paths to query
// during device initialization, based on the device type.
// These correspond to the Java GetParametersPostprocessor basic param lists.
func GetBasicParamPaths(deviceType string) []string {
	switch deviceType {
	case DeviceTypeCPE:
		return cpeBasicParamPaths
	case DeviceTypeENB, DeviceTypeNR:
		return nrBasicParamPaths
	case DeviceTypeLTE:
		return lteBasicParamPaths
	default:
		// Default to CPE paths
		return cpeBasicParamPaths
	}
}

// cpeBasicParamPaths are the basic parameter paths queried during CPE device initialization.
var cpeBasicParamPaths = []string{
	// Device Info
	igd + ".DeviceInfo.Manufacturer",
	igd + ".DeviceInfo.ModelName",
	igd + ".DeviceInfo.ProductClass",
	igd + ".DeviceInfo.SerialNumber",
	igd + ".DeviceInfo.HardwareVersion",
	igd + ".DeviceInfo.SoftwareVersion",
	igd + ".DeviceInfo.AdditionalHardwareVersion",
	igd + ".DeviceInfo.AdditionalSoftwareVersion",
	igd + ".DeviceInfo.ProvisioningCode",
	igd + ".DeviceInfo.SpecVersion",
	igd + ".DeviceInfo.UpTime",
	igd + ".DeviceInfo.FirstUseDate",
	igd + ".DeviceInfo.X_BAIKEL_Temperature",
	// Management Server
	igd + ".ManagementServer.ConnectionRequestURL",
	igd + ".ManagementServer.ConnectionRequestUsername",
	igd + ".ManagementServer.PeriodicInformInterval",
	igd + ".ManagementServer.ParameterKey",
	igd + ".ManagementServer.PeriodicInformTime",
	// WAN Device
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.MACAddress",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.DefaultGateway",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.Name",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.NATEnabled",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.AddressingType",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.IPAddress",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.SubnetMask",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.Gateway",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.DNSServers",
	// WAN Ethernet
	igd + ".WANDevice.1.WANConnectionDevice.1.WANEthernetLinkConfig.MACAddress",
	// LAN
	igd + ".LANDevice.1.LANEthernetInterfaceConfig.1.MACAddress",
	igd + ".LANDevice.1.LANEthernetInterfaceConfig.1.Status",
	igd + ".LANDevice.1.LANEthernetInterfaceConfig.1.MaxBitRate",
	// Hosts
	igd + ".Hosts.HostNumberOfEntries",
	// STB
	igd + ".STBService.1.Component.VideoDecoder.Status",
	// Services
	igd + ".Services.VoiceService.1.VoiceProfile.1.Enable",
	igd + ".Services.VoiceService.1.VoiceProfile.1.Line.1.Enable",
	igd + ".Services.VoiceService.1.VoiceProfile.1.Line.1.Status",
	// Firewall
	igd + ".X_BAIKEL_Firewall.Enable",
	igd + ".X_BAIKEL_Firewall.Level",
	// IP
	igd + ".Layer3Forwarding.DefaultConnectionService",
}

// nrBasicParamPaths are the basic parameter paths queried during NR (gNB) device initialization.
var nrBasicParamPaths = []string{
	// Device Info
	igd + ".DeviceInfo.Manufacturer",
	igd + ".DeviceInfo.ModelName",
	igd + ".DeviceInfo.ProductClass",
	igd + ".DeviceInfo.SerialNumber",
	igd + ".DeviceInfo.HardwareVersion",
	igd + ".DeviceInfo.SoftwareVersion",
	igd + ".DeviceInfo.AdditionalHardwareVersion",
	igd + ".DeviceInfo.AdditionalSoftwareVersion",
	igd + ".DeviceInfo.ProvisioningCode",
	igd + ".DeviceInfo.SpecVersion",
	igd + ".DeviceInfo.UpTime",
	igd + ".DeviceInfo.FirstUseDate",
	igd + ".DeviceInfo.X_BAIKEL_Temperature",
	// Management Server
	igd + ".ManagementServer.ConnectionRequestURL",
	igd + ".ManagementServer.ConnectionRequestUsername",
	igd + ".ManagementServer.PeriodicInformInterval",
	igd + ".ManagementServer.ParameterKey",
	// WAN
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.MACAddress",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.DefaultGateway",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.Name",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.NATEnabled",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.AddressingType",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.IPAddress",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.SubnetMask",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.Gateway",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.DNSServers",
	// NR-specific: Radio
	igd + ".Radio.1.Status",
	igd + ".Radio.1.OperatingFrequencyBand",
	igd + ".Radio.1.OperatingChannelBandwidth",
	igd + ".Radio.1.TxPower",
	// NR-specific: Cell
	igd + ".CellConfig.1.PhyCellId",
	igd + ".CellConfig.1.CellIdentity",
	igd + ".CellConfig.1.TAC",
	igd + ".CellConfig.1.MCC",
	igd + ".CellConfig.1.MNC",
	igd + ".CellConfig.1.NSSAI",
	// NR-specific: DU/CU
	igd + ".DUConfig.1.Status",
	igd + ".CUConfig.1.Status",
	// LAN
	igd + ".LANDevice.1.LANEthernetInterfaceConfig.1.MACAddress",
	igd + ".LANDevice.1.LANEthernetInterfaceConfig.1.Status",
}

// lteBasicParamPaths are the basic parameter paths queried during LTE (eNB) device initialization.
var lteBasicParamPaths = []string{
	// Device Info
	igd + ".DeviceInfo.Manufacturer",
	igd + ".DeviceInfo.ModelName",
	igd + ".DeviceInfo.ProductClass",
	igd + ".DeviceInfo.SerialNumber",
	igd + ".DeviceInfo.HardwareVersion",
	igd + ".DeviceInfo.SoftwareVersion",
	igd + ".DeviceInfo.AdditionalHardwareVersion",
	igd + ".DeviceInfo.AdditionalSoftwareVersion",
	igd + ".DeviceInfo.ProvisioningCode",
	igd + ".DeviceInfo.SpecVersion",
	igd + ".DeviceInfo.UpTime",
	igd + ".DeviceInfo.FirstUseDate",
	igd + ".DeviceInfo.X_BAIKEL_Temperature",
	// Management Server
	igd + ".ManagementServer.ConnectionRequestURL",
	igd + ".ManagementServer.ConnectionRequestUsername",
	igd + ".ManagementServer.PeriodicInformInterval",
	igd + ".ManagementServer.ParameterKey",
	// WAN
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.MACAddress",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.DefaultGateway",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.Name",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.NATEnabled",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.AddressingType",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.IPAddress",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.SubnetMask",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.Gateway",
	igd + ".WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.DNSServers",
	// LTE-specific: Radio
	igd + ".Radio.1.Status",
	igd + ".Radio.1.OperatingFrequencyBand",
	igd + ".Radio.1.OperatingChannelBandwidth",
	igd + ".Radio.1.TxPower",
	igd + ".Radio.1.EARFCN",
	// LTE-specific: Cell
	igd + ".CellConfig.1.PhyCellId",
	igd + ".CellConfig.1.CellIdentity",
	igd + ".CellConfig.1.TAC",
	igd + ".CellConfig.1.MCC",
	igd + ".CellConfig.1.MNC",
	// LAN
	igd + ".LANDevice.1.LANEthernetInterfaceConfig.1.MACAddress",
	igd + ".LANDevice.1.LANEthernetInterfaceConfig.1.Status",
}

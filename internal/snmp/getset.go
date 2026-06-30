package snmp

import (
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
	"nmsappsrv/pkg/logger"
)

// SendGet performs an SNMP GET request for the given OIDs and returns the results.
func SendGet(connInfo SnmpConnectionInfo, oids []string) ([]SnmpParameter, error) {
	snmpClient := buildGoSNMP(connInfo)

	if err := snmpClient.Connect(); err != nil {
		return nil, fmt.Errorf("SNMP connect failed to %s:%d: %w", connInfo.IP, connInfo.Port, err)
	}
	defer snmpClient.Conn.Close()

	// Ensure OIDs start with a dot
	prefixedOIDs := make([]string, len(oids))
	for i, oid := range oids {
		if len(oid) > 0 && oid[0] != '.' {
			prefixedOIDs[i] = "." + oid
		} else {
			prefixedOIDs[i] = oid
		}
	}

	result, err := snmpClient.Get(prefixedOIDs)
	if err != nil {
		return nil, fmt.Errorf("SNMP GET failed to %s:%d: %w", connInfo.IP, connInfo.Port, err)
	}

	var params []SnmpParameter
	for _, variable := range result.Variables {
		params = append(params, SnmpParameter{
			OID:   variable.Name,
			Type:  gosnmpTypeToName(variable.Type),
			Value: formatSnmpValue(variable),
		})
	}

	logger.Infof("SNMP GET completed for %s:%d, got %d results", connInfo.IP, connInfo.Port, len(params))
	return params, nil
}

// SendSet performs an SNMP SET request with the given parameters.
func SendSet(connInfo SnmpConnectionInfo, params []SnmpParameter) error {
	snmpClient := buildGoSNMP(connInfo)

	if err := snmpClient.Connect(); err != nil {
		return fmt.Errorf("SNMP connect failed to %s:%d: %w", connInfo.IP, connInfo.Port, err)
	}
	defer snmpClient.Conn.Close()

	pdus, err := buildTrapVariables(params)
	if err != nil {
		return fmt.Errorf("failed to build SET PDUs: %w", err)
	}

	result, err := snmpClient.Set(pdus)
	if err != nil {
		return fmt.Errorf("SNMP SET failed to %s:%d: %w", connInfo.IP, connInfo.Port, err)
	}

	if len(result.Variables) > 0 {
		for _, v := range result.Variables {
			if v.Type == gosnmp.NoSuchObject || v.Type == gosnmp.NoSuchInstance || v.Type == gosnmp.Null {
				logger.Warnf("SNMP SET returned error for OID %s: type=%v", v.Name, v.Type)
			}
		}
	}

	logger.Infof("SNMP SET completed for %s:%d, %d variables", connInfo.IP, connInfo.Port, len(pdus))
	return nil
}

// buildGoSNMP creates a gosnmp.GoSNMP instance from connection info (shared by GET/SET/WALK).
func buildGoSNMP(connInfo SnmpConnectionInfo) *gosnmp.GoSNMP {
	snmpClient := &gosnmp.GoSNMP{
		Target:    connInfo.IP,
		Port:      uint16(connInfo.Port),
		Transport: "udp",
		Timeout:   5 * time.Second,
		Retries:   3,
	}

	switch connInfo.Version {
	case 1:
		snmpClient.Version = gosnmp.Version1
		snmpClient.Community = connInfo.Community
	case 2:
		snmpClient.Version = gosnmp.Version2c
		snmpClient.Community = connInfo.Community
	case 3:
		snmpClient.Version = gosnmp.Version3
		snmpClient.SecurityModel = gosnmp.UserSecurityModel
		snmpClient.MsgFlags = getMsgFlags(connInfo.SecurityLevel)

		usmParams := &gosnmp.UsmSecurityParameters{
			UserName:                 connInfo.Username,
			AuthenticationPassphrase: connInfo.Password,
			PrivacyPassphrase:        connInfo.PrivacyPassphrase,
			AuthenticationProtocol:   getAuthProtocol(connInfo.AuthenticationProtocol),
			PrivacyProtocol:          getPrivProtocol(connInfo.PrivacyProtocol),
			AuthoritativeEngineID:    connInfo.EngineId,
		}
		snmpClient.SecurityParameters = usmParams
	default:
		snmpClient.Version = gosnmp.Version2c
		snmpClient.Community = connInfo.Community
	}

	return snmpClient
}

// gosnmpTypeToName converts a gosnmp ASN type to a human-readable type name.
func gosnmpTypeToName(t gosnmp.Asn1BER) string {
	switch t {
	case gosnmp.OctetString:
		return "string"
	case gosnmp.Integer:
		return "int32"
	case gosnmp.Counter32:
		return "uint32"
	case gosnmp.Gauge32:
		return "gauge32"
	case gosnmp.Counter64:
		return "counter64"
	case gosnmp.ObjectIdentifier:
		return "oid"
	case gosnmp.TimeTicks:
		return "timeTicks"
	case gosnmp.IPAddress:
		return "ipv4"
	default:
		return "unknown"
	}
}

// formatSnmpValue converts a gosnmp variable value to a string representation.
func formatSnmpValue(variable gosnmp.SnmpPDU) string {
	switch variable.Type {
	case gosnmp.OctetString:
		if b, ok := variable.Value.([]byte); ok {
			return string(b)
		}
		return fmt.Sprintf("%v", variable.Value)
	case gosnmp.Integer:
		return fmt.Sprintf("%d", variable.Value)
	case gosnmp.Counter32, gosnmp.Gauge32, gosnmp.TimeTicks:
		return fmt.Sprintf("%d", variable.Value)
	case gosnmp.Counter64:
		return fmt.Sprintf("%d", variable.Value)
	case gosnmp.ObjectIdentifier:
		return fmt.Sprintf("%s", variable.Value)
	case gosnmp.IPAddress:
		return fmt.Sprintf("%s", variable.Value)
	default:
		return fmt.Sprintf("%v", variable.Value)
	}
}

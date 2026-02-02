package scanner

import (
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
)

// SNMPClient represents an SNMP connection to a network device
type SNMPClient struct {
	client *gosnmp.GoSNMP
}

// SNMPConfig holds SNMP connection configuration
type SNMPConfig struct {
	Host      string
	Port      uint16
	Version   string // 2c or 3
	Community string // For v2c
	// SNMPv3 parameters
	Username  string
	AuthPass  string
	PrivPass  string
	Timeout   time.Duration
}

// NewSNMPClient creates a new SNMP client
func NewSNMPClient(config SNMPConfig) (*SNMPClient, error) {
	client := &gosnmp.GoSNMP{
		Target:    config.Host,
		Port:      config.Port,
		Timeout:   config.Timeout,
		Retries:   2,
		MaxOids:   60,
	}

	switch config.Version {
	case "2c":
		client.Version = gosnmp.Version2c
		client.Community = config.Community
	case "3":
		client.Version = gosnmp.Version3
		client.SecurityModel = gosnmp.UserSecurityModel
		client.MsgFlags = gosnmp.AuthPriv
		client.SecurityParameters = &gosnmp.UsmSecurityParameters{
			UserName:                 config.Username,
			AuthenticationProtocol:   gosnmp.SHA,
			AuthenticationPassphrase: config.AuthPass,
			PrivacyProtocol:          gosnmp.AES,
			PrivacyPassphrase:        config.PrivPass,
		}
	default:
		return nil, fmt.Errorf("unsupported SNMP version: %s", config.Version)
	}

	err := client.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &SNMPClient{client: client}, nil
}

// GetSystemInfo retrieves basic system information
func (c *SNMPClient) GetSystemInfo() (map[string]string, error) {
	info := make(map[string]string)

	oids := map[string]string{
		"sysDescr":    "1.3.6.1.2.1.1.1.0",
		"sysName":     "1.3.6.1.2.1.1.5.0",
		"sysLocation": "1.3.6.1.2.1.1.6.0",
		"sysContact":  "1.3.6.1.2.1.1.4.0",
	}

	for name, oid := range oids {
		result, err := c.client.Get([]string{oid})
		if err != nil {
			continue
		}
		if len(result.Variables) > 0 {
			switch v := result.Variables[0].Value.(type) {
			case string:
				info[name] = v
			case []byte:
				info[name] = string(v)
			}
		}
	}

	return info, nil
}

// GetInterfaceInfo retrieves interface information using IF-MIB
func (c *SNMPClient) GetInterfaceInfo() ([]map[string]interface{}, error) {
	var interfaces []map[string]interface{}

	// Walk ifDescr (1.3.6.1.2.1.2.2.1.2)
	ifDescrOID := "1.3.6.1.2.1.2.2.1.2"
	err := c.client.Walk(ifDescrOID, func(pdu gosnmp.SnmpPDU) error {
		iface := make(map[string]interface{})
		
		switch v := pdu.Value.(type) {
		case string:
			iface["name"] = v
		case []byte:
			iface["name"] = string(v)
		}
		
		interfaces = append(interfaces, iface)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return interfaces, nil
}

// GetLLDPNeighbors retrieves LLDP neighbor information
func (c *SNMPClient) GetLLDPNeighbors() ([]map[string]string, error) {
	var neighbors []map[string]string

	// LLDP-MIB remTable OID: 1.0.8802.1.1.2.1.4.1
	lldpRemTableOID := "1.0.8802.1.1.2.1.4.1"
	
	err := c.client.Walk(lldpRemTableOID, func(pdu gosnmp.SnmpPDU) error {
		neighbor := make(map[string]string)
		neighbor["oid"] = pdu.Name
		
		switch v := pdu.Value.(type) {
		case string:
			neighbor["value"] = v
		case []byte:
			neighbor["value"] = string(v)
		}
		
		neighbors = append(neighbors, neighbor)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return neighbors, nil
}

// Close closes the SNMP connection
func (c *SNMPClient) Close() error {
	if c.client != nil {
		return c.client.Conn.Close()
	}
	return nil
}

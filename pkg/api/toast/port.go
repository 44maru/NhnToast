package toast

import (
	"encoding/json"
	"log"

	"nhn-toast/pkg/constants"
	"nhn-toast/pkg/infrastructure/http"
)

type PortListResponse struct {
	Ports []struct {
		Status              string        `json:"status"`
		Name                string        `json:"name"`
		AllowedAddressPairs []interface{} `json:"allowed_address_pairs"`
		AdminStateUp        bool          `json:"admin_state_up"`
		NetworkID           string        `json:"network_id"`
		TenantID            string        `json:"tenant_id"`
		ExtraDhcpOpts       []interface{} `json:"extra_dhcp_opts"`
		BindingVnicType     string        `json:"binding:vnic_type"`
		DeviceOwner         string        `json:"device_owner"`
		MacAddress          string        `json:"mac_address"`
		PortSecurityEnabled bool          `json:"port_security_enabled"`
		FixedIps            []struct {
			SubnetID  string `json:"subnet_id"`
			IPAddress string `json:"ip_address"`
		} `json:"fixed_ips"`
		ID             string   `json:"id"`
		SecurityGroups []string `json:"security_groups"`
		DeviceID       string   `json:"device_id"`
	} `json:"ports"`
}

func GetPortList(token string) (*PortListResponse, error) {
	httpReqHeader := map[string]string{}
	httpReqHeader["X-Auth-Token"] = token
	jsonRes, err := http.Get(constants.PORT_URL, httpReqHeader, nil)
	if err != nil {
		return nil, err
	}

	portList := new(PortListResponse)
	err = json.Unmarshal(jsonRes, &portList)
	if err != nil {
		log.Println("Port list response json unmarshal err")
		return nil, err
	}

	return portList, nil
}

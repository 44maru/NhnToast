package toast

import (
	"encoding/json"
	"fmt"
	"log"

	"nhn-toast/pkg/constants"
	"nhn-toast/pkg/infrastructure/http"
)

type GetSubnetResponse struct {
	Subnets []struct {
		Name            string        `json:"name"`
		EnableDhcp      bool          `json:"enable_dhcp"`
		NetworkID       string        `json:"network_id"`
		TenantID        string        `json:"tenant_id"`
		DNSNameservers  []interface{} `json:"dns_nameservers"`
		GatewayIP       string        `json:"gateway_ip"`
		Ipv6RaMode      interface{}   `json:"ipv6_ra_mode"`
		AllocationPools []struct {
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"allocation_pools"`
		HostRoutes      []interface{} `json:"host_routes"`
		IPVersion       int           `json:"ip_version"`
		Ipv6AddressMode interface{}   `json:"ipv6_address_mode"`
		Cidr            string        `json:"cidr"`
		ID              string        `json:"id"`
		SubnetpoolID    interface{}   `json:"subnetpool_id"`
	} `json:"subnets"`
}
type GetNetworkResponse struct {
	Networks []struct {
		Name                string   `json:"name"`
		ID                  string   `json:"id"`
		Status              string   `json:"status"`
		Shared              bool     `json:"shared"`
		Subnets             []string `json:"subnets"`
		AdminStateUp        bool     `json:"admin_state_up"`
		PortSecurityEnabled bool     `json:"port_security_enabled"`
		RouterExternal      bool     `json:"router:external"`
		TenantID            string   `json:"tenant_id"`
		Mtu                 int      `json:"mtu"`
	} `json:"networks"`
}

func GetNetworkInfo(networkName, token string) (*GetNetworkResponse, error) {
	httpReqHeader := map[string]string{}
	httpReqHeader["X-Auth-Token"] = token
	queryParam := map[string]string{}
	queryParam["name"] = networkName
	jsonRes, err := http.Get(constants.NETWORK_URL, httpReqHeader, queryParam)
	if err != nil {
		return nil, err
	}

	networkInfoList := new(GetNetworkResponse)
	err = json.Unmarshal(jsonRes, &networkInfoList)
	if err != nil {
		log.Println("Get network info response json unmarshal err")
		return nil, err
	}

	return networkInfoList, nil
}

func GetNetworkId(networkName, token string) (string, error) {
	networkInfoList, err := GetNetworkInfo(networkName, token)
	if err != nil {
		return "", err
	}

	if len(networkInfoList.Networks) < 1 {
		return "", fmt.Errorf("Not found network id for '%s'\n", networkName)
	}

	return networkInfoList.Networks[0].ID, nil
}

func GetSubnetIdList(networkName, token string) ([]string, error) {
	networkInfoList, err := GetNetworkInfo(networkName, token)
	if err != nil {
		return nil, err
	}

	if len(networkInfoList.Networks) < 1 {
		return nil, fmt.Errorf("Not found network id for '%s'\n", networkName)
	}

	var subnetIdList []string

	for _, subnetInfo := range networkInfoList.Networks[0].Subnets {
		subnetIdList = append(subnetIdList, subnetInfo)
	}

	return subnetIdList, nil
}

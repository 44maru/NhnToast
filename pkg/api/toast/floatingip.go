package toast

import (
	"encoding/json"
	"fmt"
	"log"
	"nhn-toast/pkg/config"
	"nhn-toast/pkg/constants"
	"nhn-toast/pkg/infrastructure/http"
	"sync"
	"time"
)

type CreateFloatingIpRequest struct {
	Floatingip struct {
		FloatingNetworkID string `json:"floating_network_id"`
		PortID            string `json:"port_id"`
	} `json:"floatingip"`
}

type JointFloatingIpRequest struct {
	Floatingip struct {
		PortID string `json:"port_id"`
	} `json:"floatingip"`
}

type CreateFloatingIpResponse struct {
	Floatingip struct {
		FloatingNetworkID string `json:"floating_network_id"`
		RouterID          string `json:"router_id"`
		FixedIPAddress    string `json:"fixed_ip_address"`
		FloatingIPAddress string `json:"floating_ip_address"`
		TenantID          string `json:"tenant_id"`
		Status            string `json:"status"`
		PortID            string `json:"port_id"`
		ID                string `json:"id"`
	} `json:"floatingip"`
}

func CreateFloatingIps(serverList *ServerListDetailResponse, config *config.Config, publicNetworkId, token string) error {
	isSuccess := true
	portList, err := GetPortList(token)
	if err != nil {
		log.Println("ERROR : Failed to get port list")
		return err
	}

	var wg sync.WaitGroup
	limitCh := make(chan struct{}, config.Thread.ThreadNum)
	for _, serverInfo := range serverList.Servers {
		globalIpMacAddrMap := getGlobalIpMacAddrMap(serverInfo)

		for _, vpcInfo := range serverInfo.Addresses.DefaultNetwork {
			if vpcInfo.OSEXTIPSType == constants.OS_EXT_IP_TYPE_FLOATING {
				continue
			}

			if _, doesExist := globalIpMacAddrMap[vpcInfo.OSEXTIPSMACMacAddr]; doesExist {
				log.Printf("%s mac address '%s' already associated floating ip. skip create and joint floating ip.\n", serverInfo.Name, vpcInfo.OSEXTIPSMACMacAddr)
				continue
			}

			portId := getPortId(vpcInfo.OSEXTIPSMACMacAddr, portList)
			if portId == "" {
				isSuccess = false
				log.Printf("Not found port id for %s. skip create and joint floating ip.\n", serverInfo.Name)
				continue
			}

			wg.Add(1)
			limitCh <- struct{}{}
			go func(serverName, portId string) {
				defer wg.Done()
				createRes, err := createFloatingIp(portId, publicNetworkId, token)
				if err != nil {
					log.Printf("ERROR : instance '%s' failed to create floating ip. %s\n", serverName, err.Error())
					isSuccess = false
				} else {
					time.Sleep(time.Second * config.Thread.SleepSecBeforeJointFloatingIp)
					_, err := jointFloatingIp(createRes.Floatingip.ID, portId, token)
					if err != nil {
						log.Printf("ERROR : instance '%s' failed to joint floating ip. %s\n", serverName, err.Error())
						isSuccess = false
					} else {
						log.Printf("instance '%s' successed to create and joint floating ip.\n", serverName)
					}
				}
				<-limitCh
			}(serverInfo.Name, portId)
		}
	}

	wg.Wait()

	if isSuccess {
		return nil
	} else {
		return fmt.Errorf("Failed to create and joint any of floating ip.")
	}
}

func createFloatingIp(portId, publicNetworkId, token string) (*CreateFloatingIpResponse, error) {
	createReq := new(CreateFloatingIpRequest)
	createReq.Floatingip.FloatingNetworkID = publicNetworkId
	createReq.Floatingip.PortID = portId
	reqJsonBytes, err := json.MarshalIndent(createReq, "", "  ")
	if err != nil {
		log.Println("Request json marshal error")
		return nil, err
	}

	httpReqHeader := map[string]string{}
	httpReqHeader["Content-Type"] = "application/json"
	httpReqHeader["X-Auth-Token"] = token
	resJsonBytes, err := http.Post(constants.FLOATING_IP_URL, reqJsonBytes, httpReqHeader)
	if err != nil {
		log.Printf("ERROR : FloatingIP作成リクエスト失敗\n%s\n", err.Error())
		return nil, err
	}

	createRes := new(CreateFloatingIpResponse)
	err = json.Unmarshal(resJsonBytes, &createRes)
	if err != nil {
		log.Println("Create floating ip response json unmarshal err")
		return nil, err
	}

	return createRes, nil
}

func jointFloatingIp(floatingIpId, portId, token string) (*CreateFloatingIpResponse, error) {
	jointReq := new(JointFloatingIpRequest)
	jointReq.Floatingip.PortID = portId
	reqJsonBytes, err := json.MarshalIndent(jointReq, "", "  ")
	if err != nil {
		log.Println("Request json marshal error")
		return nil, err
	}

	httpReqHeader := map[string]string{}
	httpReqHeader["Content-Type"] = "application/json"
	httpReqHeader["X-Auth-Token"] = token
	resJsonBytes, err := http.Put(constants.FLOATING_IP_URL+"/"+floatingIpId, reqJsonBytes, httpReqHeader)
	if err != nil {
		log.Printf("ERROR : FloatingIP接続リクエスト失敗\n%s\n", err.Error())
		return nil, err
	}

	jointRes := new(CreateFloatingIpResponse)
	err = json.Unmarshal(resJsonBytes, &jointRes)
	if err != nil {
		log.Println("Joint floating ip response json unmarshal err")
		return nil, err
	}

	return jointRes, nil
}

func getPortId(macAddress string, portList *PortListResponse) string {
	for _, portInfo := range portList.Ports {
		if macAddress == portInfo.MacAddress {
			return portInfo.ID
		}
	}

	return ""
}

func DeleteFloatingIpList(floatinIpList *FloatingIpListResponse, config *config.Config, token string) error {
	isSuccess := true
	var wg sync.WaitGroup
	limitCh := make(chan struct{}, config.Thread.ThreadNum)

	httpReqHeader := map[string]string{}
	httpReqHeader["Content-Type"] = "application/json"
	httpReqHeader["X-Auth-Token"] = token

	for _, floatingIpInfo := range floatinIpList.Floatingips {
		deleteInstanceUrl := constants.FLOATING_IP_URL + "/" + floatingIpInfo.ID
		wg.Add(1)
		limitCh <- struct{}{}
		go func(floatingIpAddress, deleteFloatingIpUrl string) {
			defer wg.Done()
			err := http.Delete(deleteFloatingIpUrl, httpReqHeader)
			if err != nil {
				log.Printf("ERROR : floating ip '%s' failed to delete. %s\n", floatingIpAddress, err.Error())
				isSuccess = false
			} else {
				log.Printf("floating ip '%s' successed to delete.\n", floatingIpAddress)
			}
			time.Sleep(time.Second * config.Thread.SleepSecondsAfterDeleteFloatingIp)
			<-limitCh
		}(floatingIpInfo.FloatingIPAddress, deleteInstanceUrl)
	}

	wg.Wait()

	if isSuccess {
		return nil
	} else {
		return fmt.Errorf("Failed to delete any of floating ip.")
	}
}

type FloatingIpListResponse struct {
	Floatingips []struct {
		FloatingNetworkID string      `json:"floating_network_id"`
		RouterID          interface{} `json:"router_id"`
		FixedIPAddress    interface{} `json:"fixed_ip_address"`
		FloatingIPAddress string      `json:"floating_ip_address"`
		TenantID          string      `json:"tenant_id"`
		Status            string      `json:"status"`
		PortID            interface{} `json:"port_id"`
		ID                string      `json:"id"`
	} `json:"floatingips"`
}

func GetFloatingIpList(config *config.Config, token string) (*FloatingIpListResponse, error) {
	httpReqHeader := map[string]string{}
	httpReqHeader["X-Auth-Token"] = token
	jsonRes, err := http.Get(constants.FLOATING_IP_URL, httpReqHeader, nil)
	if err != nil {
		return nil, err
	}

	floatingipList := new(FloatingIpListResponse)
	err = json.Unmarshal(jsonRes, &floatingipList)
	if err != nil {
		log.Println("Floating ip list response json unmarshal err")
		return nil, err
	}

	return floatingipList, nil
}

func getGlobalIpMacAddrMap(serverInfo Server) map[string]struct{} {
	globalIpMacAddrMap := make(map[string]struct{})
	for _, vpcInfo := range serverInfo.Addresses.DefaultNetwork {
		if vpcInfo.OSEXTIPSType == constants.OS_EXT_IP_TYPE_FLOATING {
			globalIpMacAddrMap[vpcInfo.OSEXTIPSMACMacAddr] = struct{}{}
		}
	}

	return globalIpMacAddrMap
}

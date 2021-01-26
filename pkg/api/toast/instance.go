package toast

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"nhn-toast/pkg/config"
	"nhn-toast/pkg/constants"
	"nhn-toast/pkg/infrastructure/http"

	"github.com/tealeg/xlsx"
)

const CREATE_INSTANCE_INFO_EXCEL_PATH = "./create-instance-list.xlsx"
const U2_C1M1_FLAVER_REF = "b41750b4-d819-487d-84bc-89fc7a6d0df1"

type CreateInstanceInfo struct {
	Name          string
	SecurityGroup string
	KeyPairName   string
	NumOfInstance int
}

type ServerCreateRequest struct {
	Server struct {
		Name      string    `json:"name"`
		ImageRef  string    `json:"imageRef"`
		FlavorRef string    `json:"flavorRef"`
		Networks  []Network `json:"networks"`
		//AvailabilityZone     string                 `json:"availability_zone`
		KeyName  string `json:"key_name"`
		MaxCount int    `json:"max_count"`
		MinCount int    `json:"min_count"`
		//BlockDeviceMappingV2 []BlockDeviceMappingV2 `json:"block_device_mapping_v2"`
		SecurityGroups []SecurityGroup `json:"security_groups"`
	} `json:"server"`
}

type ServerStopRequest struct {
	OsStop []struct{} `json:"os-stop"`
}

type AutoGenerated struct {
	Server struct {
	} `json:"server"`
}

type BlockDeviceMappingV2 struct {
	UUID                string `json:"uuid"`
	BootIndex           int    `json:"boot_index"`
	VolumeSize          int    `json:"volume_size"`
	DeviceName          string `json:"device_name"`
	SourceType          string `json:"source_type"`
	DestinationType     string `json:"destination_type"`
	DeleteOnTermination int    `json:"delete_on_termination"`
}

type SecurityGroup struct {
	Name string `json:"name"`
}
type Network struct {
	Subnet string `json:"subnet"`
}

type ServerCreateResponse struct {
	Server struct {
		SecurityGroups  []SecurityGroup `json:"name"`
		OSDCFDiskConfig string          `json:"OS-DCF:diskConfig"`
		ID              string          `json:"id"`
		Links           []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	} `json:"server"`
}

type ServerBootRequest struct {
	OsStart interface{} `json:"os-start"`
}

type ServerShutdownRequest struct {
	OsStop interface{} `json:"os-stop"`
}

func CreateInstance(config *config.Config, subnetIdList []string, token string) error {
	isSuccess := true
	instanceCreateUrl := constants.COMPUTE_ENDPOINT + "/v2/" + config.UserInfo.TenantId + "/servers"

	createInstanceInfoList, err := getCreateInstanceInfoList(CREATE_INSTANCE_INFO_EXCEL_PATH)
	if err != nil {
		return err
	}

	networks := createNetworks(subnetIdList)

	var wg sync.WaitGroup
	limitCh := make(chan struct{}, config.Thread.ThreadNum)
	for _, instanceInfo := range createInstanceInfoList {
		wg.Add(1)
		limitCh <- struct{}{}
		go func(instanceInfo *CreateInstanceInfo) {
			defer wg.Done()
			err := createInstance(instanceInfo, config.Instance.ImageId, networks, token, instanceCreateUrl)
			if err != nil {
				log.Printf("ERROR : instance '%s' failed to create. %s\n", instanceInfo.Name, err.Error())
				isSuccess = false
			} else {
				log.Printf("instance '%s' successed to create.\n", instanceInfo.Name)
			}
			<-limitCh
		}(instanceInfo)
	}

	wg.Wait()

	if isSuccess {
		return nil
	} else {
		return fmt.Errorf("Failed to delete any of instance.")
	}
}

func createInstance(instanceInfo *CreateInstanceInfo, proxyImageId string, networks []Network, token, instanceCreateUrl string) error {
	createReq := new(ServerCreateRequest)
	createReq.Server.Name = instanceInfo.Name
	createReq.Server.ImageRef = proxyImageId
	createReq.Server.FlavorRef = U2_C1M1_FLAVER_REF
	createReq.Server.SecurityGroups = []SecurityGroup{{Name: instanceInfo.SecurityGroup}}
	createReq.Server.Networks = networks
	createReq.Server.KeyName = instanceInfo.KeyPairName
	createReq.Server.MinCount = instanceInfo.NumOfInstance
	createReq.Server.MaxCount = instanceInfo.NumOfInstance
	reqJsonBytes, err := json.MarshalIndent(createReq, "", "  ")
	if err != nil {
		log.Println("Request json marshal error")
		return err
	}

	httpReqHeader := map[string]string{}
	httpReqHeader["Content-Type"] = "application/json"
	httpReqHeader["X-Auth-Token"] = token
	_, err = http.Post(instanceCreateUrl, reqJsonBytes, httpReqHeader)
	if err != nil {
		log.Printf("インスタンス作成エラー\n%s\n", err.Error())
		return err
	}

	return nil
}

func getCreateInstanceInfoList(excelFilePath string) ([]*CreateInstanceInfo, error) {
	excel, err := xlsx.OpenFile(excelFilePath)
	if err != nil {
		log.Printf("%sのオープンに失敗", excelFilePath)
		return nil, err
	}

	var createInstanceInfoList []*CreateInstanceInfo
	sheet := excel.Sheets[0]
	for i, row := range sheet.Rows {
		if i == 0 {
			continue
		}

		name := row.Cells[0].String()
		if name == "" {
			log.Printf("%d行目、インスタンス名が空のためスキップ", i+1)
			continue
		}

		securityGroup := row.Cells[1].String()
		if securityGroup == "" {
			log.Printf("%d行目、セキュリティグループが空のためスキップ", i+1)
			continue
		}

		keyPairName := row.Cells[2].String()
		if keyPairName == "" {
			log.Printf("%d行目、キーペア名が空のためスキップ", i+1)
			continue
		}

		numOfInstance, err := row.Cells[3].Int()
		if err != nil {
			log.Println("作成インスタンス数取得エラー")
			return nil, err
		}

		createInstanceInfo := &CreateInstanceInfo{
			Name:          name,
			SecurityGroup: securityGroup,
			KeyPairName:   keyPairName,
			NumOfInstance: numOfInstance,
		}
		createInstanceInfoList = append(createInstanceInfoList, createInstanceInfo)
	}

	return createInstanceInfoList, nil
}

func DeleteInstanceList(serverInfoList *ServerListDetailResponse, config *config.Config, token string) error {
	isSuccess := true
	var wg sync.WaitGroup
	limitCh := make(chan struct{}, config.Thread.ThreadNum)

	httpReqHeader := map[string]string{}
	httpReqHeader["Content-Type"] = "application/json"
	httpReqHeader["X-Auth-Token"] = token

	for _, serverInfo := range serverInfoList.Servers {
		deleteInstanceUrl := constants.COMPUTE_ENDPOINT + "/v2/" + config.UserInfo.TenantId + "/servers/" + serverInfo.ID
		wg.Add(1)
		limitCh <- struct{}{}
		go func(serverName, deleteInstanceUrl string) {
			defer wg.Done()
			err := http.Delete(deleteInstanceUrl, httpReqHeader)
			if err != nil {
				log.Printf("ERROR : instance '%s' failed to delete. %s\n", serverName, err.Error())
				isSuccess = false
			} else {
				log.Printf("instance '%s' successed to delete.\n", serverName)
			}
			time.Sleep(time.Second * config.Thread.SleepSecondsAfterDeleteInstance)
			<-limitCh
		}(serverInfo.Name, deleteInstanceUrl)
	}

	wg.Wait()

	if isSuccess {
		return nil
	} else {
		return fmt.Errorf("Failed to delete any of instance.")
	}

}

func StartInstanceList(serverInfoList *ServerListDetailResponse, config *config.Config, token string) error {
	isSuccess := true
	var wg sync.WaitGroup

	bootRequest := new(ServerBootRequest)
	reqJsonBytes, err := json.MarshalIndent(bootRequest, "", "  ")

	httpReqHeader := map[string]string{}
	httpReqHeader["Content-Type"] = "application/json"
	httpReqHeader["X-Auth-Token"] = token

	limitCh := make(chan struct{}, config.Thread.ThreadNum)
	for _, serverInfo := range serverInfoList.Servers {
		instanceUrl := constants.COMPUTE_ENDPOINT + "/v2/" + config.UserInfo.TenantId + "/servers/" + serverInfo.ID + "/action"
		wg.Add(1)
		limitCh <- struct{}{}
		go func(serverName, instanceUrl string) {
			defer wg.Done()
			_, err = http.Post(instanceUrl, reqJsonBytes, httpReqHeader)
			if err != nil {
				log.Printf("ERROR : instance '%s' failed to start up. %s\n", serverName, err.Error())
				isSuccess = false
			} else {
				log.Printf("instance '%s' successed to start up.\n", serverName)
			}
			time.Sleep(time.Second * config.Thread.SleepSecondsAfterDeleteInstance)
			<-limitCh
		}(serverInfo.Name, instanceUrl)
	}

	wg.Wait()

	if isSuccess {
		return nil
	} else {
		return fmt.Errorf("Failed to start up any of instance.")
	}

}

type ServerListDetailResponse struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Status    string    `json:"status"`
	Updated   time.Time `json:"updated"`
	HostID    string    `json:"hostId"`
	Addresses struct {
		DefaultNetwork []struct {
			OSEXTIPSMACMacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
			Version            int    `json:"version"`
			Addr               string `json:"addr"`
			OSEXTIPSType       string `json:"OS-EXT-IPS:type"`
		} `json:"Default Network"`
	} `json:"addresses"`
	Links []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	KeyName string `json:"key_name"`
	Image   struct {
		ID    string `json:"id"`
		Links []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	} `json:"image"`
	OSEXTSTSTaskState  interface{} `json:"OS-EXT-STS:task_state"`
	OSEXTSTSVMState    string      `json:"OS-EXT-STS:vm_state"`
	OSSRVUSGLaunchedAt string      `json:"OS-SRV-USG:launched_at"`
	Flavor             struct {
		ID    string `json:"id"`
		Links []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	} `json:"flavor"`
	ID             string `json:"id"`
	SecurityGroups []struct {
		Name string `json:"name"`
	} `json:"security_groups"`
	OSSRVUSGTerminatedAt             interface{} `json:"OS-SRV-USG:terminated_at"`
	OSEXTAZAvailabilityZone          string      `json:"OS-EXT-AZ:availability_zone"`
	UserID                           string      `json:"user_id"`
	Name                             string      `json:"name"`
	Created                          time.Time   `json:"created"`
	TenantID                         string      `json:"tenant_id"`
	OSDCFDiskConfig                  string      `json:"OS-DCF:diskConfig"`
	OsExtendedVolumesVolumesAttached []struct {
		ID string `json:"id"`
	} `json:"os-extended-volumes:volumes_attached"`
	AccessIPv4         string `json:"accessIPv4"`
	AccessIPv6         string `json:"accessIPv6"`
	Progress           int    `json:"progress"`
	OSEXTSTSPowerState int    `json:"OS-EXT-STS:power_state"`
	ConfigDrive        string `json:"config_drive"`
	Metadata           struct {
		OsDistro        string `json:"os_distro"`
		Description     string `json:"description"`
		OsVersion       string `json:"os_version"`
		ProjectDomain   string `json:"project_domain"`
		HypervisorType  string `json:"hypervisor_type"`
		MonitoringAgent string `json:"monitoring_agent"`
		ImageName       string `json:"image_name"`
		VolumeSize      string `json:"volume_size"`
		OsArchitecture  string `json:"os_architecture"`
		LoginUsername   string `json:"login_username"`
		OsType          string `json:"os_type"`
		TcEnv           string `json:"tc_env"`
	} `json:"metadata"`
}

func GetInstanceListDetail(config *config.Config, token string) (*ServerListDetailResponse, error) {
	instanceListDetailUrl := constants.COMPUTE_ENDPOINT + "/v2/" + config.UserInfo.TenantId + "/servers/detail"
	httpReqHeader := map[string]string{}
	httpReqHeader["Content-Type"] = "application/json"
	httpReqHeader["X-Auth-Token"] = token
	jsonRes, err := http.Get(instanceListDetailUrl, httpReqHeader, nil)
	if err != nil {
		return nil, err
	}

	serverList := new(ServerListDetailResponse)
	err = json.Unmarshal(jsonRes, &serverList)
	if err != nil {
		log.Println("Server list response json unmarshal err")
		return nil, err
	}

	/*
			var buf bytes.Buffer
			err = json.Indent(&buf, jsonRes, "", "  ")
			if err != nil {
				println("Response JSON format error.")
				return nil, err
			}

		fmt.Println(buf.String())
	*/

	return serverList, nil
}

func createNetworks(subnetIdList []string) []Network {
	var networks []Network

	for _, subnetId := range subnetIdList {
		network := Network{Subnet: subnetId}
		networks = append(networks, network)
	}

	return networks
}

func StopInstanceList(serverInfoList *ServerListDetailResponse, config *config.Config, token string) error {
	isSuccess := true
	var wg sync.WaitGroup

	shutdownRequest := new(ServerShutdownRequest)
	reqJsonBytes, err := json.MarshalIndent(shutdownRequest, "", "  ")

	httpReqHeader := map[string]string{}
	httpReqHeader["Content-Type"] = "application/json"
	httpReqHeader["X-Auth-Token"] = token

	limitCh := make(chan struct{}, config.Thread.ThreadNum)
	for _, serverInfo := range serverInfoList.Servers {
		instanceUrl := constants.COMPUTE_ENDPOINT + "/v2/" + config.UserInfo.TenantId + "/servers/" + serverInfo.ID + "/action"
		wg.Add(1)
		limitCh <- struct{}{}
		go func(serverName, instanceUrl string) {
			defer wg.Done()
			_, err = http.Post(instanceUrl, reqJsonBytes, httpReqHeader)
			if err != nil {
				log.Printf("ERROR : instance '%s' failed to shutdown. %s\n", serverName, err.Error())
				isSuccess = false
			} else {
				log.Printf("instance '%s' successed to shutdown.\n", serverName)
			}
			time.Sleep(time.Second * config.Thread.SleepSecondsAfterDeleteInstance)
			<-limitCh
		}(serverInfo.Name, instanceUrl)
	}

	wg.Wait()

	if isSuccess {
		return nil
	} else {
		return fmt.Errorf("Failed to shutdown any of instance.")
	}

}

package toast

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"nhn-toast/pkg/constants"
	"nhn-toast/pkg/infrastructure/http"
)

type GetImageResponse struct {
	Images []struct {
		ContainerFormat string        `json:"container_format"`
		MinRAM          int           `json:"min_ram"`
		UpdatedAt       time.Time     `json:"updated_at"`
		LoginUsername   string        `json:"login_username"`
		File            string        `json:"file"`
		Owner           string        `json:"owner"`
		ID              string        `json:"id"`
		Size            int           `json:"size"`
		OsDistro        string        `json:"os_distro"`
		Self            string        `json:"self"`
		DiskFormat      string        `json:"disk_format"`
		OsVersion       string        `json:"os_version"`
		Schema          string        `json:"schema"`
		Status          string        `json:"status"`
		Description     string        `json:"description"`
		Tags            []interface{} `json:"tags"`
		Visibility      string        `json:"visibility"`
		OsArchitecture  string        `json:"os_architecture"`
		MinDisk         int           `json:"min_disk"`
		VirtualSize     interface{}   `json:"virtual_size"`
		Name            string        `json:"name"`
		HypervisorType  string        `json:"hypervisor_type"`
		CreatedAt       time.Time     `json:"created_at"`
		Protected       bool          `json:"protected"`
		Checksum        string        `json:"checksum"`
		OsType          string        `json:"os_type"`
	} `json:"images"`
	Schema string `json:"schema"`
	First  string `json:"first"`
	Next   string `json:"next"`
}

func GetImageId(imageName, token string) (string, error) {
	httpReqHeader := map[string]string{}
	httpReqHeader["X-Auth-Token"] = token
	queryParam := map[string]string{}
	queryParam["name"] = imageName
	jsonRes, err := http.Get(constants.IMAGE_URL, httpReqHeader, queryParam)
	if err != nil {
		return "", err
	}

	imageInfoList := new(GetImageResponse)
	err = json.Unmarshal(jsonRes, &imageInfoList)
	if err != nil {
		log.Println("Get image list response json unmarshal err")
		return "", err
	}

	if len(imageInfoList.Images) < 1 {
		return "", fmt.Errorf("Not found image id for '%s'\n", imageName)
	}

	return imageInfoList.Images[0].ID, nil
}

func GetImageList(token string) (*GetImageResponse, error) {
	httpReqHeader := map[string]string{}
	httpReqHeader["X-Auth-Token"] = token
	jsonRes, err := http.Get(constants.IMAGE_URL, httpReqHeader, nil)
	if err != nil {
		return nil, err
	}

	imageInfoList := new(GetImageResponse)
	err = json.Unmarshal(jsonRes, &imageInfoList)
	if err != nil {
		log.Println("Get image list response json unmarshal err")
		return nil, err
	}

	return imageInfoList, nil
}

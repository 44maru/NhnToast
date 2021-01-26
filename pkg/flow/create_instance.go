package flow

import (
	"log"

	"nhn-toast/pkg/api/toast"
	"nhn-toast/pkg/config"
	"nhn-toast/pkg/constants"
)

func CreateInstance(config *config.Config, token string) {
	/*
		imageId, err := toast.GetImageId(config.Instance.ImageName, token)
		if err != nil {
			log.Printf("ERROR: イメージ '%s'のID取得失敗: %s\n", config.Instance.ImageName, err.Error())
			return
		}
	*/

	subnetIdList, err := toast.GetSubnetIdList(constants.DEFAULT_SUBNET_NAME, token)
	if err != nil {
		log.Printf("ERROR: ネットワーク '%s'のID取得失敗: %s\n", constants.DEFAULT_SUBNET_NAME, err.Error())
		return
	}

	err = toast.CreateInstance(config, subnetIdList, token)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return
	}

	log.Println("インスタンス作成処理成功")
}

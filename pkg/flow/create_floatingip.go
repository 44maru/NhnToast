package flow

import (
	"log"

	"nhn-toast/pkg/api/toast"
	"nhn-toast/pkg/config"
	"nhn-toast/pkg/constants"
)

const LogFile = "./info.log"

func CreateFloatingIp(config *config.Config, token string) {
	publicNetworkId, err := toast.GetNetworkId(constants.PUBLIC_NETWORK_NAME, token)
	if err != nil {
		log.Printf("ERROR: ネットワーク '%s'のID取得失敗: %s\n", constants.PUBLIC_NETWORK_NAME, err.Error())
		return
	}

	instanceList, err := toast.GetInstanceListDetail(config, token)
	if err != nil {
		log.Printf("ERROR: インスタンス情報取得処理でエラーが発生しました。: %s\n", err.Error())
		return
	}

	err = toast.CreateFloatingIps(instanceList, config, publicNetworkId, token)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return
	}

	log.Println("FloatingIP作成処理成功")
}

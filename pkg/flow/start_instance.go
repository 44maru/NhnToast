package flow

import (
	"log"

	"nhn-toast/pkg/api/toast"
	"nhn-toast/pkg/config"
)

func StartInstance(config *config.Config, token string) {
	instanceList, err := toast.GetInstanceListDetail(config, token)
	if err != nil {
		log.Printf("ERROR: インスタンス情報取得処理でエラーが発生しました。: %s\n", err.Error())
		return
	}

	err = toast.StartInstanceList(instanceList, config, token)
	if err != nil {
		log.Printf("ERROR: 1台以上のインスタンス起動に失敗しました。: %s\n", err.Error())
		return
	}

	log.Println("インスタンス全台の起動に成功しました")
}

package flow

import (
	"log"

	"nhn-toast/pkg/api/toast"
	"nhn-toast/pkg/config"
)

func DumpInstanceGlobalIpList(config *config.Config, token string) {
	instanceList, err := toast.GetInstanceListDetail(config, token)
	if err != nil {
		log.Printf("ERROR: インスタンス情報取得処理でエラーが発生しました。: %s\n", err.Error())
		return
	}

	err = toast.DumpGloabalIPList(instanceList)
	if err != nil {
		log.Printf("ERROR: インスタンGlobalIP出力に失敗しました。: %s\n", err.Error())
		return
	}

	log.Println("インスタンス全台のGlobalIPの出力に成功しました")
}

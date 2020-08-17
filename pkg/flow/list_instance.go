package flow

import (
	"encoding/json"
	"log"

	"nhn-toast/pkg/api/toast"
	"nhn-toast/pkg/config"
)

func ListInstance(config *config.Config, token string) {
	instanceList, err := toast.GetInstanceListDetail(config, token)
	if err != nil {
		log.Printf("ERROR: インスタンス情報取得処理でエラーが発生しました。: %s\n", err.Error())
		return
	}

	instanceListJson, err := json.MarshalIndent(instanceList, "", "  ")
	if err != nil {
		log.Printf("Marshal err")
		return
	}

	log.Println(string(instanceListJson))
	log.Printf("合計%d台のインスタンス情報を取得しました\n", len(instanceList.Servers))
}

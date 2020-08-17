package flow

import (
	"log"

	"nhn-toast/pkg/api/toast"
	"nhn-toast/pkg/config"
)

func DeleteInstance(config *config.Config, token string) {
	deleteInstanceList, err := toast.GetInstanceListDetail(config, token)
	if err != nil {
		log.Printf("ERROR: インスタンス情報取得処理でエラーが発生しました。: %s\n", err.Error())
		return
	}

	err = toast.DeleteInstanceList(deleteInstanceList, config, token)
	if err != nil {
		log.Printf("ERROR: 1台以上のインスタンス削除に失敗しました。: %s\n", err.Error())
		return
	}

	log.Println("インスタンス全台の削除に成功しました")
}

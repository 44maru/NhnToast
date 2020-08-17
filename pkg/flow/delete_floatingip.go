package flow

import (
	"log"

	"nhn-toast/pkg/api/toast"
	"nhn-toast/pkg/config"
)

func DeleteFloatingIp(config *config.Config, token string) {
	deleteFloatingIpList, err := toast.GetFloatingIpList(config, token)
	if err != nil {
		log.Printf("ERROR: FloatingIP情報取得処理でエラーが発生しました。: %s\n", err.Error())
		return
	}

	err = toast.DeleteFloatingIpList(deleteFloatingIpList, config, token)
	if err != nil {
		log.Printf("ERROR: 1台以上のインスタンス削除に失敗しました。: %s\n", err.Error())
		return
	}

	log.Println("全FloatingIPの削除に成功しました")
}

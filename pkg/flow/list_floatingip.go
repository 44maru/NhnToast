package flow

import (
	"encoding/json"
	"log"

	"nhn-toast/pkg/api/toast"
	"nhn-toast/pkg/config"
)

func ListFloatingIp(config *config.Config, token string) {
	floatingIpList, err := toast.GetFloatingIpList(config, token)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return
	}

	floatingIpListJson, err := json.MarshalIndent(floatingIpList, "", "  ")
	if err != nil {
		log.Printf("Marshal err")
		return
	}

	log.Println(string(floatingIpListJson))
	log.Printf("合計%d個のFloatingIP情報を取得しました\n", len(floatingIpList.Floatingips))
}

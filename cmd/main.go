package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"nhn-toast/pkg/api/toast"
	"nhn-toast/pkg/config"
	"nhn-toast/pkg/constants"
	"nhn-toast/pkg/flow"
	"nhn-toast/pkg/infrastructure/util"
)

const LogFile = "./info.log"

func main() {
	flowType := flag.String("flow", "create-instance", "flow type")
	flag.Parse()

	logfile, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		util.FailOnError(fmt.Sprintf("%sのオープンに失敗しました", LogFile), err)
	}
	defer logfile.Close()
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)

	config, err := config.LoadConfig()
	if err != nil {
		util.FailOnError("コンフィグファイルのロードに失敗しました", err)
	}

	log.Println("Token取得処理開始...")
	token, err := toast.GenerateToken(config.UserInfo.TenantId, config.UserInfo.UserName, config.UserInfo.ApiPassword)
	if err != nil {
		log.Printf("Token取得エラー : %s\n", err.Error())
		util.WaitEnter()
		return
	}
	log.Printf("Token取得成功 : %s\n", token)

	if *flowType == constants.FLOW_TYPE_CREATE_INSTANCE {
		flow.CreateInstance(config, token)

	} else if *flowType == constants.FLOW_TYPE_CREATE_FLOATINGIP {
		flow.CreateFloatingIp(config, token)

	} else if *flowType == constants.FLOW_TYPE_DELETE_INSTANCE {
		flow.DeleteInstance(config, token)

	} else if *flowType == constants.FLOW_TYPE_DELETE_FLOATINGIP {
		flow.DeleteFloatingIp(config, token)

	} else if *flowType == constants.FLOW_TYPE_LIST_INSTANCE {
		flow.ListInstance(config, token)

	} else if *flowType == constants.FLOW_TYPE_LIST_FLOATINGIP {
		flow.ListFloatingIp(config, token)

	} else if *flowType == constants.FLOW_TYPE_LIST_IMAGE {
		flow.ListImage(token)

	} else if *flowType == constants.FLOW_TYPE_START_INSTANCE {
		flow.StartInstance(config, token)

	} else if *flowType == constants.FLOW_TYPE_STOP_INSTANCE {
		flow.StopInstance(config, token)

	} else if *flowType == constants.FLOW_TYPE_DUMP_GLOBAL_IP_LIST {
		flow.DumpInstanceGlobalIpList(config, token)

	}

	util.WaitEnter()
}

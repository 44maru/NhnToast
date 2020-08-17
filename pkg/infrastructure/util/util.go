package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const LogFile = "./info.log"
const ConfigFile = "./config.json"

func FailOnError(errMsg string, err error) {
	//errs := errors.WithStack(err)
	log.Println(errMsg)
	if err != nil {
		log.Printf("%s\n", err.Error())
	}
	WaitEnter()
	os.Exit(1)
}

func WaitEnter() {
	fmt.Println("エンターを押すと処理を終了します。")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}

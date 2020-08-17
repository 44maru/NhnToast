package pkg/log

import (
	"fmt"
	"io"
	"log"
	"os"
)

const LogFile = "./info.log"

func Init() {
	logfile, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		failOnError(fmt.Sprintf("%sのオープンに失敗しました", LogFile), err)
	}
	defer logfile.Close()
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)
}

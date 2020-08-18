package config

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	UserInfo UserInfo
	Thread   Thread
	Instance Instance
}

type UserInfo struct {
	TenantId    string `toml:"tenantId"`
	UserName    string `toml:"user"`
	ApiPassword string `toml:"apiPassword"`
}

type Thread struct {
	ThreadNum                     int           `toml:"threadNum"`
	SleepSecBeforeJointFloatingIp time.Duration `toml:"sleepSecondsBeforeJointFloatingIp"`
}

type Instance struct {
	ImageName string `toml:"imageName"`
}

const CONFIG_FILE_PATH = "./config.toml"

func LoadConfig() (*Config, error) {
	config := new(Config)
	_, err := toml.DecodeFile(CONFIG_FILE_PATH, config)
	if err != nil {
		log.Println("config parse error.")
		return nil, err
	}

	return config, nil
}

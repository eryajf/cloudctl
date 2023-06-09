package public

import (
	"io/ioutil"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Configuration 项目配置
type Configuration struct {
	TcSecretID  string `yaml:"tc_secretId"`
	TcSecretKey string `yaml:"tc_secretKey"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		data, err := ioutil.ReadFile("config.yml")
		// 此处如果使用go-bindata将本地配置文件转成二进制，那么使用如下代码读取配置文件
		// data, err := Asset("config.yml")
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			log.Fatal(err)
		}

		// 如果环境变量有配置，读取环境变量
		tc_secret_id := os.Getenv("TC_SERCRET_ID")
		if tc_secret_id != "" {
			config.TcSecretID = tc_secret_id
		}
		tc_secret_key := os.Getenv("TC_SERCRET_KEY")
		if tc_secret_key != "" {
			config.TcSecretKey = tc_secret_key
		}
	})

	return config
}

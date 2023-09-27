package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var ConfigObj *Config

type Config struct {
	viper *viper.Viper
}

//type GrpcConfig struct {
//	Name string
//	Addr string
//}

func InitConfig(params ...string) *Config {
	//params[0]:  path    配置文件的路径
	log.Printf("run config")
	ConfigObj := &Config{viper: viper.New()}
	workDir := ""
	if params != nil && params[0] != "" {
		workDir = params[0]
	} else {
		workDir, _ = os.Getwd()
	}
	if params != nil && params[1] != "" {
		ConfigObj.viper.SetConfigName(params[1]) // name of config file (without extension)
	} else {
		ConfigObj.viper.SetConfigName("config") // name of config file (without extension)
	}
	//conf.viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	ConfigObj.viper.SetConfigType("json") // REQUIRED if the config file does not have the extension in the name

	//viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	//viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	//conf.viper.AddConfigPath(workDir + "/config") // optionally look for config in the working directory
	ConfigObj.viper.AddConfigPath(workDir) // optionally look for config in the working directory

	err := ConfigObj.viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Fatalln(err)
	}
	return ConfigObj
}

package main

import (
	"fmt"
	"github.com/spf13/viper"
	"httpsproxy/httpsserve"
	"log"
	"net"
	"os"
)

var logger = log.New(os.Stderr, "httpsproxy:", log.Llongfile|log.LstdFlags)

func init() {
	initConfig()
}

func main() {
	listenAddress := viper.GetString("app.listenAddress")
	fmt.Println("开始监听端口:", listenAddress)
	if !checkAddress(listenAddress) {
		logger.Fatal("监听端口失败")
	}

	ipList := viper.GetStringSlice("app.ip_list")
	fmt.Println("代理ip数:", len(ipList))
	if len(ipList) < 1 {
		logger.Fatal("代理ip不能为空")
	}

	if !checkAddress(listenAddress) {
		logger.Fatal("监听端口失败")
	}

	httpsserve.Serve(listenAddress)

}

func checkAddress(listenAddress string) bool {
	_, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		return false
	}
	return true

}

func initConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath(".") // 添加搜索路径

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

package main

import (
	"flag"
	"httpsproxy/httpsserve"
	"log"
	"net"
	"os"
)

var logger = log.New(os.Stderr, "httpsproxy:", log.Llongfile|log.LstdFlags)

func main(){
	var listenAdress string
	flag.StringVar(&listenAdress, "L", "0.0.0.0:8080", "listen address.eg: 127.0.0.1:8080")
	flag.Parse()

	if !checkAdress(listenAdress){
		logger.Fatal("-L listen address format incorrect.Please check it")
	}

	httpsserve.Serve(listenAdress)

}

func checkAdress(adress string) bool{
	_, err := net.ResolveTCPAddr("tcp", adress)
	if err != nil{
		return false
	}
	return true

}

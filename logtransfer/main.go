package main

import (
	"LogCollectionProject/logtransfer/config"
	"LogCollectionProject/logtransfer/es"
	"LogCollectionProject/logtransfer/kafka"
	"fmt"
	"github.com/sirupsen/logrus"
)

//logtransfer从kafka中消费数据，并传入es用于页面展示与搜索
func main() {
	//0.ini配置解析
	err,configObj:= config.LoadConfig()
	if err !=nil {
		logrus.Error("config.loadConfig,err:%v", err)
	}
	fmt.Printf("parse config success, configOBJ:%v \n",configObj)

	//1. 连接kafka
	err = kafka.InitKafka([]string{configObj.KafkaConfig.Address},configObj.KafkaConfig.Topic)
	if err !=nil {
		logrus.Error("kafka.InitKafka,err:%v", err)
	}
	fmt.Println("inint kafka consumer success")

	//2. 连接ES
	err = es.InitES(configObj.ESConf.Address,
		configObj.ESConf.Index,
		configObj.ESConf.ChanSize,
		configObj.ESConf.GorutineNum)
	if err !=nil {
		logrus.Error("es.InitES,err:%v", err)
	}
	fmt.Println("inint ES client success")

	select  {
		
	}
}




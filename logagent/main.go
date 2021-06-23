package main

import (
	"LogCollectionProject/logagent/etcd"
	"LogCollectionProject/logagent/kafka"
	"LogCollectionProject/logagent/tailfile"
	"fmt"
	"github.com/sirupsen/logrus" //日志打印包
	"gopkg.in/ini.v1"            //ini包用于映射配置文件成go中的结构体
)


//日志收集的客户端
//target:收集指定目录下的日志文件，发送到kafka中
func main() {
	//0.ini配置文件解析(初始化连接kafka、读取文件准备)
	var configObj = new(Config)
	cfg , err := ini.Load("./config/config.ini")
	if err != nil {
		logrus.Error("load config failed,err:%v", err)
		return
	}

	//0.1 将配置文件设置成结构体
	err = cfg.MapTo(configObj)
	if err != nil {
		logrus.Error("cfg.MapTo failed,err:%v", err)
		return
	}
	fmt.Printf("configObj:%s\n",configObj)

	//1. 初始化通过sarama连接kafka
	err = kafka.InitKafka([]string{configObj.KafkaConfig.Address},configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Error("kafka:InitKafka failed,err:%v", err)
		return
	}
	fmt.Println("inint kafka client success")

	//2.根据配置初始化etcd
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Error("etcd:InitEtcd failed,err:%v", err)
		return
	}
	fmt.Println("inint etcd client success")

	//2.1从etcd中拉去需要收集的日志的配置项，方便tial去根据配置项读取日志内容
	allConfig:=etcd.GetCollectionConfig(configObj.EtcdConfig.CollectKey)
	fmt.Printf("allConfig:%s\n",allConfig)
	//3.根据配置中的日志路径初始化tial
	err = tailfile.Init(configObj.CollectionConfig.LogFilepath)
	if err != nil {
		logrus.Error("tailfile:InitTailfile failed,err:%v", err)
		return
	}
	fmt.Println("inint tailfile client success")

	//4.把tial读取的日志内容通过sarama发送到kafka中
	//tailObj --> log --> kafkaClient -->kafka
	err = run()
	if err!=nil{
		logrus.Error("main:run failed,err:%v", err)
		return
	}
	fmt.Println("run logagent success")
}

//真正的业务逻辑
func run()(err error){
	//1.通过tialObj读取日记文件,放入kafkaClient的channel中
	tailfile.CollectLogMsg()

	//2.从kafkaClient的channel中取出消息，发送给kafka
	//在kafka的init中就直接执行下面的函数
	//go kafka.SendMsgToKafka()
	return
}



package config

import (
	modle "LogCollectionProject/logtransfer/model"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

/**
 * @author xhli
 * @date 2021/7/4 16:25
 * @version 1.0
 * @description: TODO
 */

func LoadConfig()(err error, configObj *modle.LogTransferConfig){
	//0.ini配置文件解析
	configObj = new(modle.LogTransferConfig)
	cfg , err := ini.Load("logtransfer/config/logtransfer.ini")
	if err != nil {
		logrus.Error("load config failed,err:%v", err)
		return err,nil
	}

	//0.1 将配置文件设置成结构体
	err = cfg.MapTo(configObj)
	if err != nil {
		logrus.Error("cfg.MapTo failed,err:%v", err)
		return err,nil
	}

	return nil,configObj
}
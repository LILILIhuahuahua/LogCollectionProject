package etcd

import (
	"LogCollectionProject/logagent/common"
	"LogCollectionProject/logagent/tailfile"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//etcd相关操作
var(
	//etcd的客户端对象
	Client *clientv3.Client
)

//初始化etcd连接
func Init(address []string)(err error){
	//1.连接etcd
	Client, err = clientv3.New(clientv3.Config {
		Endpoints: address, // etcd的节点，可以传入多个
		DialTimeout: 5*time.Second, // 连接超时时间
	})

	if err != nil {
		fmt.Printf("connect to etcd failed, err: %v \n", err)
		return err
	}
	fmt.Println("etcd:connect to etcd success")

	return
}


/**
	通过etcd拉取日志收集的配置项的函数
	假定:在ectd中存有json格式的日志收集的配置项（路径）
 **/
func GetCollectionConfig(key string)(collectEntryList []common.CollectEntry){
	//get操作，设置1秒超时
	//使用context是为了设置连接的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//拉取json格式的日志收集的配置项（路径）
	resp, err := Client.Get(ctx, key)
	if err != nil {
		fmt.Printf("get conf from etcd failed, err:%v \n", err)
		return
	}

	if len(resp.Kvs) == 0 {
		logrus.Warnf("get len=0 conf from etcd by key:%s",key)
		return
	}
	//解析json映射成collectEntry对象
	result :=resp.Kvs[0]
	fmt.Printf("result.Value:%s\n",result.Value)
	err = json.Unmarshal(result.Value,&collectEntryList)
	if err != nil {
		fmt.Printf("etcd.GetCollectionConfig:json.Unmarshal failed, err:%v \n", err)
		return
	}
	return collectEntryList
}

//通过watch监听etxd中`key="collect_log_conf"`配置的变化
func WatchConfig(collectKey string){
	watchChan := Client.Watch(context.Background(),collectKey)

	// 从通道中尝试获取值（监视的信息）
	var newCollectEntryList []common.CollectEntry
	for wresp := range watchChan {
		for _, watchMsg := range wresp.Events{
			fmt.Printf("etcd.WatchConfig: etch config chang-->Type:%v key:%s value:%s \n", watchMsg.Type, watchMsg.Kv.Key, watchMsg.Kv.Value)
			//解析新的配置信息
			err := json.Unmarshal(watchMsg.Kv.Value,&newCollectEntryList)
			if err != nil {
				fmt.Printf("etcd.WatchConfig:json.Unmarshal failed, err:%v \n", err)
				continue
			}
			//告诉tailfile模块，需要收集的日志Config改变了
			tailfile.SendNewConf(newCollectEntryList)
			fmt.Printf("etcd.WatchConfig: newCollectEntryList:%s \n",newCollectEntryList)


		}
	}
}




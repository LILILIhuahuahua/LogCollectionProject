package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//watch:监听etcd中某个key的变化
func main()  {
	cli, err := clientv3.New(clientv3.Config {
		Endpoints: []string{"127.0.0.1:2379"}, // etcd的节点，可以传入多个
		DialTimeout: 5*time.Second, // 连接超时时间
	})

	if err != nil {
		fmt.Printf("connect to etcd failed, err: %v \n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()

	// watch
	// 派一个哨兵，一直监视着 s4 这个key的变化（新增，修改，删除），返回一个只读的chan
	watchChan := cli.Watch(context.Background(), "s4")

	// 从通道中尝试获取值（监视的信息）
	for wresp := range watchChan {
		for _, watchMsg := range wresp.Events{
			fmt.Printf("Type:%v key:%s value:%s \n", watchMsg.Type, watchMsg.Kv.Key, watchMsg.Kv.Value)
		}
	}
}

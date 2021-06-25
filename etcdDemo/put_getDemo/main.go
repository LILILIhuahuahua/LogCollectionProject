package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//
func main() {
	//1.连接etcd
	cli, err := clientv3.New(clientv3.Config {
		Endpoints: []string{"127.0.0.1:2379"}, // etcd的节点，可以传入多个
		DialTimeout: 5*time.Second, // 连接超时时间
	})

	if err != nil {
		fmt.Printf("connect to etcd failed, err: %v \n", err)
		return
	}
	fmt.Println("connect to etcd success")

	// 延迟关闭
	defer cli.Close()

	// put操作  设置1秒超时
	//使用context是为了设置连接的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "collect_log_conf", "[{\"path\":\"E:/Kafka/kafka_2.8/log/s5.log\",\"topic\":\"s5_log\"},{\"path\":\"E:/Kafka/kafka_2.8/log/s6.log\",\"topic\":\"s6_log\"},{\"path\":\"E:/Kafka/kafka_2.8/log/s7.log\",\"topic\":\"s7_log\"}]\n")
	//cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v \n", err)
		return
	}

	// get操作，设置1秒超时
	//使用context是为了设置连接的超时时间
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "collect_log_conf")
	defer cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v \n", err)
		return
	}
	for _, msg:=range resp.Kvs{
		fmt.Printf("key:%s,value:%s",msg.Key,msg.Value)
	}



}

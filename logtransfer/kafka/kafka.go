package kafka

import (
	"LogCollectionProject/logtransfer/es"
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
	"time"
)

/**
 * @author xhli
 * @date 2021/7/4 16:22
 * @version 1.0
 * @description: TODO
 */

//初始化Kafka连接
func InitKafka(address []string,topic string) (err error){
	// 创建新的消费者
	consumer, err:= sarama.NewConsumer(address, nil)
	if err != nil {
		fmt.Printf("kafka.InitKafka:fail to start consumer, err:%v\n", err)
		return err
	}
	// 拿到指定topic下面的所有分区列表
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return err
	}
	fmt.Printf("getPatition,topic is :%s, partitionList is %v \n",topic,partitionList)


	//消费数据
	go consumerData(consumer,partitionList,topic)
	return nil
}


func consumerData(consumer sarama.Consumer,partitionList []int32,topic string)(err error){
	//等待一秒，确保ES的chan有被初始化
	time.Sleep(time.Second)
	fmt.Println("kafka start consumerData")
	//消费数据
	var wg sync.WaitGroup //
	for partition := range partitionList{ // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition),sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n",
				partition, err)
			return err
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		wg.Add(1)
		go func(sarama.PartitionConsumer){
			for msg:=range pc.Messages(){
				//将数据发到ES的chan中
				fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s",
					msg.Partition, msg.Offset, msg.Topic, msg.Value)
				es.PutLogData(string(msg.Value))
			}
		}(pc)
	}
	wg.Wait()
	return nil
}
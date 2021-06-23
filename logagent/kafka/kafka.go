package kafka

//kafka相关操作
import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var(
	//kafka客户端对象
	client sarama.SyncProducer
	//存储需要发送到kafka的msg指针的管道
	//chan中存放指针 （直接存数据非常影响性能）
	MsgChan chan *sarama.ProducerMessage
)

//将文本类型转换成kafkaMsg类型
func ToKafkaMsg(topic,msgText string)(kafkaMsg *sarama.ProducerMessage){
	msg :=&sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(msgText)
	return msg
}

//初始化一个全局的kafka的客户端
func InitKafka(address []string,chanSize int) (err error){
	//1.生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	//2.连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		return err
	}

	//3.初始化其他数据结构
	MsgChan = make(chan *sarama.ProducerMessage,chanSize)

	//起一个goroutine专门用于发送msg到kafka
	go SendMsgToKafka()
	return nil
}


//从MsgChan中读取msg，发送给kafka
func SendMsgToKafka(){
	for {
		//从MsgChan中读取msg
		select {
		case msg := <-MsgChan:
			//fmt.Println(msg)
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Error("kafka:SendMsgToKafka failed, err:", err)
				return
			}
			logrus.Info("send msg to kafka success,pid:%v,offset:%v",pid,offset)

		}
	}
}


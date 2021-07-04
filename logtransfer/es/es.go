package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

/**
 * @author xhli
 * @date 2021/7/4 16:49
 * @version 1.0
 * @description: TODO
 */

//ES.go将日志数据写入Es，传入kibana做图像展示

//ESClient结构体
type ESClient struct {
	client *elastic.Client
	LogDataChan chan string //Kafka通过这个chan将数据发给ES
	Index string //这个ESClient对应的ES数据库
}

//
type ESMsg struct {
	msg string
}

var (
	esClient *ESClient
)
//初始化es的连接
func InitES(address string,index string,chanSize,gorutineNum int)(err error){
	//创建一个客户端
	url :="http://"+address
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		fmt.Printf("ES.InitES :fail to start consumer, url:%v  err:%v\n", url,err)
		return err
	}

	//构建esClient结构体对象
	esClient = &ESClient{client: client,Index: index}
	esClient.LogDataChan = make(chan string,chanSize)
	fmt.Println("connect to es success")

	//起多个gorutine消息ESChan中输入，写入ES中
	for i:=0;i<gorutineNum;i++{
		go sendToES(index)
	}

	return nil
}

// PutLogData 
/**
 * @Author xhli
 * @Description esClient对外暴露的方法,用于将数据传入ES的chan中
 * @Date 17:41 2021/7/4
 * @Param  msg
 * @return nil
 **/
func PutLogData(msg string){
	esClient.LogDataChan <- msg
}

func GetLogData()string{
	msg :=<- esClient.LogDataChan
	return msg
}

//从ES的chan中取数据，插入到ES中指定的index(数据库)中
func sendToES(index string){
	for msg:= range esClient.LogDataChan{
		esmsg :=ESMsg{msg: msg}
		//插入数据
		put1, err := esClient.client.Index().
			Index(index).
			BodyJson(esmsg).  //Json化数据  (这里只能传结构体)
			Do(context.Background()) //执行插入操作
		if err != nil {
			fmt.Printf("ES.sendToES.put :fail err:%v\n",err)
			return
		}
		fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}


}
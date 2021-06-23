package tailfile

import (
	"LogCollectionProject/logagent/kafka"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	//tial对象
	TailObj *tail.Tail
)
func Init(fileName string)(err error){
	//文件读取的权限
	config := tail.Config{
		ReOpen: true,
		Follow: true,
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll: true,
	}

	//打开文件
	TailObj, err = tail.TailFile(fileName, config)
	if err != nil {
		return err
	}

	return
}

//通过tialObj读取日记文件，封装成kafka中msg类型，丢到kafka的channel中
func CollectLogMsg(){
	var (
		lineMsg *tail.Line
		ok bool
	)
	//读取文件的每一行,封装成kafka中msg类型，丢到kafka的channel中
	for {
		lineMsg, ok = <-TailObj.Lines
		if !ok {
			logrus.Warn("tail file fial open file, filename:%s\n",
				TailObj.Filename)
			//读取出错，等一秒在读
			time.Sleep(time.Second)
			continue
		}
		if len(strings.Trim(lineMsg.Text,"\r"))==0{
			continue
		}
		//2.读取到的一行msg封装成kafka中msg类型，丢到kafka的channel中
		fmt.Printf("len = %v  ,readLog-lineMsg:%v\n",len(lineMsg.Text),lineMsg.Text)
		kafkaMsg := kafka.ToKafkaMsg("shopping", lineMsg.Text)
		kafka.MsgChan <- kafkaMsg
	}
}




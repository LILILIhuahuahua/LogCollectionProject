package tailfile

import (
	"LogCollectionProject/logagent/common"
	"LogCollectionProject/logagent/kafka"
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

//tailTask结构体: 每个tailTask对应一个tailObj，通过tailObj读取对应path下文件的内容
type tailTask struct {
	path string
	topic string
	tailObj *tail.Tail
	ctx context.Context
	cancel context.CancelFunc
}

var(
	Config  tail.Config
)


func Init(allConfig []common.CollectEntry)(err error){
	//allConfig中存有若干个日志的手机项
	//针对每一个日志收集项创建一个对象的tailObj

	//1.初始化一个tailTask的管理者
	tailTaskManager := InitTailTaskManager(allConfig)

	//2.设置tail读取文件的权限
	Config = tail.Config{
		ReOpen: true,
		Follow: true,
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll: true,

	}

	//2.1针对每一个日志收集项创建一个对象的tailObj
	for _,colEntry:= range allConfig{
		//根据每个收集日志项，创建并启动一个tail
		Run(tailTaskManager,colEntry)
	}


	//3.TTaskManager不断监听新的配置 (从SendNewConf不断取新的config)
	go TTaskManager.watch()

	return
}

//将新的CollectEntry传入collectEntryChan中
//便于tail根据新的CollectEntry创建tailObj读取日志文件
func SendNewConf(newConf []common.CollectEntry){
	TTaskManager.collectEntryChan <-newConf
}

//通过tialObj读取日记文件，封装成kafka中msg类型，丢到kafka的channel中
func (ttask *tailTask)CollectLogMsg(){
	var (
		lineMsg *tail.Line
		ok bool
	)
	//读取文件的每一行,封装成kafka中msg类型，丢到kafka的channel中

	for {
		//通过context优雅的关闭gorutine
		select {
		case <-ttask.ctx.Done():
			fmt.Printf("关闭一个ttask.CollectLogMsg的gorutine , path：%s\n",ttask.path)
			return
		case lineMsg, ok = <-ttask.tailObj.Lines:
			if !ok {
				logrus.Warnf("tail file fial open file, filename:%s\n",
					ttask.tailObj.Filename)
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
}

//创建并运行一个tail对象
func Run(tailTaskManager tailTaskManager,colEntry common.CollectEntry)(err error){
	//打开文件
	TailObj, err := tail.TailFile(colEntry.Path, Config)
	if err != nil {
		logrus.Errorf("tailfile.Init: create tailObj failed,path:%s\n",colEntry.Path)
		return err
	}
	ctx,cancel :=context.WithCancel(context.Background())
	ttask :=tailTask{
		path:colEntry.Path,
		topic: colEntry.Topic,
		tailObj: TailObj,
		ctx: ctx,
		cancel: cancel,
	}
	//创建一个ttask就交给tailTaskManager管理,方便后续管理
	tailTaskManager.tailTaskMap[ttask.path] = &ttask

	//创建任务成功，直接让tailObj去收集日记
	fmt.Printf("creat a tailObj,  path:%s\n",ttask.path)
	go ttask.CollectLogMsg()
	return
}





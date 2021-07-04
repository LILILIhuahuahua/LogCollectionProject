package modle

//整个项目的配置结构体
type LogTransferConfig struct {
	KafkaConfig      `ini:"kafka"`
	ESConf           `ini:"es"`
}

type KafkaConfig struct {
	Address string `ini:"address"`
	Topic string `ini:"topic"`
}

type ESConf struct {
	Address string `ini:"address"`
	Index string `ini:"index"`  //ES中index代表数据库
	ChanSize int `ini:"chan_size"` //ESClient中chan的大小
	GorutineNum int  `ini:"gorutine_num"` //ESClient最多起多少gorutine来消费ESchan中数据
}
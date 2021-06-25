package common

//整个项目的配置结构体
type Config struct {
	KafkaConfig      `ini:"kafka"`
	CollectionConfig `ini:"collect"`
	EtcdConfig       `ini:"etcd"`
}

type KafkaConfig struct {
	Address string `ini:"address"`
	Topic string `ini:"topic"`
	ChanSize int `ini:"chan_size"`
}

type CollectionConfig struct {
	LogFilepath string `ini:"logfile_path"`
}

type EtcdConfig struct {
	Address string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

//日志收集项的数据结构
type CollectEntry struct {
	Path string `json:"path"`//日志所在路径
	Topic string `json:"topic"`//日志的主题
}


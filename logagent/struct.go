package main

//整个项目的配置结构体
type Config struct {
	KafkaConfig  `ini:"kafka"`
	CollectionConfig`ini:"collect"`
	EtcdConfig `ini:"etcd"`
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
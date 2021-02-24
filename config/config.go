package config

type AppConf struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
}
type KafkaConf struct {
	Address     string `ini:"address"`
	ChanMaxSize int    `ini:"chan_max_size"`
}
type TaillogConf struct {
	FileName string `ini:"filename"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
	Key     string `ini:"collect_log_key"`
}

package main

import (
	"fmt"
	"golearn/golearn1/logagent/config"
	"golearn/golearn1/logagent/etcd"
	"golearn/golearn1/logagent/kafka"
	"golearn/golearn1/logagent/taillog"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

var (
	cfg = new(config.AppConf)
)

func run() {
	//读取日志 发生到kafka

}

func main() {
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		fmt.Println("map to ini failed,err:", err)

	}

	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.ChanMaxSize)
	if err != nil {
		fmt.Printf("init kafka failed,err:%v\n", err)
		return
	}
	fmt.Println("init kafka ok")
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Println("init etcd failed,err:=", err)
		return
	}

	fmt.Println("init etcd ok")
	//1.从ETCD中获取日志收集项的配置信息
	//ipStr,err:=utils.GetOutboundIP()
	//if err != nil {
	//	panic(err)
	//
	//}
	//etcdConfKey:=fmt.Sprintf(cfg.EtcdConf.Key,ipStr)
	logEntryConf, err := etcd.GetConf(cfg.EtcdConf.Key)
	fmt.Println(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Println("get conf from etcd failed,err:=", err)
		return
	}
	fmt.Printf("get conf from etcd suceess value=%v\n", logEntryConf)
	//2.派哨兵监控日志收集项配置信息
	taillog.Init(logEntryConf)
	newConfChan := taillog.NewConfChan() //从taillog包中获取暴露通道
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(cfg.EtcdConf.Key, newConfChan) //哨兵发现新配置会通知上面那个NewConfChan
	wg.Wait()
	//err=taillog.Init(cfg.TaillogConf.FileName)
	//if err!= nil {
	//	fmt.Printf("init taillog failed,err:%v\n", err)
	//	return
	//}
	//fmt.Println("init taillog ok")

	//run()
}

package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func Init(addr string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Println("connected to etcd fail,err:=", err)
		return
	}

	return

}

// 从ETCD中根据KEY获取配置项
func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:=%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		err = json.Unmarshal(ev.Value, &logEntryConf)
		if err != nil {
			fmt.Printf("unmarshal failed, err:=%v", err)
			return
		}
	}
	return
}

//etcd watch
func WatchConf(key string, newConfCh chan<- []*LogEntry) {
	ch := cli.Watch(context.Background(), key)
	for wresp := range ch {
		for _, evt := range wresp.Events {
			//先判断操作类型
			var newConf []*LogEntry
			if evt.Type != clientv3.EventTypeDelete {
				//如果是删除操作

				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("unmarshal failed,err+%v\n ", err)
					continue
				}
			}
			fmt.Printf("get new conf:%v\n", newConf)
			newConfCh <- newConf
		}

	}
}

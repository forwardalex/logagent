package taillog

import (
	"fmt"
	"golearn/golearn1/logagent/etcd"
	"time"
)

type taillogMgr struct {
	logEntry    []*etcd.LogEntry
	tskMap      map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

var tskMgr *taillogMgr

func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &taillogMgr{
		logEntry:    logEntryConf,
		tskMap:      make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry),
	}
	for _, logEntry := range logEntryConf {

		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		mk := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		tskMgr.tskMap[mk] = tailObj
	}
	go tskMgr.run()
}

//监听NewConfChan 有新配置做相应处理
func (t *taillogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.tskMap[mk]
				if ok {
					continue
				} else {
					tailObj := NewTailTask(conf.Path, conf.Topic)
					t.tskMap[mk] = tailObj
				}
			}
			for _, c1 := range t.logEntry {
				isDelete := true
				for _, c2 := range newConf {
					if c2.Path == c1.Path && c2.Topic == c1.Topic {
						continue
					}
				}
				if isDelete {
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.tskMap[mk].cancelFunc()
				}
			}

			fmt.Println("new conf come!", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

//向外暴露一个函数 提供tskMgr的NewConfChan
func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}

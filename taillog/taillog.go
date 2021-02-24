package taillog

import (
	"context"
	"fmt"
	"golearn/golearn1/logagent/kafka"

	"time"

	"github.com/hpcloud/tail"
)

type TailTask struct {
	path       string
	topic      string
	instance   *tail.Tail
	ctx        context.Context //为了实现推出t.run
	cancelFunc context.CancelFunc
}

//Init 初始化taillog
func NewTailTask(path, topic string) (taiObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	taiObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	taiObj.init()
	return
}
func (t *TailTask) init() {

	config := tail.Config{

		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	go t.run()

}

//ReadChan ...
func (t *TailTask) run() {
	for {
		select {
		case <-t.ctx.Done():
			fmt.Printf("tail task:%s_%s exit\n", t.path, t.topic)
			return
		case lines := <-t.instance.Lines:
			//fmt.Println("tttttt")
			//fmt.Println(string(lines.Text))
			//kafka.SendtoKafka(t.topic,lines.Text)
			//把日志数据发送到一个通道中
			kafka.SendtoChan(t.topic, lines.Text)
		default:
			time.Sleep(time.Second)
			//fmt.Println("xxxxxxxx")
		}

	}

}

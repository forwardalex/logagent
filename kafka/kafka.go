package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type logData struct {
	topic string
	data  string
}

var (
	producer    sarama.SyncProducer
	logDataChan chan *logData
)

//Init 初始化client
func Init(addrs []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	producer, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	logDataChan = make(chan *logData, maxSize)
	go sendtoKafka()
	return
}

//SendtoKafka 发送信息到kafka中
func sendtoKafka() {
	for {
		select {
		case ld := <-logDataChan:
			fmt.Printf("see topic:= %v, see data:= %v\n", ld.topic, ld.data)
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			//defer client.Close()
			//send
			pid, offset, err := producer.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed err:", err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
			fmt.Println("send msg sucess")
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

}
func SendtoChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/kppotato/kafka_monitor/g"
	"strings"
	"fmt"
	"github.com/kppotato/kafka_monitor/model"
)

func GetTopicLogSize() map[string]*model.KafkaTopic  {
	cfg := sarama.NewConfig()
	client, err := sarama.NewClient(strings.Split(g.Opts.KafkaAddress,","), cfg)
	if err != nil{
		g.Logger.Error("conn kafka faild",g.Opts.KafkaAddress,err)
	}
	defer client.Close()

	myKafkaTopic :=make(map[string]*model.KafkaTopic,0)
	topics,_ :=client.Topics()
	for _, t := range topics{
		partitions,_ :=client.Partitions(t)
		ktopic:= &model.KafkaTopic{}
		ktopic.Topic=t
		for _,i :=range partitions{
			offset, err :=client.GetOffset(t,i,sarama.OffsetNewest)
			if err != nil{
				fmt.Println(err)
			}
			ktopic.Partitions=append(ktopic.Partitions,&model.ZKoffset{ Partitions:int(i),Offset: offset})
			//fmt.Printf("topic:%s  partition:%d  offset:%d\n",t,i,offset)
			ktopic.TotalOffset += offset
		}
		myKafkaTopic[t]=ktopic
		//fmt.Printf("topic:%s total offset:%d\n",t,offset_total)
	}
	return myKafkaTopic
}
func GetTipic()[]string{
	cfg := sarama.NewConfig()
	client, err := sarama.NewClient(strings.Split(g.Opts.KafkaAddress,","), cfg)
	if err != nil{
		g.Logger.Error("conn kafka faild",g.Opts.KafkaAddress,err)
	}
	defer client.Close()

	topics,_ :=client.Topics()
	return topics
}
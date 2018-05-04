package cron

import (
	"time"
	"github.com/kppotato/kafka_monitor/zookeeper"
	"github.com/kppotato/kafka_monitor/g"
	"github.com/prometheus/client_golang/prometheus"
)

func AddPrometheus()  {
	for{
		<- time.After(1*time.Minute)
		list := zookeeper.GetConsumerGroupZKForPrometheus()
		for _, x :=range list{
			g.TopicLag.With(prometheus.Labels{"cluster":g.Opts.KafkaName,"topic":x.Topicname,"group":x.Groupname}).Set(float64(x.Lag))
			g.TopicLogSize.With(prometheus.Labels{"cluster":g.Opts.KafkaName,"topic":x.Topicname,"group":x.Groupname}).Set(float64(x.LogSize))
			g.TopicOffset.With(prometheus.Labels{"cluster":g.Opts.KafkaName,"topic":x.Topicname,"group":x.Groupname}).Set(float64(x.Offset))
			g.Logger.Infof("topic:%s group:%s logsize:%d offset:%d lag:%d",x.Topicname,x.Groupname,x.LogSize,x.Offset,x.Lag)
		}
	}
}

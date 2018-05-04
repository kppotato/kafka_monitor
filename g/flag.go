package g

import (
	"flag"
)
func SetFlag(opt *Options) *flag.FlagSet{
	flagset := flag.NewFlagSet("Kafka",flag.ExitOnError)
	flagset.String("cfg","","configuration file")
	flagset.String("zookeeper-address",opt.ZookeeperAddress,"zookeeper-address")
	flagset.String("kafka-address",opt.KafkaAddress,"kafka-address")
	flagset.String("zkpath",opt.ZkPath,"zkpath")
	flagset.String("kafka-name",opt.KafkaName,"KafkaName")
	flagset.String("grafana",opt.Grafana,"Grafana")
	flagset.Int("http-port",opt.HttpPort,"http-port")
	flagset.Int("prometheus-port",opt.PrometheusPort,"prometheus-port")
	return flagset
}

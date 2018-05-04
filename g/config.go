package g


var(
	Opts *Options
)

type Options struct{
	Config    						string    		`flag:"cfg"`
	ZookeeperAddress				string		`flag:"zookeeper-address" cfg:"zookeeper_address"`
	KafkaName						string		`flag:"kafka-name" cfg:"kafka_name"`
	KafkaAddress					string		`flag:"kafka-address" cfg:"kafka_address"`
	ZkPath							string		`flag:"zkpath" cfg:"zkpath"`
	Grafana							string		`flag:"grafana" cfg:"grafana"`
	PrometheusPort					int			`flag:"prometheus-port" cfg:"prometheus_port"`
	HttpPort						int			`flag:"http-port" cfg:"http_port"`
}
type Config map[string]interface{}

func NewOption() *Options{
	return &Options{
		Config:"",
		KafkaName:"KP-TEST-CLUSTER",
		PrometheusPort:4567,
		HttpPort:9876,
		Grafana:"http://192.168.51.110:3000/d/JSzmhLMiz/kafka?orgId=1",
		ZookeeperAddress:"cd-kafka00.host-mtime.com:2181,cd-kafka01.host-mtime.com:2181,cd-kafka02.host-mtime.com:2181,cd-kafka03.host-mtime.com:2181,cd-kafka04.host-mtime.com:2181",
	}
}
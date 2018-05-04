package model


type ZkKafka struct{
	Port int
	Host string
}
type ZKTopic struct{
	Name				string
	Partition			string
	Offset				int64
	LogSize				int64
	Lag					int64
	Owner				string
	Created				string
	LastSeen			string
}
type ZKoffset struct{
	Partitions   	int
	Offset       	int64
	Owner			string
}

type ZKConusmer struct{
	ConsumerGroup		string
	Topicname			[]*ZKTopic
}
//{"controller_epoch":8,"leader":5,"version":1,"leader_epoch":12,"isr":[5,2]}
type ZKbroker struct{
	Controller_epoch   		int
	Leader       			int
	Version					int
	Leader_epoch			int
	Isr						[]int
}
//{"jmx_port":-1,"timestamp":"1518967706758","endpoints":["PLAINTEXT://cd-kafka00.host-mtime.com:9092"],"host":"cd-kafka00.host-mtime.com","version":3,"port":9092}
type ZKKakfkaBroker struct{
	Id					int					`json:",omitempty"`
	Host   				string
	Port       			int
}

type ZKConsumerGroup struct{
	ConsumerGroupName			string
	Brokers						[]*ZKKakfkaBroker
	Topics						[]*ZKTopic
}
type KafkaTopic struct{
	Partitions					[]*ZKoffset
	TotalOffset					int64
	Topic						string
}

type Topic struct{
	Name			 string
	Group			 []string
}

func NewZKConsumerGroup() *ZKConsumerGroup{
	return &ZKConsumerGroup{ Topics:make([]*ZKTopic,0)}
}
type TopicPrometheus struct{
	Topicname				string
	Offset					int64
	LogSize					int64
	Lag						int64
	Groupname				string
}




//{"jmx_port":-1,"timestamp":"1515958524775","endpoints":["PLAINTEXT://db-kafka00.host-mtime.com:9092"],"host":"db-kafka00.host-mtime.com","version":3,"port":9092}
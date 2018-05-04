package zookeeper

import (
	"github.com/kppotato/kafka_monitor/g"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
	"encoding/json"
	"github.com/kppotato/kafka_monitor/model"
	"fmt"
	fpath "path"
	//"os"
	"strconv"
	"github.com/kppotato/kafka_monitor/kafka"
	"os"
)

func Init(){
	checkZookeeper()
	getKafkaAddrByZk()
	if g.Opts.KafkaAddress == ""{
		g.Logger.Println("init kafka failed")
		os.Exit(1)
	}
}

func checkZookeeper()  {
	if  g.Opts.ZookeeperAddress == ""{
		g.Logger.Panicln("zookeeper address is empty")
	}
	g.Logger.Info(g.Opts.ZookeeperAddress)
	zkConn,_,err :=zk.Connect(strings.Split(g.Opts.ZookeeperAddress,","),30*time.Second)
	if err != nil{
		g.Logger.Panicln("zookeeper conn is failed",err)
	}
	defer zkConn.Close()
}

func getKafkaAddrByZk(){
	zkConn,_,err :=zk.Connect(strings.Split(g.Opts.ZookeeperAddress,","),30*time.Second)
	if err != nil{
		g.Logger.Error("zookeeper conn is failed")
	}
	defer zkConn.Close()
	g.Logger.Info(combineKakfaPath("/brokers/ids"))
	ids,_,err := zkConn.Children(combineKakfaPath("/brokers/ids"));
	//ids,_,err := zkConn.Children("/kafka/brokers/ids");
	if err != nil{
		g.Logger.Error("zookeeper get kafka address failed")
		return
	}
	kakfaaddress := make([]string,len(ids))
	zkkafka := model.ZkKafka{}
	for i,id := range ids{
		g.Logger.Info(combineKakfaPath("/brokers/ids/",id))
		result,_,err := zkConn.Get(combineKakfaPath("/brokers/ids/",id));
		if err != nil{
			g.Logger.Error("get ids error",id)
		}
		g.Logger.Info(string(result))
		err =json.Unmarshal(result,&zkkafka)
		if err != nil{
			g.Logger.Error("json error",string(result))
			continue
		}
		kakfaaddress[i]=fmt.Sprintf("%s:%d",zkkafka.Host,zkkafka.Port)
		//kakfaaddress =append(kakfaaddress,fmt.Sprintf("%s:%d",zkkafka.Host,zkkafka.Port))
	}
	//g.Logger.Info(len(kakfaaddress))
	g.Logger.Info("kafka address :",strings.Join(kakfaaddress,","))
	g.Opts.KafkaAddress = strings.Join(kakfaaddress,",")
}

func combineKakfaPath(path... string) string{
	if g.Opts.ZkPath == ""{
		return fpath.Join(path...)
	}
	result := make([]string,0)
	result=append(result,g.Opts.ZkPath)
	for _,i :=range path{
		result =append(result,i)
	}
	return fpath.Join(result...)
}

func getKafkaBrokers(zkConn *zk.Conn) map[int]*model.ZKKakfkaBroker{
	// kafka brokers的详细信息
	KafkaBrokers := make(map[int]*model.ZKKakfkaBroker,0)
	kafkabrokerslist, _, err := zkConn.Children(combineKakfaPath("/brokers/ids")) //获取  kafaka服务器 具体的 host port
	if err != nil {
		g.Logger.Error("zookeeper get brokers ids failed",err)
		return nil
	}
	for _,ids :=range kafkabrokerslist{
		broker, _, err := zkConn.Get(combineKakfaPath("/brokers/ids",ids))
		if err != nil {
			g.Logger.Error("zookeeper get brokers ids failed",err)
			continue
		}
		v :=&model.ZKKakfkaBroker{}
		json.Unmarshal(broker,v)
		intids ,_ :=strconv.Atoi(ids)
		v.Id=intids
		KafkaBrokers[intids]=v
	}
	return KafkaBrokers
}
func GetConsumerGroupZK() []*model.ZKConsumerGroup{
	zkConn,_,err :=zk.Connect(strings.Split(g.Opts.ZookeeperAddress,","),100*time.Second)
	if err != nil{
		g.Logger.Panicln("zookeeper conn is failed")
		return nil
	}
	KafkaBrokers := getKafkaBrokers(zkConn)    //获取 kafka集群的 brokers 信息
	if KafkaBrokers == nil{
		g.Logger.Panicln("zookeeper get brokers failed")
		return nil
	}
	defer  zkConn.Close()
	consumerGropList :=make([]*model.ZKConsumerGroup,0)
	consumersGroups,_,err := zkConn.Children(combineKakfaPath("/consumers"))
	if err != nil {
		g.Logger.Error("zookeeper get consumers failed",err)
		return nil
	}
	kafkalogSizeTopic := kafka.GetTopicLogSize()
	for _,cgroup :=range consumersGroups{
		result := model.NewZKConsumerGroup()
		consumerGropList=append(consumerGropList, result)
		result.ConsumerGroupName = cgroup
		topics,_,err := zkConn.Children(combineKakfaPath("/consumers/",cgroup,"offsets"))
		if err != nil {
			g.Logger.Error("zookeeper get topics failed",err)
			continue
		}
		for _, topic :=range topics{
			topicTotal := &model.ZKTopic{}
			topicTotal.Partition =""
			result.Topics =append(result.Topics,topicTotal)
			topicTotal.Name=topic
			owners, _, err := zkConn.Children(combineKakfaPath("/consumers/",cgroup,"/owners",topic))
			myOwners :=make([]string,0)
			for _,owner :=range owners{
				ownername, _, err := zkConn.Get(combineKakfaPath("/consumers/",cgroup,"/owners",topic,owner))
				if err != nil {
					g.Logger.Error("zookeeper get topics failed",err)
					continue
				}
				myOwners =append(myOwners,string(ownername))
			}

			partitions,_,err := zkConn.Children(combineKakfaPath("/consumers",cgroup,"offsets",topic))
			if err != nil {
				g.Logger.Error("zookeeper get topics failed",err)
				continue
			}
			topicTotal.LogSize =getTopicTotalLogSize(kafkalogSizeTopic,topic)
			for i , partition :=range partitions{
				ptopic := &model.ZKTopic{}
				b,_,err := zkConn.Get(combineKakfaPath("/consumers",cgroup,"offsets",topic,partition))
				if err != nil {
					g.Logger.Error("zookeeper get topics failed",err)
				}
				size,_:=strconv.ParseInt(string(b),10,64)
				topicTotal.Offset += size
				ptopic.Offset=size
				ptopic.Name=""
				ptopic.Partition = partition
				ptopic.LogSize=getPartitionsLogSize(kafkalogSizeTopic,topic,partition)
				ptopic.Lag=ptopic.LogSize-ptopic.Offset
				if len(myOwners)>i{
					ptopic.Owner = myOwners[i]
				}
				result.Topics =append(result.Topics,ptopic)
			}
			topicTotal.Lag=topicTotal.LogSize-topicTotal.Offset
			brokers, _, err := zkConn.Children(combineKakfaPath("/brokers/topics/",topic,"partitions")) //获取拥有着
			if err != nil {
				g.Logger.Error("zookeeper get brokers failed",err)
				continue
			}
			//{"controller_epoch":8,"leader":5,"version":1,"leader_epoch":12,"isr":[5,2]}
			mybrokers := make([]int,0)
			for _,broker := range brokers{
				strbroker, _, err := zkConn.Get(combineKakfaPath("/brokers/topics/",topic,"partitions",broker,"state"))
				if err != nil {
					g.Logger.Error("zookeeper get brokers failed",err)
				}
				v :=&model.ZKbroker{}
				json.Unmarshal(strbroker,v)
				mybrokers=append(mybrokers, v.Leader)
			}
			removeDup := make(map[int]struct{},0)
			for _, i :=range mybrokers{
				v:= KafkaBrokers[i]
				if _,ok :=removeDup[i];ok{
					continue
				}
				removeDup[i]= struct{}{}
				result.Brokers =append(result.Brokers, v)
			}
		}
	}
	return consumerGropList
}


func GetConsumerGroupZKForPrometheus() []*model.TopicPrometheus{
	zkConn,_,err :=zk.Connect(strings.Split(g.Opts.ZookeeperAddress,","),100*time.Second)
	if err != nil{
		g.Logger.Panicln("zookeeper conn is failed")
		return nil
	}
	defer  zkConn.Close()
	consumerGropList :=make([]*model.TopicPrometheus,0)
	consumersGroups,_,err := zkConn.Children(combineKakfaPath("/consumers"))
	if err != nil {
		g.Logger.Error("zookeeper get consumers failed",err)
		return nil
	}
	kafkalogSizeTopic := kafka.GetTopicLogSize()
	for _,cgroup :=range consumersGroups{
		result := &model.TopicPrometheus{}
		consumerGropList=append(consumerGropList, result)
		result.Groupname = cgroup
		topics,_,err := zkConn.Children(combineKakfaPath("/consumers/",cgroup,"offsets"))
		if err != nil {
			g.Logger.Error("zookeeper get topics failed",err)
			continue
		}
		for _, topic :=range topics{
			result.Topicname=topic
			partitions,_,err := zkConn.Children(combineKakfaPath("/consumers",cgroup,"offsets",topic))
			if err != nil {
				g.Logger.Error("zookeeper get topics failed",err)
				continue
			}
			result.LogSize = getTopicTotalLogSize(kafkalogSizeTopic,topic)
			for _ , partition :=range partitions{
				b,_,err := zkConn.Get(combineKakfaPath("/consumers",cgroup,"offsets",topic,partition))
				if err != nil {
					g.Logger.Error("zookeeper get topics failed",err)
				}
				size,_:=strconv.ParseInt(string(b),10,64)
				result.Offset += size
			}
			result.Lag=result.LogSize-result.Offset
		}
	}
	return consumerGropList
}

func getPartitionsLogSize(m map[string]*model.KafkaTopic,topic string ,partitions string )int64{

	p,_ := strconv.Atoi(partitions)
	Topic, ok := m[topic]
	if ok{
			for _,i :=range Topic.Partitions{
				if i.Partitions == p{
					return i.Offset
				}
			}
		}
	return 0
}

func getTopicTotalLogSize(m map[string]*model.KafkaTopic,topic string) int64{
	Topic, ok := m[topic]
	if ok{
		return Topic.TotalOffset
	}
	return 0
}
package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kppotato/kafka_monitor/g"
	"net/http"
	"fmt"
	"html/template"
	"github.com/kppotato/kafka_monitor/zookeeper"
	"encoding/json"
	"github.com/kppotato/kafka_monitor/model"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"strconv"
)

func HttpStart(){
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/consumerGroup", ConsumerGroup)
	router.GET("/topic", Topic)
	//router.GET("/consumerdetail/:name",Consumerdetail)
	g.Logger.Fatal(http.ListenAndServe("0.0.0.0:"+ strconv.Itoa(g.Opts.HttpPort), router))
}
func HttpPrometheusStart(){
	http.Handle("/metrics", promhttp.Handler())
	g.Logger.Fatal(http.ListenAndServe(":"+strconv.Itoa(g.Opts.PrometheusPort), nil))
}
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//tmpl, err := template.New("").Delims("[[", "]]").ParseFiles("./view/index.html")
	tmpl, err := template.New("index.html").Delims("[[","]]") .ParseFiles("view/index.html")
	//t, err :=template ParseFiles("./view/index.html")
	//t, err := template.ParseFiles("./view/index.html")
	if err != nil{
		g.Logger.Error("template error",err)
	}
	fmt.Println(g.Opts.Grafana)
	err = tmpl.Execute(w,map[string]string{"Grafana":g.Opts.Grafana,"KafkaName":g.Opts.KafkaName})
	if err != nil{
		g.Logger.Error("template error",err)
	}
}

func ConsumerGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))

	//list := make([]*model.ZKConsumerGroup,0)
	//f := &model.ZKConsumerGroup{}
	//list = append(list,f)
	//f.Brokers = make([]*model.ZKKakfkaBroker,0)
	//f.Topics = make([]*model.ZKTopic,0)
	//
	//f.ConsumerGroupName="es_group"
	//f.Brokers =append(f.Brokers,&model.ZKKakfkaBroker{ Id:1,Host:"www.host.com",Port:22 })
	//f.Brokers =append(f.Brokers,&model.ZKKakfkaBroker{ Id:2,Host:"ff.host.com",Port:22 })
	//f.Topics =append(f.Topics,&model.ZKTopic{ Name:"topic_my",Partition:2,Offset:2,LogSize:4,Lag:2,Owner:"" })
	//f.Topics =append(f.Topics,&model.ZKTopic{ Name:"",Partition:2,Offset:2,LogSize:4,Lag:2,Owner:"2" })
	//f.Topics =append(f.Topics,&model.ZKTopic{ Name:"",Partition:2,Offset:2,LogSize:4,Lag:2,Owner:"3" })
	//f.Topics =append(f.Topics,&model.ZKTopic{ Name:"",Partition:2,Offset:2,LogSize:4,Lag:2,Owner:"4" })
	//fmt.Printf("%s",result)
	//b,_:=json.Marshal(list)
	ret,ok :=g.Mycache.Get("kafka")
	if ok {
		b, _ := json.Marshal(ret)
		fmt.Fprintf(w, string(b))
	}else{
		consumer := zookeeper.GetConsumerGroupZK()
		b,_:=json.Marshal(consumer)
		fmt.Fprintf(w,string(b))
		g.Mycache.Set("kafka",consumer,cache.DefaultExpiration)
	}

}

func Topic(w http.ResponseWriter, r *http.Request, ps httprouter.Params){

	//consumer := zookeeper.GetConsumerGroupZK()
	//list := make([]*model.ZKConsumerGroup,0)
	//f := &model.ZKConsumerGroup{}
	//list = append(list,f)
	//f.Brokers = make([]*model.ZKKakfkaBroker,0)
	//f.Topics = make([]*model.ZKTopic,0)
	//
	//f.ConsumerGroupName="es_group"
	//f.Brokers =append(f.Brokers,&model.ZKKakfkaBroker{ Id:1,Host:"www.host.com",Port:22 })
	//f.Brokers =append(f.Brokers,&model.ZKKakfkaBroker{ Id:2,Host:"ff.host.com",Port:22 })
	//f.Topics =append(f.Topics,&model.ZKTopic{ Name:"topic_my",Partition:2,Offset:2,LogSize:4,Lag:2,Owner:"" })
	//f.Topics =append(f.Topics,&model.ZKTopic{ Name:"",Partition:2,Offset:2,LogSize:4,Lag:2,Owner:"2" })
	//f.Topics =append(f.Topics,&model.ZKTopic{ Name:"",Partition:2,Offset:2,LogSize:4,Lag:2,Owner:"3" })
	//f.Topics =append(f.Topics,&model.ZKTopic{ Name:"",Partition:2,Offset:2,LogSize:4,Lag:2,Owner:"4" })

	var consumer []*model.ZKConsumerGroup
	ret,ok :=g.Mycache.Get("kafka")
	if ok {
		consumer,_ = ret.([]*model.ZKConsumerGroup)
	}else{
		consumer = zookeeper.GetConsumerGroupZK()
		g.Mycache.Set("kafka",consumer,cache.DefaultExpiration)
	}
	topicList := make([]*model.Topic,0)
	topicMap :=make(map[string][]string,0)
	for _,x :=range consumer{
		for _,t :=range x.Topics{
			if t.Name == ""{
				continue
			}
			topicMap[t.Name] = append(topicMap[t.Name], x.ConsumerGroupName)
		}
	}

	for k,v :=range topicMap{
		f := &model.Topic{}
		f.Name =k
		f.Group = v
		topicList = append(topicList, f)
	}
	//fmt.Printf("%s",result)
	b,_:=json.Marshal(topicList)
	fmt.Fprintf(w,string(b))
}
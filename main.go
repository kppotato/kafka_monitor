package main

import (
	"log"
	"runtime"
	"syscall"
	"fmt"
	"github.com/judwhite/go-svc/svc"
	"github.com/mreiferson/go-options"
	"github.com/kppotato/kafka_monitor/g"
	"os"
	"github.com/BurntSushi/toml"
	"github.com/kppotato/kafka_monitor/zookeeper"
	"github.com/kppotato/kafka_monitor/http"
	"github.com/kppotato/kafka_monitor/cron"
)
func init()  {
	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())
}
type program struct {

}
func (p *program) Init(env svc.Environment) error {
	g.Opts = g.NewOption()
	flagset := g.SetFlag(g.Opts)
	flagset.Parse(os.Args[1:])

	configFile :=flagset.Lookup("cfg").Value.String()
	var cfg g.Config
	if configFile != "" {
		_,err :=toml.DecodeFile(configFile,&cfg)
		if err != nil {
			log.Printf("ERROR: failed to load config file %s - %s\n",configFile,err.Error())
			os.Exit(1)
		}
	}
	options.Resolve(g.Opts,flagset,cfg)
	zookeeper.Init()
	return nil
}
func (p *program) Start() error {
	cron.Init()
	go http.HttpStart()
	go http.HttpPrometheusStart()

	//fmt.Println(zookeeper.GetConsumerGroupZK())
	//fmt.Println(zookeeper.GetKafkaByZk)
	//fmt.Printf("%+v",zookeeper.GetConsumerGroup())
	//fmt.Printf("%+v\n",kafka.GetTopicLogSize())
	//for _,i :=range ff{
	//	for _,t :=range i.Topicname{
	//		g.Logger.Info(t.Name,t.Offset,i.ConsumerGroup)
	//	}
	//}
	return nil
}
func (p *program) Stop() error {
	return nil
}

func main()  {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
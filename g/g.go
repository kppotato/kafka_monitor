package g

import (
	"github.com/patrickmn/go-cache"
	"time"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Mycache = cache.New(1*time.Minute,1*time.Minute)
	TopicOffset= prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Kafka_topic_Offset",
			Help: "Current topic offset.",
		},
		[]string{"topic","cluster","group"},
		)
	TopicLogSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Kafka_topic_Logsize",
			Help: "Current topic logsize.",
		},
		[]string{"topic","cluster","group"},
	)
	TopicLag = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Kafka_topic_Lag",
			Help: "Current topic Lag.",
		},
		[]string{"topic","cluster","group"},
	)
)

func init(){
	prometheus.MustRegister(TopicOffset)
	prometheus.MustRegister(TopicLogSize)
	prometheus.MustRegister(TopicLag)
}

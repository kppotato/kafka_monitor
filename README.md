# kafka_monitor
本系统主要是收集kafka集群的的消费情况。


使用是技术：
前端使用bootstrap, angular.js, 后端使用:golang   时序数据库：prometheus,图形：Grafana


使用说明书：

      运行golang主程序 主程序需要的参数，
-kafka-name=“你的集群名字”  
-zookeeper-address=“zk地址” 
-zkpath=“zk路径名” 没有可以为空
-grafana="grafana地址"
-prometheus-port=“prometheus获取数据端口”
-http-port=“页面端口地址”

注意：需要把view和static和go程序放在一起使用。
目录结构：main view static

grafana.json是grafana的模板文件，导入后作为-grafana 参数给程序用来展示数据，在grafana里自己配置prometheus地址
![Alt text](/path/img.jpg)







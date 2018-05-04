# kafka_monitor


> 本系统主要是收集kafka集群的的消费情况。


> 技术：
- 前端使用bootstrap, angular.js, 后端使用:golang   时序数据库：prometheus,图形：Grafana


> 使用说明书：

      - 运行golang主程序 主程序需要的参数，
 <table>
 <thead>
 <tr>
  <td>参数名</td>
  <td>描述</td>
  <td>例子</td>
 </tr>
  </thead>
 <tbody>
  <tr>
  <td>kafka-name</td>
  <td>kafka集群名字</td>
  <td>./kafka_monitor -kafka-name=名称</td>
 </tr>
    <tr>
  <td>zookeeper-address</td>
  <td>zk地址 (多个用,间隔)</td>
  <td>./kafka_monitor -zookeeper-address=xxx:000,xxx:000</td>
 </tr>
  <tr>
  <td>zkpath</td>
  <td>zk里路径前缀（可以为空）</td>
  <td>./kafka_monitor -zkpath= </td>
 </tr>
    <tr>
  <td>grafana</td>
  <td>grafana模板地址</td>
  <td>./kafka_monitor -grafana=xxx </td>
 </tr>
      <tr>
  <td>prometheus-port</td>
  <td>prometheus获取数据端口</td>
  <td>./kafka_monitor -prometheus-port=1111 </td>
 </tr>
        <tr>
  <td>http-port</td>
  <td>页面http端口</td>
  <td>./kafka_monitor -http-port=1111 </td>
 </tr>
   </tbody>
 </table>

注意：需要把view和static和go程序放在一起使用。
目录结构：main view static

grafana.json是grafana的模板文件，导入后作为-grafana 参数给程序用来展示数据，在grafana里自己配置prometheus地址
![展示](/path/img.jpg)
![展示2](/path/img2.jpg)







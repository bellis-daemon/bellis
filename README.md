# bellis

分布式云监控告警系统，采用各模块分工明确的设计，模块间通过 GRPC、消息队列（Redis）通信，采样数据通过 InfluxDB 存储。

| 主页面 | 添加监控目标 |
|-|-|
|![image](https://github.com/bellis-daemon/bellis/assets/55825043/d79e7d3c-f01f-46e6-9bbf-22999de58ae6)|![image](https://github.com/bellis-daemon/bellis/assets/55825043/89182639-126c-4762-abb0-cf1d3b22ae11)|



### 架构设计

![整体架构](https://github.com/bellis-daemon/bellis/assets/55825043/4cf373b0-a416-4776-8d6f-61c5b907be99)


### 模板模式

![模板模式](https://github.com/bellis-daemon/bellis/assets/55825043/6fda1272-3d71-455d-b165-1e8b0d4e133c)

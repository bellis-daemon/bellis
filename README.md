# bellis （正在开发）

分布式云监控告警系统，采用各模块分工明确的设计，模块间通过 GRPC、消息队列（Redis）通信，采样数据通过 InfluxDB 存储。

| 主页面                                                                                                    | 添加监控目标                                                                                                 |
|--------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------|
| ![image](https://github.com/bellis-daemon/bellis/assets/55825043/d79e7d3c-f01f-46e6-9bbf-22999de58ae6) | ![image](https://github.com/bellis-daemon/bellis/assets/55825043/89182639-126c-4762-abb0-cf1d3b22ae11) |

### 架构设计

![整体架构](https://github.com/bellis-daemon/bellis/assets/55825043/4cf373b0-a416-4776-8d6f-61c5b907be99)

### 模板模式

```mermaid
classDiagram
    Application --|> Implement: 包含
    Implement --> BT: 实现
    Implement --> Ping: 实现
    Implement --> VPS: 实现
    Implement --> Minecraft: 实现
    class BT{
        Options btOptions
        Client  btgosdk.Client
        Fetch() // 获取监控目标信息
        Init() // 初始化
    }
    class Ping{
        Options pingOptions
        Fetch() // 获取监控目标信息
        Init() // 初始化
    }
    class VPS{
        Options vpsOptions
        client  *http.Client
        Fetch() // 获取监控目标信息
        Init() // 初始化
    }
    class Minecraft{
	    Options minecraftOptions
        Fetch() // 获取监控目标信息
        Init() // 初始化
    }
    class Implement{
        Fetch() // 获取监控目标信息
        Init() // 初始化
    }
    class Application{
        ctx         context.Context // 上下文
        measurement string // 监控目标类型名称
	    deadline    time.Time // 任期到期时间
	    failedCount uint // 失败计数
        handler Implement // 子处理器
        Run() // 开始运行
        refresh() // 刷新监控目标状态
        reclaim() // 将任期重新分配
        alert() // 发出离线警告
        UpdateOptions() // 更新设置
    }
```
# tlog-Go 日志功能及实现 

## 日志配置
```yaml
plugins:
  log:                                      #所有日志配置
    default:                                  #默认日志配置，log.Debug("xxx")
      - writer: console                         #控制台标准输出 默认
        level: debug                            #标准输出日志的级别
      - writer: file                              #本地文件日志
        level: debug                              #本地文件滚动日志的级别
        formatter: json                           #标准输出日志的格式
        formatter_config:
          time_fmt: 2006-01-02 15:04:05           #日志时间格式。"2006-01-02 15:04:05"为常规时间格式，"seconds"为秒级时间戳，"milliseconds"为毫秒时间戳，"nanoseconds"为纳秒时间戳
          time_key: Time                          #日志时间字段名称，不填默认"T"
          level_key: Level                        #日志级别字段名称，不填默认"L"
          name_key: Name                          #日志名称字段名称， 不填默认"N"
          caller_key: Caller                      #日志调用方字段名称， 不填默认"C"
          function_key: Function                  #日志调用方字段名称， 不填默认不打印函数名
          message_key: Message                    #日志消息体字段名称，不填默认"M"
          stacktrace_key: StackTrace              #日志堆栈字段名称， 不填默认"S"
        writer_config:                            #本地文件输出具体配置
          filename: ../log/tlog_size.log          #本地文件滚动日志存放的路径
          write_mode: 3                           #日志写入模式，1-同步，2-异步，3-极速(异步丢弃), 不配置默认极速模式
          roll_type: size                         #文件滚动类型,size为按大小滚动
          max_age: 7                              #最大日志保留天数
          max_backups: 10                         #最大日志文件数
          compress:  false                        #日志文件是否压缩
          max_size: 10                            #本地文件滚动日志的大小 单位 MB
      - writer: file                              #本地文件日志
        level: debug                              #本地文件滚动日志的级别
        formatter: json                           #标准输出日志的格式
        writer_config:                            #本地文件输出具体配置
          filename: ../log/tlog_time.log          #本地文件滚动日志存放的路径
          write_mode: 3                           #日志写入模式，1-同步，2-异步，3-极速(异步丢弃), 不配置默认极速模式
          roll_type: time                         #文件滚动类型,time为按时间滚动
          max_age: 7                              #最大日志保留天数
          max_backups: 10                         #最大日志文件数
          compress:  false                        #日志文件是否压缩
          max_size: 10                            #本地文件滚动日志的大小 单位 MB
          time_unit: day                          #滚动时间间隔，支持：minute/hour/day/month/year
      - writer: atta                                #atta远程日志输出
        remote_config:                              #远程日志配置，业务自定义结构，每一种远程日志都有自己独立的配置
          atta_id: '05e00006180'                    #atta id 每个业务自己申请
          atta_token: '6851146865'                  #atta token 业务自己申请
          message_key: msg                          #日志打印包体的对应atta的field
          field:                                    #申请atta id时，业务自己定义的表结构字段，顺序必须一致
            - msg
            - uid
            - cmd
    custom:                                   #业务自定义的logger配置，名字随便定，每个服务可以有多个logger，可使用 log.Get("custom").Debug("xxx") 打日志
      - writer: file                              #业务自定义的core配置，名字随便定
        caller_skip: 1                            #用于定位日志的调用处
        level: debug                              #业务自定义core输出的级别
        writer_config:                            #本地文件输出具体配置
          filename: ../log/tlog1.log              #本地文件滚动日志存放的路径
      - writer: file                              #本地文件日志
        caller_skip: 1                            #用于定位日志的调用处
        level: debug                              #本地文件滚动日志的级别
        writer_config:                            #本地文件输出具体配置
          filename: ../log/tlog2.log              #本地文件滚动日志存放的路径
```

## 相关概念解析
- logger: 具体打日志的对外接口，每个日志都可以有多个输出，如上配置，log.Debug可以同时输出到console终端和file本地文件，可以任意多个
- log factory: 日志插件工厂，每个日志都是一个插件，一个服务可以有多个日志插件，需要通过日志工厂读取配置信息实例化具体logger并注册到框架中，没有使用日志工厂，默认只输出到终端
- writer factory: 日志输出插件工厂，每个输出流都是一个插件，一个日志可以有多个输出，需要通过输出工厂读取具体配置实例化具体core
- core: zap的具体日志输出实例，有终端，本地日志，远程日志等等
- zap: uber的一个log开源实现，tlog框架直接使用的zap的实现
- with fields: 设置一些业务自定义数据到每条log里:比如uid，imei等，每个请求入口设置一次


## 日志插件实例化流程
- 1. 首先框架会提前注册好 default log factory, console writer factory, file writer factory
- 2. log factory 解析logger配置，遍历writer配置
- 3. 逐个writer调用 writer factory 加载writer配置
- 4. writer factory实例化core
- 5. 多个core组合成一个logger
- 6. 注册logger到框架中

## 支持多logger
1. 首先需要在main函数入口注册插件
```go
	import (
		"github.com/hyperits/tlog/log"
        "github.com/hyperits/tlog/plugin"
	)
	func main() {
		plugin.Register("custom", log.DefaultLogFactory) 
		
		s := tlog.NewServer()
	}
```
2. 配置文件定义自己的logger，如上custom 
3. 业务代码具体场景Get使用不同的logger
```go
	log.Get("custom").Debug("message")
```
4. 由于一个context只能保存一个logger，所以DebugContext等接口只能打印default logger，需要使用XxxContext接口打印自定义logger时，可以在请求入口get logger后设置到ctx里面，如
```go
    tlog.Message(ctx).WithLogger(log.Get("custom"))
    log.DebugContext(ctx, "custom log msg")
```
## 框架日志
1. 框架以尽量不打日志为原则，将错误一直往上抛交给用户自己处理
2. 底层严重问题才会打印trace日志，需要设置环境变量才会开启：export tlog_LOG_TRACE=1

## 关于 `caller_skip` 的说明

首先捋一下使用 `logger` 的几种方式：

1. Default Logger：

```go
log.Debug("default logger") // 使用默认的 logger
```

此时该 log 使用的配置为 `default`:

```yaml
    default:                              #默认日志配置，log.Debug("xxx")
      - writer: console                   #控制台标准输出 默认
        level: debug                      #标准输出日志的级别
      - writer: file                      #本地文件日志
        level: debug                      #本地文件滚动日志的级别
        formatter: json                   #标准输出日志的格式
        writer_config:                    #本地文件输出具体配置
          filename: ../log/tlog_time.log  #本地文件滚动日志存放的路径
          write_mode: 3                   #日志写入模式，1-同步，2-异步，3-极速(异步丢弃), 不配置默认极速模式
          roll_type: time                 #文件滚动类型,time为按时间滚动
          max_age: 7                      #最大日志保留天数
          max_backups: 10                 #最大日志文件数
          compress:  false                #日志文件是否压缩
          max_size: 10                    #本地文件滚动日志的大小 单位 MB
          time_unit: day                  #滚动时间间隔，支持：minute/hour/day/month/year
```

此时不需要关注或者去设置 `caller_skip` 的值，该值默认为 2，意思是在 `zap.Logger.Debug` 上套了两层（`tlog.log.Debug -> tlog.log.zapLog.Debug -> zap.Logger.Debug`）

2. 将自定义的 logger 放到 context 中进行使用：

```go
    tlog.Message(ctx).WithLogger(log.Get("custom"))
    log.DebugContext(ctx, "custom log msg")
```

此时也不需要关注或者去设置 `caller_skip` 的值，该值默认为 2，意思是在 `zap.Logger.Debug` 上套了两层（`tlog.log.DebugContext -> tlog.log.zapLog.Debug -> zap.Logger.Debug`）

配置例子如下：

```yaml
    custom:                           #业务自定义的logger配置，名字随便定，每个服务可以有多个logger，可使用 log.Get("custom").Debug("xxx") 打日志
      - writer: file                  #业务自定义的core配置，名字随便定
        level: debug                  #业务自定义core输出的级别
        writer_config:                #本地文件输出具体配置
          filename: ../log/tlog1.log  #本地文件滚动日志存放的路径
```


3. 不在 context 中使用自定义的 logger：

```go
	log.Get("custom").Debug("message")
```

此时需要将 `custom` logger 的 `caller_skip` 值设置为 1，因为 `log.Get("custom")` 直接返回的是 `tlog.log.zapLog`，调用 `tlog.log.zapLog.Debug` 只在 `zap.Logger.Debug` 上套了一层（`tlog.log.zapLog.Debug -> zap.Logger.Debug`）

配置例子如下：

```yaml
    custom:                           #业务自定义的logger配置，名字随便定，每个服务可以有多个logger，可使用 log.Get("custom").Debug("xxx") 打日志
      - writer: file                  #业务自定义的core配置，名字随便定
        caller_skip: 1                #用于定位日志的调用处
        level: debug                  #业务自定义core输出的级别
        writer_config:                #本地文件输出具体配置
          filename: ../log/tlog1.log  #本地文件滚动日志存放的路径
```

要注意 `caller_skip` 放置的位置（不要放在 `writer_config` 里面），并且对于多个 `writer` 都有 `caller_skip` 时，该 logger 的 `caller_skip` 的值以最后一个为准，比如：

```yaml
    custom:                           #业务自定义的logger配置，名字随便定，每个服务可以有多个logger，可使用 log.Get("custom").Debug("xxx") 打日志
      - writer: file                  #业务自定义的core配置，名字随便定
        caller_skip: 1                #用于定位日志的调用处
        level: debug                  #业务自定义core输出的级别
        writer_config:                #本地文件输出具体配置
          filename: ../log/tlog1.log  #本地文件滚动日志存放的路径
      - writer: file                  #本地文件日志
        caller_skip: 2                #用于定位日志的调用处
        level: debug                  #本地文件滚动日志的级别
        writer_config:                #本地文件输出具体配置
          filename: ../log/tlog2.log  #本地文件滚动日志存放的路径
``` 

最终 `custom` 这个 logger 的 `caller_skip` 值会被设置为 2

__注意：__ 上述用法 2 和用法 3 是冲突的，只能同时用其中的一种


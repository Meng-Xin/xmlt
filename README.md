# XMLT（畜牧论坛） Coding...
> Gin + DDD 模式开发的一个论坛项目（学习使用）
> 畜牧论坛主要用来收集全国畜牧业市场行情，例如 小麦价格等…… 前期主要优化论坛建设、帖子模块、用户模块等。
> 学习交流群：852991268

## 主要技术点
> + DDD 规范设计(我也是初次在Go中使用，可能存在设计不规范、不合理)
> + 面向接口编程，Repository层接口实现、Service层接口实现，并划分Router类集成。
> + JWT身份认证、双Token无感刷新。
> + Cron定时任务使用，定时从Redis获取点赞数据同步到MySQL
> + 使用RabbitMQ消息队列对高并发评论场景的异步解耦。
> + API双限流策略:RedisIP限流 + TokenBucket负载限流；Gorm限流中间件。
> + Logger日志中间件使用、Gorm 慢查询中间件使用、Gorm Scope通用函数使用。
> + 泛型函数使用：自定义Compare泛型的交、并、差 ……Utils工具函数。
> + go-mock 使用对Repository层、Service层接口实现进行打桩测试。
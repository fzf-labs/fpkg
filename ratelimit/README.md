# ratelimit

限流分为单机限流和分布式限流。

# 限流算法

* 令牌桶
* 漏斗
* 滑动窗口
* 自适应限流

## 扩展

https://pkg.go.dev/golang.org/x/time/rate

https://github.com/alibaba/sentinel-golang

https://github.com/uber-go/ratelimit

https://github.com/go-kratos/aegis/tree/main/ratelimit

https://github.com/RussellLuo/slidingwindow

## Reference

- [Sentinel Go版本](https://github.com/alibaba/sentinel-golang)
- [常用微服务框架的适配](https://github.com/sentinel-group/sentinel-go-adapters)
- [动态数据源扩展支持](https://github.com/sentinel-group/sentinel-go-datasources)
- [阿里双11同款流控降级组件 Sentinel Go简介](https://mp.weixin.qq.com/s/j1kTArkROXlymghR1hkDFA)
- [分布式高并发服务三种常用限流方案简介](https://mp.weixin.qq.com/s/zIhQuK1jmHcn5eIqhJfNkw)

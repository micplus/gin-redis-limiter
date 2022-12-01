# Gin Redis Limiter

## Description

Gin middleware with redis to limit client request by IP.

The limiter can limit some path request counts within a sliding window duration, and block the client.

If too many request to the path detected, block the client that it will not have access to the path for a while.

Gin + Redis实现IP黑名单功能。在一段时间内，如果某个IP对注册了中间件的路由访问频率过高，就将这个IP拉黑，限制它的访问。

## Dependency

github.com/gin-gonic/gin

github.com/go-redis/redis/v9

## Inspiration

https://github.com/imtoori/gin-redis-ip-limiter

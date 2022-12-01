# Gin Redis Limiter

## Description

Gin middleware with redis to limit client request by IP.

The limiter can limit some path request counts within a sliding window duration, and block the client.

If too many request to the path detected, block the client that it will not have access to the path for a while.

Gin + Redis实现IP黑名单功能。在一段时间内，如果某个IP对注册了中间件的路由访问频率过高，就将这个IP拉黑，限制它的访问。

## Example

```go
package main

import (
	limiter "github.com/micplus/gin-redis-limiter"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func main() {
    rdb := redis.NewClient(&redis.Options{
        Addr: ":6379",
        DB: 0,
    })
    r := gin.Default()
    // if some IP requested for over 30 times within 5 seconds, 
    // ban it for 15 seconds.
    limit := 30
	slidingWindow := time.Second * 5
	banTime := time.Second * 15
    r.Use(limiter.Limiter(rdb, "test", 
        limit, slidingWindow, 
        "testban", banTime))
    // r.GET("/", SomeHandler)
    // ...
    r.Run(":8080")
}
```

## Dependency

github.com/gin-gonic/gin

github.com/go-redis/redis/v9

## Inspiration

https://github.com/imtoori/gin-redis-ip-limiter

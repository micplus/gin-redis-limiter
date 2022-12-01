package limiter

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func Limiter(rdb *redis.Client, userKey string, limitCount int, slidingWindow time.Duration, banKey string, banTime time.Duration) gin.HandlerFunc {
	bg := context.Background()
	if err := rdb.Ping(bg).Err(); err != nil {
		panic(err)
	}

	return func(ctx *gin.Context) {
		// if client is already banned
		banKey := ctx.ClientIP() + "+" + banKey
		_, err := rdb.Get(bg, banKey).Int64()
		if err == nil { // not redis.Nil, already banned!
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		// not banned yet
		userKey := ctx.ClientIP() + "-" + userKey
		now := time.Now()

		// make it atomic
		pipe := rdb.TxPipeline()
		// remove outdate history
		_ = pipe.ZRemRangeByScore(bg, userKey, "0",
			strconv.Itoa(int(now.UnixNano()-slidingWindow.Nanoseconds())))
		// list of request history
		historyCmd := pipe.ZRange(bg, userKey, 0, -1)
		// add current request, append only
		_ = pipe.ZAddNX(bg, userKey, redis.Z{
			Score:  float64(now.UnixNano()),
			Member: now.UnixNano(),
		})
		_ = pipe.Expire(bg, userKey, slidingWindow)
		_, _ = pipe.Exec(bg)

		// get history request list
		history, _ := historyCmd.Result()
		if len(history) >= limitCount { // request too frequently
			// ban it!
			rdb.Set(bg, banKey, now.UnixNano(), banTime)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		ctx.Next()
	}
}

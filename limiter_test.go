package limiter_test

import (
	"net/http"
	"testing"
	"time"

	limiter "github.com/micplus/gin-redis-limiter"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func TestLimiter(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
		DB:   1,
	})
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	limit := 30
	slidingWindow := time.Second * 5
	banTime := time.Second * 15
	r.Use(limiter.Limiter(rdb, "test", limit, slidingWindow, "testban", banTime))
	r.GET("/", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

	go r.Run(":8080")

	for i := 0; i < limit+3; i++ {
		resp, err := http.Get("http://localhost:8080")
		if err != nil {
			t.Error("request error: ", err)
			return
		}

		switch i {
		case limit:
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Error("limit not works")
				return
			}
			time.Sleep(slidingWindow)
		case limit + 1:
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Error("ban not works")
				return
			}
			time.Sleep(banTime - slidingWindow)
		case limit + 2:
			if resp.StatusCode != http.StatusOK {
				t.Error("unwanted ban")
				return
			}
		}
	}
}

package ratelimit

import (
	"goweb/utils/parsecfg"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

// gin 限流
var rl = ratelimit.New(parsecfg.GlobalConfig.RateLimit) // allow RateLimit req per second
// RateLimit protected service
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		rl.Take()
	}
}

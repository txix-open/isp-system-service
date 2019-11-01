package invoker

import (
	"github.com/integration-system/isp-lib/redis"
)

var (
	RedisClient = redis.NewRxClient()
)

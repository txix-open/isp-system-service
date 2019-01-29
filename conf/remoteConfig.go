package conf

import (
	"github.com/integration-system/isp-lib/database"
	"github.com/integration-system/isp-lib/structure"
)

type RemoteConfig struct {
	DB                     database.DBConfiguration       `schema:"Database"`
	RedisAddress           structure.AddressConfiguration `schema:"Redis"`
	DefaultTokenExpireTime int64                          `schema:"Default token expire time,In milliseconds"`
	Metrics                structure.MetricConfiguration  `schema:"Metrics"`
}

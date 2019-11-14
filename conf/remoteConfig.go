package conf

import (
	"github.com/integration-system/isp-lib/structure"
)

type RemoteConfig struct {
	Database               structure.DBConfiguration     `schema:"Настройка базы данных"`
	Redis                  structure.RedisConfiguration  `schema:"Настройка Redis" valid:"required~Required"`
	DefaultTokenExpireTime int64                         `schema:"Время жизни токена по умолчанию,время жизни токена в миллисекундах с момента его создания. если = -1 - время жизни неограниченно"`
	Metrics                structure.MetricConfiguration `schema:"Настройка метрик"`
	ApplicationSecret      string                        `schema:"Ключ для подписи application токенов" valid:"required~Required"`
}

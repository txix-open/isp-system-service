package conf

import (
	"reflect"

	"github.com/integration-system/isp-kit/dbx"
	"github.com/integration-system/isp-kit/log"
	"github.com/integration-system/isp-kit/rc/schema"
	"github.com/integration-system/isp-lib/v2/structure"
	"github.com/integration-system/jsonschema"
)

func init() {
	schema.CustomGenerators.Register("logLevel", func(field reflect.StructField, t *jsonschema.Type) {
		t.Type = "string"
		t.Enum = []interface{}{"debug", "info", "error", "fatal"}
	})
}

type Remote struct {
	Database               dbx.Config                   `schema:"Настройка базы данных"`
	LogLevel               log.Level                    `schemaGen:"logLevel" schema:"Уровень логирования"`
	Redis                  structure.RedisConfiguration `schema:"Настройка Redis" valid:"required~Required"`
	DefaultTokenExpireTime int                          `schema:"Время жизни токена по умолчанию,время жизни токена в миллисекундах с момента его создания. если = -1 - время жизни неограниченно"`
	ApplicationSecret      string                       `schema:"Ключ для подписи application токенов" valid:"required~Required"`
}

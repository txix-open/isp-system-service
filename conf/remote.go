package conf

import (
	"reflect"

	"github.com/integration-system/isp-kit/dbx"
	"github.com/integration-system/isp-kit/log"
	"github.com/integration-system/isp-kit/rc/schema"
	"github.com/integration-system/jsonschema"
)

func init() {
	schema.CustomGenerators.Register("logLevel", func(field reflect.StructField, t *jsonschema.Type) {
		t.Type = "string"
		t.Enum = []interface{}{"debug", "info", "error", "fatal"}
	})
}

type Remote struct {
	Database dbx.Config `schema:"Настройка базы данных"`
	LogLevel log.Level  `schemaGen:"logLevel" schema:"Уровень логирования"`
}

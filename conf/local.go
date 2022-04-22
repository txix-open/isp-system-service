package conf

import (
	"github.com/integration-system/isp-kit/bootstrap"
)

type Local struct {
	*bootstrap.LocalConfig
	InstanceUuid string
}

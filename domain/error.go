package domain

import (
	"github.com/pkg/errors"
)

var (
	ErrDomainNotFound      = errors.New("domain not found")
	ErrDomainDuplicateName = errors.New("domain name already exist")

	ErrServiceNotFound      = errors.New("service not found")
	ErrServiceDuplicateName = errors.New("service name already exist")

	ErrSystemNotFound = errors.New("system not found")

	ErrApplicationNotFound      = errors.New("application not found")
	ErrApplicationDuplicateName = errors.New("application name already exist")
)

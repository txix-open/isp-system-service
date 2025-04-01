package domain

import (
	"github.com/pkg/errors"
)

var (
	ErrDomainNotFound      = errors.New("domain not found")
	ErrDomainDuplicateName = errors.New("domain name already exist")

	ErrAppGroupNotFound      = errors.New("application group not found")
	ErrAppGroupDuplicateName = errors.New("application group name already exist")

	ErrSystemNotFound = errors.New("system not found")

	ErrApplicationNotFound      = errors.New("application not found")
	ErrApplicationDuplicateName = errors.New("application name already exist")
	ErrApplicationDuplicateId   = errors.New("application id already exist")

	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token is expired")

	ErrAccessListNotFound = errors.New("access_list not found")
)

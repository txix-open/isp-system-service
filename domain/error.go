package domain

import (
	"github.com/txix-open/isp-kit/errors"
)

const (
	ErrCodeInvalidRequest = 600

	ErrCodeApplicationNotFound      = 602
	ErrCodeApplicationDuplicateName = 603
	ErrCodeApplicationDuplicateId   = 604

	ErrCodeAppGroupNotFound      = 605
	ErrCodeAppGroupDuplicateName = 606
)

var (
	ErrAppGroupNotFound      = errors.New("application group not found")
	ErrAppGroupDuplicateName = errors.New("application group name already exist")

	ErrApplicationNotFound      = errors.New("application not found")
	ErrApplicationDuplicateName = errors.New("application name already exist")
	ErrApplicationDuplicateId   = errors.New("application id already exist")

	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token is expired")

	ErrAccessListNotFound = errors.New("access_list not found")
)

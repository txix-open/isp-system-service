package redis

import (
	"context"
	"fmt"
	"time"

	rd "github.com/go-redis/redis/v8"
	"github.com/integration-system/isp-lib/v2/redis"
	"github.com/pkg/errors"
	"isp-system-service/domain"
	"isp-system-service/entity"
)

const (
	systemIdentityFieldInDb      = "1"
	domainIdentityFieldInDb      = "2"
	serviceIdentityFieldInDb     = "3"
	applicationIdentityFieldInDb = "4"
)

type Client struct {
	instanceUuid string
	cli          *redis.RxClient
}

func NewClient(instanceUuid string, cli *redis.RxClient) Client {
	return Client{
		cli:          cli,
		instanceUuid: instanceUuid,
	}
}

func (r Client) SetApplicationToken(ctx context.Context, req entity.RedisSetToken) error {
	_, err := r.cli.UseDbTx(redis.ApplicationTokenDb, func(pipeline rd.Pipeliner) error {
		idMap := map[string]interface{}{
			systemIdentityFieldInDb:      domain.DefaultSystemId,
			domainIdentityFieldInDb:      req.DomainIdentity,
			serviceIdentityFieldInDb:     req.ServiceIdentity,
			applicationIdentityFieldInDb: req.ApplicationIdentity,
		}

		key := fmt.Sprintf("%s|%s", req.Token, r.instanceUuid)
		stat := pipeline.HMSet(ctx, key, idMap)
		err := stat.Err()
		if err != nil {
			return errors.WithMessagef(err, "pipeline hmset")
		}

		if req.ExpireTime > 0 {
			err = pipeline.Expire(ctx, req.Token, time.Duration(req.ExpireTime)*time.Millisecond).Err()
			if err != nil {
				return errors.WithMessage(err, "pipeline expire token")
			}
		}

		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "application_token transaction")
	}

	return nil
}

func (r Client) UpdateApplicationPermission(ctx context.Context, req entity.RedisApplicationPermission) error {
	_, err := r.cli.UseDb(redis.ApplicationPermissionDb, func(pipeline rd.Pipeliner) error {
		key := fmt.Sprintf("%d|%s", req.AppId, req.Method)
		_, err := pipeline.Set(ctx, key, req.Value, 0).Result()
		if err != nil {
			return errors.WithMessage(err, "pipeline set")
		}

		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "application_permission transaction")
	}

	return nil
}

func (r Client) UpdateApplicationPermissionList(ctx context.Context, removed []string, added []interface{}) error {
	_, err := r.cli.UseDb(redis.ApplicationPermissionDb, func(pipeline rd.Pipeliner) error {
		if len(removed) > 0 {
			_, err := pipeline.Del(ctx, removed...).Result()
			if err != nil {
				return errors.WithMessage(err, "pipeline del")
			}
		}

		if len(added) > 0 {
			_, err := pipeline.MSet(ctx, added...).Result()
			if err != nil {
				return errors.WithMessage(err, "pipeline mset")
			}
		}

		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "application_permission transaction")
	}

	return nil
}

func (r Client) DeleteApplication(ctx context.Context, appTokens []string, applicationPermissionList []string) error {
	keys := make([]string, len(appTokens))
	for i, token := range appTokens {
		keys[i] = fmt.Sprintf("%s|%s", token, r.instanceUuid)
	}

	_, err := r.cli.UseDbTx(redis.ApplicationTokenDb, func(pipeline rd.Pipeliner) error {
		if len(appTokens) > 0 {
			_, err := pipeline.Del(ctx, keys...).Result()
			if err != nil {
				return errors.WithMessage(err, "pipeline del")
			}
		}

		if len(applicationPermissionList) > 0 {
			_, err := r.cli.UseDb(redis.ApplicationPermissionDb, func(pipeliner rd.Pipeliner) error {
				_, err := pipeliner.Del(ctx, applicationPermissionList...).Result()
				if err != nil {
					return errors.WithMessage(err, "pipeline del")
				}

				return nil
			})
			if err != nil {
				return errors.WithMessage(err, "application_permission transaction")
			}
		}

		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "application_token transaction")
	}

	return nil
}

func (r Client) DeleteToken(ctx context.Context, tokens []string) error {
	keys := make([]string, len(tokens))
	for i, token := range tokens {
		keys[i] = fmt.Sprintf("%s|%s", token, r.instanceUuid)
	}

	_, err := r.cli.UseDbTx(redis.ApplicationTokenDb, func(p rd.Pipeliner) error {
		_, err := p.Del(ctx, keys...).Result()
		if err != nil {
			return errors.WithMessage(err, "del keys")
		}

		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "application_token transaction")
	}

	return nil
}

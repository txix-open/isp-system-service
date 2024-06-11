package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/isp-kit/dbx/migration"
	"github.com/txix-open/isp-kit/grpc/client"
	"github.com/txix-open/isp-kit/test"
	"github.com/txix-open/isp-kit/test/dbt"
	"github.com/txix-open/isp-kit/test/grpct"
	"isp-system-service/assembly"
	"isp-system-service/conf"
	"isp-system-service/domain"
	"isp-system-service/entity"
	"isp-system-service/migrations"
)

func TestSecureSuite(t *testing.T) {
	suite.Run(t, &SecureSuite{})
}

type SecureSuite struct {
	suite.Suite
	test    *test.Test
	require *require.Assertions
	testDb  *dbt.TestDb
	api     *client.Client
}

func (s *SecureSuite) SetupSuite() {
	s.test, s.require = test.New(s.T())

	s.testDb = dbt.New(s.test)
	migrations.Initialize.SetParams("../migrations", s.testDb.Schema())

	db, err := s.testDb.DB()
	s.require.NoError(err)

	err = migration.NewRunner(db.DB.DB, "../migrations").Run()
	s.require.NoError(err)

	locator := assembly.NewLocator(s.testDb, s.test.Logger())
	handler := locator.Handler(conf.Remote{})
	_, s.api = grpct.TestServer(s.test, handler)

	createdTime := time.Now().UTC()
	InsertDomain(s.testDb, entity.Domain{
		Id: 3, Name: "test_domain", SystemId: 1, CreatedAt: createdTime, UpdatedAt: createdTime,
	})
	InsertService(s.testDb, entity.Service{
		Id: 5, Name: "test_service", DomainId: 3, CreatedAt: createdTime, UpdatedAt: createdTime,
	})
	InsertApplication(s.testDb, entity.Application{
		Id: 7, Name: "test_application", ServiceId: 5, CreatedAt: createdTime, UpdatedAt: createdTime,
	})
}

func (s *SecureSuite) TestAuthenticate_Success() {
	InsertToken(s.testDb, entity.Token{
		Token: "test_token_success", AppId: 7, ExpireTime: -1, CreatedAt: time.Now().UTC(),
	})

	result := domain.AuthenticateResponse{}
	err := s.api.Invoke("system/secure/authenticate").
		JsonRequestBody(domain.AuthenticateRequest{
			Token: "test_token_success",
		}).
		ReadJsonResponse(&result).
		Do(context.Background())
	s.require.NoError(err)
	s.require.Equal(domain.AuthenticateResponse{
		Authenticated: true,
		ErrorReason:   "",
		AuthData: &domain.AuthData{
			SystemId:      1,
			DomainId:      3,
			ServiceId:     5,
			ApplicationId: 7,
		},
	}, result)
}

func (s *SecureSuite) TestAuthenticate_NotFound() {
	result := domain.AuthenticateResponse{}
	err := s.api.Invoke("system/secure/authenticate").
		JsonRequestBody(domain.AuthenticateRequest{
			Token: "test_token_not_found",
		}).
		ReadJsonResponse(&result).
		Do(context.Background())
	s.require.NoError(err)
	s.require.Equal(domain.AuthenticateResponse{
		Authenticated: false,
		ErrorReason:   domain.ErrTokenNotFound.Error(),
		AuthData:      nil,
	}, result)
}

func (s *SecureSuite) TestAuthenticate_NotExpired() {
	InsertToken(s.testDb, entity.Token{
		Token: "test_token_not_expired", AppId: 7, ExpireTime: int((time.Hour).Milliseconds()), CreatedAt: time.Now().UTC(),
	})

	result := domain.AuthenticateResponse{}
	err := s.api.Invoke("system/secure/authenticate").
		JsonRequestBody(domain.AuthenticateRequest{
			Token: "test_token_not_expired",
		}).
		ReadJsonResponse(&result).
		Do(context.Background())
	s.require.NoError(err)
	s.require.Equal(domain.AuthenticateResponse{
		Authenticated: true,
		ErrorReason:   "",
		AuthData: &domain.AuthData{
			SystemId:      1,
			DomainId:      3,
			ServiceId:     5,
			ApplicationId: 7,
		},
	}, result)
}

func (s *SecureSuite) TestAuthenticate_Expired() {
	InsertToken(s.testDb, entity.Token{
		Token: "test_token_expired", AppId: 7, ExpireTime: 0, CreatedAt: time.Now().UTC(),
	})

	result := domain.AuthenticateResponse{}
	err := s.api.Invoke("system/secure/authenticate").
		JsonRequestBody(domain.AuthenticateRequest{
			Token: "test_token_expired",
		}).
		ReadJsonResponse(&result).
		Do(context.Background())
	s.require.NoError(err)
	s.require.Equal(domain.AuthenticateResponse{
		Authenticated: false,
		ErrorReason:   domain.ErrTokenExpired.Error(),
		AuthData:      nil,
	}, result)
}

func (s *SecureSuite) TestAuthorize_Success_True() {
	InsertAccessList(s.testDb, entity.AccessList{
		AppId:  7,
		Method: "endpoint/available",
		Value:  true,
	})

	result := domain.AuthorizeResponse{}
	err := s.api.Invoke("system/secure/authorize").
		JsonRequestBody(domain.AuthorizeRequest{
			ApplicationId: 7,
			Endpoint:      "endpoint/available",
		}).
		ReadJsonResponse(&result).
		Do(context.Background())
	s.require.NoError(err)
	s.require.Equal(domain.AuthorizeResponse{
		Authorized: true,
	}, result)
}

func (s *SecureSuite) TestAuthorize_Success_False() {
	InsertAccessList(s.testDb, entity.AccessList{
		AppId:  7,
		Method: "endpoint/not_available",
		Value:  false,
	})

	result := domain.AuthorizeResponse{}
	err := s.api.Invoke("system/secure/authorize").
		JsonRequestBody(domain.AuthorizeRequest{
			ApplicationId: 7,
			Endpoint:      "endpoint/not_available",
		}).
		ReadJsonResponse(&result).
		Do(context.Background())
	s.require.NoError(err)
	s.require.Equal(domain.AuthorizeResponse{
		Authorized: false,
	}, result)
}

func (s *SecureSuite) TestAuthorize_NotFound() {
	result := domain.AuthorizeResponse{}
	err := s.api.Invoke("system/secure/authorize").
		JsonRequestBody(domain.AuthorizeRequest{
			ApplicationId: 7,
			Endpoint:      "endpoint/not_found",
		}).
		ReadJsonResponse(&result).
		Do(context.Background())
	s.require.NoError(err)
	s.require.Equal(domain.AuthorizeResponse{
		Authorized: false,
	}, result)
}

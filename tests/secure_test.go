package tests_test

import (
	"context"
	"testing"
	"time"

	"github.com/txix-open/isp-kit/dbx"

	"isp-system-service/assembly"
	"isp-system-service/conf"
	"isp-system-service/domain"
	"isp-system-service/entity"

	"github.com/stretchr/testify/suite"
	"github.com/txix-open/isp-kit/grpc/client"
	"github.com/txix-open/isp-kit/test"
	"github.com/txix-open/isp-kit/test/dbt"
	"github.com/txix-open/isp-kit/test/grpct"
)

func TestSecureSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &SecureSuite{})
}

type SecureSuite struct {
	suite.Suite

	test   *test.Test
	testDb *dbt.TestDb
	api    *client.Client
}

func (s *SecureSuite) SetupSuite() {
	s.test, _ = test.New(s.T())

	s.testDb = dbt.New(s.test, dbx.WithMigrationRunner("../migrations", s.test.Logger()))

	locator := assembly.NewLocator(s.testDb, s.test.Logger())
	config := locator.Config(conf.Remote{})
	_, s.api = grpct.TestServer(s.test, config.Handler)

	createdTime := time.Now().UTC()
	InsertDomain(s.testDb, entity.Domain{
		Id: 3, Name: "test_domain", SystemId: 1, CreatedAt: createdTime, UpdatedAt: createdTime,
	})
	InsertAppGroup(s.testDb, entity.AppGroup{
		Id: 5, Name: "test_application_group", DomainId: 3, CreatedAt: createdTime, UpdatedAt: createdTime,
	})
	InsertApplication(s.testDb, entity.Application{
		Id: 7, Name: "test_application", ApplicationGroupId: 5, CreatedAt: createdTime, UpdatedAt: createdTime,
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
		JsonResponseBody(&result).
		Do(context.Background())
	s.Require().NoError(err)
	s.Require().Equal(domain.AuthenticateResponse{
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
		JsonResponseBody(&result).
		Do(context.Background())
	s.Require().NoError(err)
	s.Require().Equal(domain.AuthenticateResponse{
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
		JsonResponseBody(&result).
		Do(context.Background())
	s.Require().NoError(err)
	s.Require().Equal(domain.AuthenticateResponse{
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
		JsonResponseBody(&result).
		Do(context.Background())
	s.Require().NoError(err)
	s.Require().Equal(domain.AuthenticateResponse{
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
		JsonResponseBody(&result).
		Do(context.Background())
	s.Require().NoError(err)
	s.Require().Equal(domain.AuthorizeResponse{
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
		JsonResponseBody(&result).
		Do(context.Background())
	s.Require().NoError(err)
	s.Require().Equal(domain.AuthorizeResponse{
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
		JsonResponseBody(&result).
		Do(context.Background())
	s.Require().NoError(err)
	s.Require().Equal(domain.AuthorizeResponse{
		Authorized: false,
	}, result)
}

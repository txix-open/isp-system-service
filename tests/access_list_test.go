package tests_test

import (
	"testing"

	"github.com/txix-open/isp-kit/dbx"

	"isp-system-service/assembly"
	"isp-system-service/conf"
	"isp-system-service/domain"
	"isp-system-service/entity"
	"isp-system-service/repository"

	"github.com/stretchr/testify/suite"
	"github.com/txix-open/isp-kit/grpc/client"
	"github.com/txix-open/isp-kit/test"
	"github.com/txix-open/isp-kit/test/dbt"
	"github.com/txix-open/isp-kit/test/fake"
	"github.com/txix-open/isp-kit/test/grpct"
)

func TestAccessListSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &AccessListSuite{})
}

type AccessListSuite struct {
	suite.Suite
	test           *test.Test
	testDb         *dbt.TestDb
	api            *client.Client
	accessListRepo repository.AccessList
	appId          int
}

func (s *AccessListSuite) SetupTest() {
	s.test, _ = test.New(s.T())

	s.testDb = dbt.New(s.test, dbx.WithMigrationRunner("../migrations", s.test.Logger()))

	locator := assembly.NewLocator(s.testDb, s.test.Logger())
	config := locator.Config(conf.Remote{})
	_, s.api = grpct.TestServer(s.test, config.Handler)
	s.accessListRepo = repository.NewAccessList(s.testDb)

	createdDomain, err := repository.NewDomain(s.testDb).CreateDomain(
		s.T().Context(),
		fake.It[string](),
		fake.It[string](),
		1,
	)
	s.Require().NoError(err)

	appGroup, err := repository.NewAppGroup(s.testDb).CreateAppGroup(
		s.T().Context(),
		fake.It[string](),
		fake.It[string](),
		createdDomain.Id,
	)
	s.Require().NoError(err)

	s.appId = fake.It[int]()

	_, err = repository.NewApplication(s.testDb).CreateApplication(
		s.T().Context(),
		s.appId,
		fake.It[string](),
		fake.It[string](),
		appGroup.Id,
		domain.ApplicationSystemType,
	)
	s.Require().NoError(err)
}

func (s *AccessListSuite) TestGetById_HappyPath() {
	expectedAccessList := []entity.AccessList{
		{
			AppId:  s.appId,
			Method: fake.It[string](),
			Value:  false,
		},
		{
			AppId:  s.appId,
			Method: fake.It[string](),
			Value:  true,
		},
		{
			AppId:  s.appId,
			Method: fake.It[string](),
			Value:  true,
		},
	}

	err := s.accessListRepo.InsertArrayAccessList(s.T().Context(), expectedAccessList)
	s.Require().NoError(err)

	req := domain.Identity{
		Id: s.appId,
	}
	var accessList []domain.MethodInfo
	err = s.api.Invoke("system/access_list/get_by_id").
		JsonRequestBody(&req).
		JsonResponseBody(&accessList).
		Do(s.T().Context())
	s.Require().NoError(err)
	s.Require().Len(accessList, len(expectedAccessList))

	actualAccessList := s.convertAccessList(accessList)
	s.Require().ElementsMatch(expectedAccessList, actualAccessList)
}

func (s *AccessListSuite) TestSetOne_HappyPath() {
	expectedAccessList := []entity.AccessList{
		{
			AppId:  s.appId,
			Method: fake.It[string](),
			Value:  false,
		},
		{
			AppId:  s.appId,
			Method: fake.It[string](),
			Value:  true,
		},
		{
			AppId:  s.appId,
			Method: fake.It[string](),
			Value:  true,
		},
	}

	err := s.accessListRepo.InsertArrayAccessList(s.T().Context(), expectedAccessList)
	s.Require().NoError(err)

	req := domain.AccessListSetOneRequest{
		AppId:  s.appId,
		Method: expectedAccessList[0].Method,
		Value:  true,
	}
	err = s.api.Invoke("system/access_list/set_one").
		JsonRequestBody(&req).
		Do(s.T().Context())
	s.Require().NoError(err)

	actualAccessList, err := s.accessListRepo.GetAccessListByAppId(s.T().Context(), s.appId)
	s.Require().NoError(err)

	for i := range expectedAccessList {
		expectedAccessList[i].Value = true
	}
	s.Require().ElementsMatch(expectedAccessList, actualAccessList)
}

func (s *AccessListSuite) TestSetList_HappyPath() {
	req := domain.AccessListSetListRequest{
		AppId: s.appId,
		Methods: []domain.MethodInfo{
			{
				Method: fake.It[string](),
				Value:  true,
			},
			{
				Method: fake.It[string](),
				Value:  true,
			},
		},
	}

	err := s.api.Invoke("system/access_list/set_list").
		JsonRequestBody(&req).
		Do(s.T().Context())
	s.Require().NoError(err)

	actualAccessList, err := s.accessListRepo.GetAccessListByAppId(s.T().Context(), s.appId)
	s.Require().NoError(err)
	s.Require().Len(actualAccessList, len(req.Methods))

	expectedAccessList := s.convertAccessList(req.Methods)
	s.Require().ElementsMatch(expectedAccessList, actualAccessList)
}

func (s *AccessListSuite) TestDeleteList_HappyPath() {
	accessList := []entity.AccessList{
		{
			AppId:  s.appId,
			Method: fake.It[string](),
			Value:  false,
		},
		{
			AppId:  s.appId,
			Method: fake.It[string](),
			Value:  true,
		},
		{
			AppId:  s.appId,
			Method: fake.It[string](),
			Value:  true,
		},
	}

	err := s.accessListRepo.InsertArrayAccessList(s.T().Context(), accessList)
	s.Require().NoError(err)

	req := domain.AccessListDeleteListRequest{
		AppId: s.appId,
		Methods: []string{
			accessList[0].Method,
			accessList[1].Method,
		},
	}
	err = s.api.Invoke("system/access_list/delete_list").
		JsonRequestBody(&req).
		Do(s.T().Context())
	s.Require().NoError(err)

	actualAccessList, err := s.accessListRepo.GetAccessListByAppId(s.T().Context(), s.appId)
	s.Require().NoError(err)
	s.Require().Len(actualAccessList, 1)
	s.Require().Equal(accessList[2], actualAccessList[0])
}

func (s *AccessListSuite) convertAccessList(methods []domain.MethodInfo) []entity.AccessList {
	converted := make([]entity.AccessList, 0, len(methods))
	for _, method := range methods {
		converted = append(converted, entity.AccessList{
			AppId:  s.appId,
			Method: method.Method,
			Value:  method.Value,
		})
	}
	return converted
}

package tests_test

import (
	"isp-system-service/assembly"
	"isp-system-service/conf"
	"isp-system-service/domain"
	"isp-system-service/entity"
	"isp-system-service/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/txix-open/isp-kit/dbx"
	"github.com/txix-open/isp-kit/grpc/apierrors"
	"github.com/txix-open/isp-kit/grpc/client"
	"github.com/txix-open/isp-kit/test"
	"github.com/txix-open/isp-kit/test/dbt"
	"github.com/txix-open/isp-kit/test/fake"
	"github.com/txix-open/isp-kit/test/grpct"
)

func TestApplicationGroupSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &AppGroupSuite{})
}

type AppGroupSuite struct {
	suite.Suite
	test         *test.Test
	testDb       *dbt.TestDb
	domainRepo   repository.Domain
	appGroupRepo repository.AppGroup
	api          *client.Client
}

func (s *AppGroupSuite) SetupTest() {
	s.test, _ = test.New(s.T())

	s.testDb = dbt.New(s.test, dbx.WithMigrationRunner("../migrations", s.test.Logger()))

	locator := assembly.NewLocator(s.testDb, s.test.Logger())
	config := locator.Config(conf.Remote{})
	_, s.api = grpct.TestServer(s.test, config.Handler)
	s.domainRepo = repository.NewDomain(s.testDb)
	s.appGroupRepo = repository.NewAppGroup(s.testDb)

	q := `
	INSERT INTO domain
	(id, name, description, system_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, description, system_id, created_at, updated_at
	`
	_, err := s.testDb.Exec(s.T().Context(), q, 1, fake.It[string](), fake.It[string](), 1)
	s.Require().NoError(err)
}

func (s *AppGroupSuite) TestGetByIdList() {
	toGenerate := 10
	appGroups := make([]domain.AppGroup, 0, toGenerate)
	ids := make([]int, 0, toGenerate)
	for range 10 {
		appGroup := s.createAppGroup()
		appGroups = append(appGroups, domain.AppGroup{
			Id:          appGroup.Id,
			Name:        appGroup.Name,
			Description: appGroup.Description.String,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		})
		ids = append(ids, appGroup.Id)
	}
	result := []domain.AppGroup{}
	apiReq := domain.IdListRequest{IdList: ids}
	err := s.api.Invoke("system/application_group/get_by_id_list").
		JsonRequestBody(apiReq).
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)
	for i := range result {
		s.Require().NotEmpty(result[i].CreatedAt)
		result[i].CreatedAt = time.Time{}

		s.Require().NotEmpty(result[i].UpdatedAt)
		result[i].UpdatedAt = time.Time{}
	}
	s.Require().ElementsMatch(appGroups, result)
}

func (s *AppGroupSuite) TestGetByIdList_EmptyDescription() {
	toGenerate := 10
	appGroups := make([]domain.AppGroup, 0, toGenerate)
	ids := make([]int, 0, toGenerate)
	for range 10 {
		appGroup, err := s.appGroupRepo.CreateAppGroup(
			s.T().Context(),
			fake.It[string](),
			"",
			1,
		)
		s.Require().NoError(err)
		appGroups = append(appGroups, domain.AppGroup{
			Id:          appGroup.Id,
			Name:        appGroup.Name,
			Description: "",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		})
		ids = append(ids, appGroup.Id)
	}
	result := []domain.AppGroup{}
	apiReq := domain.IdListRequest{IdList: ids}
	err := s.api.Invoke("system/application_group/get_by_id_list").
		JsonRequestBody(apiReq).
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)
	for i := range result {
		s.Require().NotEmpty(result[i].CreatedAt)
		result[i].CreatedAt = time.Time{}

		s.Require().NotEmpty(result[i].UpdatedAt)
		result[i].UpdatedAt = time.Time{}
	}
	s.Require().ElementsMatch(appGroups, result)
}

func (s *AppGroupSuite) TestGetByIdList_NotFoundAppGroups() {
	appGroup := s.createAppGroup()
	result := []domain.AppGroup{}

	apiReq := domain.IdListRequest{IdList: []int{appGroup.Id + 1}}
	err := s.api.Invoke("system/application_group/get_by_id_list").
		JsonRequestBody(apiReq).
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)
	s.Require().Empty(result)
}

func (s *AppGroupSuite) TestCreate_HappyPath() {
	apiReq := domain.CreateAppGroupRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
	}
	var result domain.AppGroup
	err := s.api.Invoke("system/application_group/create").
		JsonRequestBody(apiReq).
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)

	s.Require().Equal(apiReq.Name, result.Name)
	s.Require().Equal(apiReq.Description, result.Description)

	appGroup, err := s.appGroupRepo.GetAppGroupById(s.T().Context(), result.Id)
	s.Require().NoError(err)
	s.Require().Equal(apiReq.Name, appGroup.Name)
	s.Require().Equal(apiReq.Description, appGroup.Description.String)
}

func (s *AppGroupSuite) TestCreate_AppGroupNameNotUnique() {
	appGroup := s.createAppGroup()

	apiReq := domain.CreateAppGroupRequest{
		Name:        appGroup.Name,
		Description: fake.It[string](),
	}
	var result domain.AppGroup
	err := s.api.Invoke("system/application_group/create").
		JsonRequestBody(apiReq).
		JsonResponseBody(&result).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeAppGroupDuplicateName, apiError.ErrorCode)
}

func (s *AppGroupSuite) TestUpdate_HappyPath() {
	appGroupToUpdate := s.createAppGroup()

	apiReq := domain.UpdateAppGroupRequest{
		Id:          appGroupToUpdate.Id,
		Name:        fake.It[string](),
		Description: fake.It[string](),
	}
	var result domain.AppGroup
	err := s.api.Invoke("system/application_group/update").
		JsonRequestBody(apiReq).
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)

	s.Require().Equal(apiReq.Name, result.Name)
	s.Require().Equal(apiReq.Description, result.Description)

	appGroup, err := s.appGroupRepo.GetAppGroupById(s.T().Context(), apiReq.Id)
	s.Require().NoError(err)
	s.Require().Equal(apiReq.Name, appGroup.Name)
	s.Require().Equal(apiReq.Description, appGroup.Description.String)
}

func (s *AppGroupSuite) TestUpdate_AppGroupNameNotUnique() {
	appGroupToUpdate := s.createAppGroup()
	appGroup := s.createAppGroup()

	apiReq := domain.UpdateAppGroupRequest{
		Id:          appGroupToUpdate.Id,
		Name:        appGroup.Name,
		Description: fake.It[string](),
	}
	err := s.api.Invoke("system/application_group/update").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeAppGroupDuplicateName, apiError.ErrorCode)
}

func (s *AppGroupSuite) TestUpdate_AppGroupNotFound() {
	apiReq := domain.UpdateAppGroupRequest{
		Id:          fake.It[int](),
		Name:        fake.It[string](),
		Description: fake.It[string](),
	}
	err := s.api.Invoke("system/application_group/update").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeAppGroupNotFound, apiError.ErrorCode)
}

func (s *AppGroupSuite) TestDeleteList() {
	appGroupToDelete := s.createAppGroup()
	appGroup := s.createAppGroup()

	var result domain.DeleteResponse
	apiReq := domain.IdListRequest{IdList: []int{appGroupToDelete.Id}}
	err := s.api.Invoke("system/application_group/delete_list").
		JsonRequestBody(apiReq).
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)
	s.Require().Equal(1, result.Deleted)

	_, err = s.appGroupRepo.GetAppGroupById(s.T().Context(), appGroup.Id)
	s.Require().NoError(err)
}

func (s *AppGroupSuite) TestGetAll() {
	toGenerate := 10
	appGroups := make([]domain.AppGroup, 0, toGenerate)
	for range 10 {
		appGroup := s.createAppGroup()
		appGroups = append(appGroups, domain.AppGroup{
			Id:          appGroup.Id,
			Name:        appGroup.Name,
			Description: appGroup.Description.String,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		})
	}
	result := []domain.AppGroup{}
	err := s.api.Invoke("system/application_group/get_all").
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)
	for i := range result {
		s.Require().NotEmpty(result[i].CreatedAt)
		result[i].CreatedAt = time.Time{}

		s.Require().NotEmpty(result[i].UpdatedAt)
		result[i].UpdatedAt = time.Time{}
	}
	s.Require().ElementsMatch(appGroups, result)
}

func (s *AppGroupSuite) createAppGroup() *entity.AppGroup {
	appGroup, err := s.appGroupRepo.CreateAppGroup(
		s.T().Context(),
		fake.It[string](),
		fake.It[string](),
		1,
	)
	s.Require().NoError(err)
	return appGroup
}

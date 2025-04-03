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

func TestApplicationSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ApplicationSuite{})
}

type ApplicationSuite struct {
	suite.Suite
	test         *test.Test
	testDb       *dbt.TestDb
	domainRepo   repository.Domain
	appRepo      repository.Application
	appGroupRepo repository.AppGroup
	api          *client.Client
}

func (s *ApplicationSuite) SetupTest() {
	s.test, _ = test.New(s.T())

	s.testDb = dbt.New(s.test, dbx.WithMigrationRunner("../migrations", s.test.Logger()))

	locator := assembly.NewLocator(s.testDb, s.test.Logger())
	config := locator.Config(conf.Remote{})
	_, s.api = grpct.TestServer(s.test, config.Handler)
	s.domainRepo = repository.NewDomain(s.testDb)
	s.appRepo = repository.NewApplication(s.testDb)
	s.appGroupRepo = repository.NewAppGroup(s.testDb)
}

func (s *ApplicationSuite) TestGetAllApplications() {
	toInsert := fake.It[[]domain.Application](fake.MaxSliceSize(2), fake.MinSliceSize(2))
	expectedApps := s.insertApps(toInsert)

	result := []domain.Application{}
	err := s.api.Invoke("system/application/get_all").
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)
	for i := range result {
		s.Require().NotEmpty(result[i].CreatedAt)
		result[i].CreatedAt = time.Time{}

		s.Require().NotEmpty(result[i].UpdatedAt)
		result[i].UpdatedAt = time.Time{}
	}
	s.Require().ElementsMatch(expectedApps, result)
}

func (s *ApplicationSuite) TestCreate_HappyPath() {
	result := domain.ApplicationWithTokens{}
	apiReq := fake.It[domain.CreateApplicationRequest]()
	apiReq.Type = domain.ApplicationSystemType
	appGroup := s.createAppGroup()
	apiReq.ApplicationGroupId = appGroup.Id

	err := s.api.Invoke("system/application/create_application").
		JsonRequestBody(apiReq).
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)

	s.Require().NotEmpty(result.App.CreatedAt)
	result.App.CreatedAt = time.Time{}

	s.Require().NotEmpty(result.App.UpdatedAt)
	result.App.UpdatedAt = time.Time{}

	expectedApp := domain.Application{
		Id:          apiReq.Id,
		Name:        apiReq.Name,
		Description: apiReq.Description,
		Type:        apiReq.Type,
		ServiceId:   apiReq.ApplicationGroupId,
	}
	s.Require().Equal(expectedApp, result.App)
	s.Require().Empty(result.Tokens)
}

func (s *ApplicationSuite) TestCreate_AppGroupNotFound() {
	apiReq := fake.It[domain.CreateApplicationRequest]()
	apiReq.Type = domain.ApplicationSystemType

	err := s.api.Invoke("system/application/create_application").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeAppGroupNotFound, apiError.ErrorCode)
}

func (s *ApplicationSuite) TestCreate_ApplicationNameNotUniqueInAppGroup() {
	apiReq := fake.It[domain.CreateApplicationRequest]()
	apiReq.Type = domain.ApplicationSystemType

	inserted := s.insertApps(fake.It[[]domain.Application](fake.MinSliceSize(1), fake.MaxSliceSize(1)))

	apiReq.ApplicationGroupId = inserted[0].ServiceId
	apiReq.Name = inserted[0].Name

	err := s.api.Invoke("system/application/create_application").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeApplicationDuplicateName, apiError.ErrorCode)
}

func (s *ApplicationSuite) TestCreate_ApplicationIdNotUnique() {
	apiReq := fake.It[domain.CreateApplicationRequest]()
	apiReq.Type = "SYSTEM"
	appGroup := s.createAppGroup()
	apiReq.ApplicationGroupId = appGroup.Id

	inserted := s.insertApps(fake.It[[]domain.Application](fake.MinSliceSize(1), fake.MaxSliceSize(1)))
	apiReq.Id = inserted[0].Id

	err := s.api.Invoke("system/application/create_application").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeApplicationDuplicateId, apiError.ErrorCode)
}

func (s *ApplicationSuite) TestUpdate_HappyPath() {
	inserted := s.insertApps(fake.It[[]domain.Application](fake.MinSliceSize(1), fake.MaxSliceSize(1)))
	apiReq := domain.UpdateApplicationRequest{
		OldId:       inserted[0].Id,
		NewId:       fake.It[int](),
		Name:        fake.It[string](),
		Description: fake.It[string](),
	}

	var result domain.ApplicationWithTokens
	err := s.api.Invoke("system/application/update_application").
		JsonRequestBody(apiReq).
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)

	s.Require().NotEmpty(result.App.CreatedAt)
	result.App.CreatedAt = time.Time{}

	s.Require().NotEmpty(result.App.UpdatedAt)
	result.App.UpdatedAt = time.Time{}

	expectedApp := domain.Application{
		Id:          apiReq.NewId,
		Name:        apiReq.Name,
		Description: apiReq.Description,
		Type:        inserted[0].Type,
		ServiceId:   inserted[0].ServiceId,
	}
	s.Require().Equal(expectedApp, result.App)
	s.Require().Empty(result.Tokens)

	app, err := s.appRepo.GetApplicationById(s.T().Context(), apiReq.NewId)
	s.Require().NoError(err)

	expectedDbApp := entity.Application{
		Id:                 apiReq.NewId,
		Name:               apiReq.Name,
		Description:        app.Description,
		Type:               inserted[0].Type,
		ApplicationGroupId: inserted[0].ServiceId,
	}

	s.Require().NotEmpty(app.CreatedAt)
	app.CreatedAt = time.Time{}

	s.Require().NotEmpty(app.UpdatedAt)
	app.UpdatedAt = time.Time{}
	s.Require().Equal(expectedDbApp, *app)
}

func (s *ApplicationSuite) TestUpdate_ApplicationNameNotUniqueInAppGroup() {
	inserted := s.insertApps(fake.It[[]domain.Application](fake.MinSliceSize(2), fake.MaxSliceSize(2)))
	apiReq := domain.UpdateApplicationRequest{
		NewId:       inserted[0].Id,
		OldId:       inserted[0].Id,
		Name:        inserted[1].Name,
		Description: fake.It[string](),
	}

	err := s.api.Invoke("system/application/update_application").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeApplicationDuplicateName, apiError.ErrorCode)
}

func (s *ApplicationSuite) TestUpdate_AppIdNotUnique() {
	inserted := s.insertApps(fake.It[[]domain.Application](fake.MinSliceSize(2), fake.MaxSliceSize(2)))
	apiReq := domain.UpdateApplicationRequest{
		OldId:       inserted[0].Id,
		NewId:       inserted[1].Id,
		Name:        fake.It[string](),
		Description: fake.It[string](),
	}

	err := s.api.Invoke("system/application/update_application").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeApplicationDuplicateId, apiError.ErrorCode)
}

func (s *ApplicationSuite) TestUpdate_AppNotFound() {
	apiReq := domain.UpdateApplicationRequest{
		OldId:       fake.It[int](),
		NewId:       fake.It[int](),
		Name:        fake.It[string](),
		Description: fake.It[string](),
	}

	err := s.api.Invoke("system/application/update_application").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeApplicationNotFound, apiError.ErrorCode)
}

func (s *ApplicationSuite) createAppGroup() *entity.AppGroup {
	createdDomain, err := s.domainRepo.CreateDomain(
		s.T().Context(),
		fake.It[string](),
		fake.It[string](),
		1,
	)
	s.Require().NoError(err)

	appGroup, err := s.appGroupRepo.CreateAppGroup(
		s.T().Context(),
		fake.It[string](),
		fake.It[string](),
		createdDomain.Id,
	)
	s.Require().NoError(err)
	return appGroup
}

func (s *ApplicationSuite) insertApps(toInsert []domain.Application) []domain.Application {
	appGroup := s.createAppGroup()
	toExpect := make([]domain.Application, 0, len(toInsert))
	for _, app := range toInsert {
		createdApp, err := s.appRepo.CreateApplication(
			s.T().Context(),
			app.Id,
			app.Name,
			app.Description,
			appGroup.Id,
			app.Type,
		)
		s.Require().NoError(err)
		s.Require().NotEmpty(createdApp.CreatedAt)
		s.Require().NotEmpty(createdApp.UpdatedAt)

		toExpect = append(toExpect, domain.Application{
			Id:          createdApp.Id,
			Name:        createdApp.Name,
			Description: createdApp.Description.String,
			ServiceId:   createdApp.ApplicationGroupId,
			Type:        createdApp.Type,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		})
	}
	return toExpect
}

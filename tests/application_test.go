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
	tokenRepo    repository.Token
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
	s.tokenRepo = repository.NewToken(s.testDb)
}

func (s *ApplicationSuite) TestGetAllApplications() {
	expectedApps := s.insertApps(2)

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
	appGroup := s.createAppGroup()

	apiReq := domain.CreateApplicationRequest{
		Id:                 1,
		Name:               fake.It[string](),
		Type:               domain.ApplicationSystemType,
		ApplicationGroupId: appGroup.Id,
	}

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
	apiReq := domain.CreateApplicationRequest{
		Id:                 1,
		Name:               fake.It[string](),
		Type:               domain.ApplicationSystemType,
		ApplicationGroupId: 1,
	}

	err := s.api.Invoke("system/application/create_application").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeAppGroupNotFound, apiError.ErrorCode)
}

func (s *ApplicationSuite) TestCreate_ApplicationNameNotUniqueInAppGroup() {
	inserted := s.insertApps(1)

	apiReq := domain.CreateApplicationRequest{
		Id:                 inserted[0].Id + 1,
		Name:               inserted[0].Name,
		Type:               domain.ApplicationSystemType,
		ApplicationGroupId: inserted[0].ServiceId,
	}

	err := s.api.Invoke("system/application/create_application").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeApplicationDuplicateName, apiError.ErrorCode)
}

func (s *ApplicationSuite) TestCreate_ApplicationIdNotUnique() {
	inserted := s.insertApps(1)

	apiReq := domain.CreateApplicationRequest{
		Id:                 inserted[0].Id,
		Name:               inserted[0].Name,
		Type:               domain.ApplicationSystemType,
		ApplicationGroupId: inserted[0].ServiceId,
	}

	err := s.api.Invoke("system/application/create_application").
		JsonRequestBody(apiReq).
		Do(s.T().Context())
	apiError := apierrors.FromError(err)
	s.Require().NotNil(apiError)
	s.Require().Equal(domain.ErrCodeApplicationDuplicateId, apiError.ErrorCode)
}

func (s *ApplicationSuite) TestUpdate_HappyPath() {
	inserted := s.insertApps(2)
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
	inserted := s.insertApps(2)
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
	inserted := s.insertApps(2)
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

func (s *ApplicationSuite) TestGetByToken() {
	insertedApps := s.insertApps(2)

	token := fake.It[string]()
	expectedApp := domain.GetApplicationByTokenResponse{
		ApplicationId:      insertedApps[0].Id,
		ApplicationGroupId: insertedApps[0].ServiceId,
	}
	_, err := s.tokenRepo.SaveToken(s.T().Context(), token, expectedApp.ApplicationId, fake.It[int]())
	s.Require().NoError(err)

	result := domain.GetApplicationByTokenResponse{}
	err = s.api.Invoke("system/application/get_application_by_token").
		JsonRequestBody(domain.GetApplicationByTokenRequest{Token: token}).
		JsonResponseBody(&result).
		Do(s.T().Context())
	s.Require().NoError(err)
	s.Require().Equal(expectedApp, result)
}

func (s *ApplicationSuite) TestGetByToken_NotFound() {
	s.insertApps(2)

	result := domain.GetApplicationByTokenResponse{}
	err := s.api.Invoke("system/application/get_application_by_token").
		JsonRequestBody(domain.GetApplicationByTokenRequest{Token: fake.It[string]()}).
		JsonResponseBody(&result).
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

func (s *ApplicationSuite) insertApps(appsCount uint) []domain.Application {
	toInsert := fake.It[[]domain.Application](fake.MinSliceSize(appsCount), fake.MaxSliceSize(appsCount))

	appGroup := s.createAppGroup()
	toExpect := make([]domain.Application, 0, len(toInsert))
	for _, app := range toInsert {
		appId, err := s.appRepo.NextApplicationId(s.T().Context())
		s.Require().NoError(err)

		createdApp, err := s.appRepo.CreateApplication(
			s.T().Context(),
			appId,
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

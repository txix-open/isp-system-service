package tests_test

import (
	"isp-system-service/assembly"
	"isp-system-service/conf"
	"isp-system-service/domain"
	"isp-system-service/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/txix-open/isp-kit/dbx"
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

func (s *ApplicationSuite) insertApps(toInsert []domain.Application) []domain.Application {
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
			Description: *createdApp.Description,
			ServiceId:   createdApp.ApplicationGroupId,
			Type:        createdApp.Type,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		})
	}
	return toExpect
}

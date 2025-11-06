package tests_test

import (
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/isp-kit/dbx"
	"github.com/txix-open/isp-kit/grpc/client"
	"github.com/txix-open/isp-kit/test"
	"github.com/txix-open/isp-kit/test/dbt"
	"github.com/txix-open/isp-kit/test/fake"
	"github.com/txix-open/isp-kit/test/grpct"
	"isp-system-service/assembly"
	"isp-system-service/conf"
	"isp-system-service/service/baseline"
	"testing"
)

func TestBaselineSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &BaselineSuite{})
}

type BaselineSuite struct {
	suite.Suite

	test     *test.Test
	testDb   *dbt.TestDb
	baseline baseline.Service
	api      *client.Client
}

func (s *BaselineSuite) SetupSuite() {
	s.test, _ = test.New(s.T())

	s.testDb = dbt.New(s.test, dbx.WithMigrationRunner("../migrations", s.test.Logger()))

	locator := assembly.NewLocator(s.testDb, s.test.Logger())
	config := locator.Config(conf.Remote{
		Baseline: conf.Baseline{
			InitialAdminUiToken: fake.It[string](),
		},
	})
	s.baseline = config.Baseline
	_, s.api = grpct.TestServer(s.test, config.Handler)
}

func (s *BaselineSuite) TestBaselineSuccess() {
	err := s.baseline.Do(s.T().Context())
	s.Require().NoError(err)

	err = s.baseline.Do(s.T().Context())
	s.Require().NoError(err)
}

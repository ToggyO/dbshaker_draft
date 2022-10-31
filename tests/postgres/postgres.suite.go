package postgres

import (
	_ "github.com/lib/pq"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/ToggyO/dbshaker/tests/suites"
)

// TODO: добавить тест на UpTo (создание миграции находу с помощью create)

type PgTestSuite struct {
	suites.ServiceFixtureSuite
}

func (s *PgTestSuite) SetupSuite() {
	s.Init(internal.PostgresDialect, CreatePgConnectionString(suites.NewDbConf()))

}

func (s *PgTestSuite) TestMigrationUp() {
	print("KEKE")
}

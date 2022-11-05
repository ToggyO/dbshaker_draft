package postgres

import (
	dbshaker "github.com/ToggyO/dbshaker/pkg"
	"github.com/stretchr/testify/require"

	"github.com/ToggyO/dbshaker/tests/suites"
	_ "github.com/lib/pq"
)

// TODO: добавить тест на UpTo (создание миграции находу с помощью create)

type PgTestSuite struct {
	suites.ServiceFixtureSuite
}

func (s *PgTestSuite) SetupSuite() {
	s.Init("postgres", CreatePgConnectionString(suites.NewDbConf("postgres/.env")))
}

func (s *PgTestSuite) TestMigrationDownTo() {
	err := dbshaker.DownTo(s.Db, s.MigrationRoot, 15102022005)
	require.NoError(s.Suite.T(), err)

	migrations, err := dbshaker.ListMigrations()
	require.NoError(s.Suite.T(), err)
	require.Len(s.T(), migrations, 2)
}
